package scheduler

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/entity"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ISyncMidsuitScheduler interface {
	AuthOneStep() (*AuthOneStepResponse, error)
	SyncOrganizationType(jwtToken string) error
	SyncOrganization(jwtToken string) error
}

type SyncMidsuitScheduler struct {
	Viper *viper.Viper
	Log   *logrus.Logger
	DB    *gorm.DB
}

func NewSyncMidsuitScheduler(viper *viper.Viper, log *logrus.Logger) ISyncMidsuitScheduler {
	db := config.NewDatabase()
	return &SyncMidsuitScheduler{
		Viper: viper,
		Log:   log,
		DB:    db,
	}
}

func SyncMidsuitSchedulerFactory(viper *viper.Viper, log *logrus.Logger) ISyncMidsuitScheduler {
	return NewSyncMidsuitScheduler(viper, log)
}

type AuthOneStepResponse struct {
	UserID       int    `json:"userId"`
	Language     string `json:"language"`
	MenuTreeID   int    `json:"menuTreeId"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type OrganizationTypeMidsuitResponse struct {
	ID        int    `json:"id"`
	Category  string `json:"category"`
	Name      string `json:"name"`
	ModelName string `json:"model-name"`
}

type OrganizationTypeMidsuitAPIResponse struct {
	PageCount   int                               `json:"page-count"`
	RecordsSize int                               `json:"records-size"`
	SkipRecords int                               `json:"skip-records"`
	RowCount    int                               `json:"row-count"`
	ArrayCount  int                               `json:"array-count"`
	Records     []OrganizationTypeMidsuitResponse `json:"records"`
}

type OrganizationMidsuitResponse struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	OrgType string `json:"org_type"`
	Region  string `json:"region"`
	Address string `json:"address"`
}

type OrganizationMidsuitAPIResponse struct {
	PageCount   int                           `json:"page-count"`
	RecordsSize int                           `json:"records-size"`
	SkipRecords int                           `json:"skip-records"`
	RowCount    int                           `json:"row-count"`
	ArrayCount  int                           `json:"array-count"`
	Records     []OrganizationMidsuitResponse `json:"records"`
}

func (s *SyncMidsuitScheduler) AuthOneStep() (*AuthOneStepResponse, error) {
	payload := map[string]interface{}{
		"userName": s.Viper.GetString("midsuit.username"),
		"password": s.Viper.GetString("midsuit.username") + "321!",
		"parameters": map[string]interface{}{
			"clientId":       s.Viper.GetString("midsuit.client_id"),
			"roleId":         s.Viper.GetString("midsuit.role_id"),
			"organizationId": 0,
		},
	}

	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + s.Viper.GetString("midsuit.auth_endpoint")
	method := "POST"

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[SyncMidsuitScheduler.AuthOneStep] Error when marshalling payload: " + err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[SyncMidsuitScheduler.AuthOneStep] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[SyncMidsuitScheduler.AuthOneStep] Error when fetching response: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return nil, errors.New("[SyncMidsuitScheduler.AuthOneStep] Error when fetching response: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var authResponse AuthOneStepResponse
	if err := json.Unmarshal(bodyBytes, &authResponse); err != nil {
		s.Log.Error(err)
		return nil, errors.New("[SyncMidsuitScheduler.AuthOneStep] Error when unmarshalling response: " + err.Error())
	}

	return &authResponse, nil
}

func (s *SyncMidsuitScheduler) SyncOrganizationType(jwtToken string) error {
	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + "/views/org-type"
	method := "GET"

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncOrganizationType] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncOrganizationType] Error when fetching response: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncOrganizationType] Error when fetching response: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var apiResponse OrganizationTypeMidsuitAPIResponse
	if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncOrganizationType] Error when unmarshalling response: " + err.Error())
	}

	// Process the records
	for _, record := range apiResponse.Records {
		// Process each record as needed
		s.Log.Infof("Processing record: %+v", record)
		orgType := &entity.OrganizationType{
			MidsuitID: strconv.Itoa(record.ID),
			Name:      record.Name,
			Category:  record.Category,
		}

		// Check if the record already exists
		var existingOrgType entity.OrganizationType
		if err := s.DB.Where("midsuit_id = ?", orgType.MidsuitID).First(&existingOrgType).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create the record if it doesn't exist
				if err := s.DB.Create(orgType).Error; err != nil {
					s.Log.Error(err)
					return errors.New("[SyncMidsuitScheduler.SyncOrganizationType] Error when creating record: " + err.Error())
				}
				s.Log.Infof("Created record: %+v", orgType)
			} else {
				s.Log.Error(err)
				return errors.New("[SyncMidsuitScheduler.SyncOrganizationType] Error when querying record: " + err.Error())
			}
		} else {
			// Update the record if it exists
			if err := s.DB.Model(&existingOrgType).Updates(orgType).Error; err != nil {
				s.Log.Error(err)
				return errors.New("[SyncMidsuitScheduler.SyncOrganizationType] Error when updating record: " + err.Error())
			}
			s.Log.Infof("Updated record: %+v", existingOrgType)
		}
	}

	return nil
}

func (s *SyncMidsuitScheduler) SyncOrganization(jwtToken string) error {
	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + "/views/org"
	method := "GET"

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncOrganization] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncOrganization] Error when fetching response: " + err.Error())
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncOrganization] Error when fetching response: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var apiResponse OrganizationMidsuitAPIResponse
	if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncOrganization] Error when unmarshalling response: " + err.Error())
	}

	// Process the records
	for _, record := range apiResponse.Records {
		// Process each record as needed
		s.Log.Infof("Processing record: %+v", record)

		var orgType entity.OrganizationType
		if err := s.DB.Where("name = ?", record.OrgType).First(&orgType).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Log.Errorf("Organization type with ID %d not found", record.OrgType)
				continue
			}
			s.Log.Error(err)
			return errors.New("[SyncMidsuitScheduler.SyncOrganization] Error when querying organization type: " + err.Error())
		}

		org := &entity.Organization{
			MidsuitID:          strconv.Itoa(record.ID),
			Name:               record.Name,
			Region:             record.Region,
			OrganizationTypeID: orgType.ID,
			Address:            record.Address,
		}

		// Check if the record already exists
		var existingOrg entity.Organization
		if err := s.DB.Where("midsuit_id = ?", org.MidsuitID).First(&existingOrg).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create the record if it doesn't exist
				if err := s.DB.Create(org).Error; err != nil {
					s.Log.Error(err)
					return errors.New("[SyncMidsuitScheduler.SyncOrganization] Error when creating record: " + err.Error())
				}
				s.Log.Infof("Created record: %+v", org)
			} else {
				s.Log.Error(err)
				return errors.New("[SyncMidsuitScheduler.SyncOrganization] Error when querying record: " + err.Error())
			}
		} else {
			// Update the record if it exists
			if err := s.DB.Model(&existingOrg).Updates(org).Error; err != nil {
				s.Log.Error(err)
				return errors.New("[SyncMidsuitScheduler.SyncOrganization] Error when updating record: " + err.Error())
			}
			s.Log.Infof("Updated record: %+v", existingOrg)
		}
	}

	return nil
}
