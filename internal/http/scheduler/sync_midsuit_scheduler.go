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
	"sort"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ISyncMidsuitScheduler interface {
	AuthOneStep() (*AuthOneStepResponse, error)
	SyncOrganizationType(jwtToken string) error
	SyncOrganization(jwtToken string) error
	SyncJobLevel(jwtToken string) error
	SyncOrganizationLocation(jwtToken string) error
	SyncOrganizationStructure(jwtToken string) error
	SyncJob(jwtToken string) error
	SyncEmployee(jwtToken string) error
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

type JobLevelMidsuitResponse struct {
	ID    int    `json:"id"`
	UID   string `json:"uid"`
	Level int    `json:"level"`
	Name  string `json:"name"`
}

type JobLevelMidsuitAPIResponse struct {
	PageCount   int                       `json:"page-count"`
	RecordsSize int                       `json:"records-size"`
	SkipRecords int                       `json:"skip-records"`
	RowCount    int                       `json:"row-count"`
	ArrayCount  int                       `json:"array-count"`
	Records     []JobLevelMidsuitResponse `json:"records"`
}

type OrganizationLocationMidsuitResponse struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	OrganizationID int    `json:"organization_id"`
}

type OrganizationLocationMidsuitAPIResponse struct {
	PageCount   int                                   `json:"page-count"`
	RecordsSize int                                   `json:"records-size"`
	SkipRecords int                                   `json:"skip-records"`
	RowCount    int                                   `json:"row-count"`
	ArrayCount  int                                   `json:"array-count"`
	Records     []OrganizationLocationMidsuitResponse `json:"records"`
}

type OrganizationStructureMidsuitResponse struct {
	ID             int    `json:"id"`
	JobLevelID     int    `json:"job_level_id"`
	Name           string `json:"name"`
	OrganizationID int    `json:"organization_id"`
	ParentID       int    `json:"parent_id"`
}

type OrganizationStructureMidsuitAPIResponse struct {
	PageCount   int                                    `json:"page-count"`
	RecordsSize int                                    `json:"records-size"`
	SkipRecords int                                    `json:"skip-records"`
	RowCount    int                                    `json:"row-count"`
	ArrayCount  int                                    `json:"array-count"`
	Records     []OrganizationStructureMidsuitResponse `json:"records"`
}

type JobMidsuitResponse struct {
	ID                      int    `json:"id"`
	Name                    string `json:"name"`
	OrganizationStructureID int    `json:"organization_structure_id"`
	ParentID                int    `json:"parent_id"`
	JobLevelID              int    `json:"job_level_id"`
	Existing                int    `json:"existing"`
	Promotion               int    `json:"promotion"`
	OrganizationID          int    `json:"organization_id"`
}

type JobMidsuitAPIResponse struct {
	PageCount   int                  `json:"page-count"`
	RecordsSize int                  `json:"records-size"`
	SkipRecords int                  `json:"skip-records"`
	RowCount    int                  `json:"row-count"`
	ArrayCount  int                  `json:"array-count"`
	Records     []JobMidsuitResponse `json:"records"`
}

type EmployeeMidsuitResponse struct {
	ID             int    `json:"id"`
	Email          string `json:"email"`
	EndDate        string `json:"end_date"`
	MobilePhone    string `json:"mobile_phone"`
	Name           string `json:"name"`
	OrganizationID int    `json:"organization_id"`
}

type EmployeeMidsuitAPIResponse struct {
	PageCount   int                       `json:"page-count"`
	RecordsSize int                       `json:"records-size"`
	SkipRecords int                       `json:"skip-records"`
	RowCount    int                       `json:"row-count"`
	ArrayCount  int                       `json:"array-count"`
	Records     []EmployeeMidsuitResponse `json:"records"`
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

func (s *SyncMidsuitScheduler) SyncJobLevel(jwtToken string) error {
	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + "/views/job-level"
	method := "GET"

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncJobLevel] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncJobLevel] Error when fetching response: " + err.Error())
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncJobLevel] Error when fetching response: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var apiResponse JobLevelMidsuitAPIResponse
	if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncJobLevel] Error when unmarshalling response: " + err.Error())
	}

	// Process the records
	for _, record := range apiResponse.Records {
		// Process each record as needed
		s.Log.Infof("Processing record: %+v", record)
		jobLevel := &entity.JobLevel{
			MidsuitID: strconv.Itoa(record.ID),
			Level:     strconv.Itoa(record.Level),
			Name:      record.Name,
		}

		// Check if the record already exists
		var existingJobLevel entity.JobLevel
		if err := s.DB.Where("midsuit_id = ?", jobLevel.MidsuitID).First(&existingJobLevel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create the record if it doesn't exist
				if err := s.DB.Create(jobLevel).Error; err != nil {
					s.Log.Error(err)
					return errors.New("[SyncMidsuitScheduler.SyncJobLevel] Error when creating record: " + err.Error())
				}
				s.Log.Infof("Created record: %+v", jobLevel)
			} else {
				s.Log.Error(err)
				return errors.New("[SyncMidsuitScheduler.SyncJobLevel] Error when querying record: " + err.Error())
			}
		} else {
			// Update the record if it exists
			if err := s.DB.Model(&existingJobLevel).Updates(jobLevel).Error; err != nil {
				s.Log.Error(err)
				return errors.New("[SyncMidsuitScheduler.SyncJobLevel] Error when updating record: " + err.Error())
			}
			s.Log.Infof("Updated record: %+v", existingJobLevel)
		}
	}

	return nil
}

func (s *SyncMidsuitScheduler) SyncOrganizationLocation(jwtToken string) error {
	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + "/views/org-location"
	method := "GET"

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncOrganizationLocation] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncOrganizationLocation] Error when fetching response: " + err.Error())
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncOrganizationLocation] Error when fetching response: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var apiResponse OrganizationLocationMidsuitAPIResponse
	if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncOrganizationLocation] Error when unmarshalling response: " + err.Error())
	}

	// Process the records
	for _, record := range apiResponse.Records {
		// Process each record as needed
		s.Log.Infof("Processing record: %+v", record)

		var org entity.Organization
		if err := s.DB.Where("midsuit_id = ?", strconv.Itoa(record.OrganizationID)).First(&org).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Log.Errorf("Organization with ID %d not found", record.OrganizationID)
				continue
			}
			s.Log.Error(err)
			return errors.New("[SyncMidsuitScheduler.SyncOrganizationLocation] Error when querying organization: " + err.Error())
		}

		orgLocation := &entity.OrganizationLocation{
			MidsuitID:      strconv.Itoa(record.ID),
			Name:           record.Name,
			OrganizationID: org.ID,
		}

		// Check if the record already exists
		var existingOrgLocation entity.OrganizationLocation
		if err := s.DB.Where("midsuit_id = ?", orgLocation.MidsuitID).First(&existingOrgLocation).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create the record if it doesn't exist
				if err := s.DB.Create(orgLocation).Error; err != nil {
					s.Log.Error(err)
					return errors.New("[SyncMidsuitScheduler.SyncOrganizationLocation] Error when creating record: " + err.Error())
				}
				s.Log.Infof("Created record: %+v", orgLocation)
			} else {
				s.Log.Error(err)
				return errors.New("[SyncMidsuitScheduler.SyncOrganizationLocation] Error when querying record: " + err.Error())
			}
		} else {
			// Update the record if it exists
			if err := s.DB.Model(&existingOrgLocation).Updates(orgLocation).Error; err != nil {
				s.Log.Error(err)
				return errors.New("[SyncMidsuitScheduler.SyncOrganizationLocation] Error when updating record: " + err.Error())
			}
			s.Log.Infof("Updated record: %+v", existingOrgLocation)
		}
	}

	return nil
}

func (s *SyncMidsuitScheduler) SyncOrganizationStructure(jwtToken string) error {
	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + "/views/org-structure"
	method := "GET"

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncOrganizationStructure] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncOrganizationStructure] Error when fetching response: " + err.Error())
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncOrganizationStructure] Error when fetching response: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var apiResponse OrganizationStructureMidsuitAPIResponse
	if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncOrganizationStructure] Error when unmarshalling response: " + err.Error())
	}

	// Sort records: parents first, then children
	sort.Slice(apiResponse.Records, func(i, j int) bool {
		return apiResponse.Records[i].ParentID < apiResponse.Records[j].ParentID
	})

	// Process the records
	for _, record := range apiResponse.Records {
		// Process each record as needed
		s.Log.Infof("Processing record: %+v", record)

		var org entity.Organization
		if err := s.DB.Where("midsuit_id = ?", strconv.Itoa(record.OrganizationID)).First(&org).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Log.Errorf("Organization with ID %d not found", record.OrganizationID)
				continue
			}
			s.Log.Error(err)
			return errors.New("[SyncMidsuitScheduler.SyncOrganizationStructure] Error when querying organization: " + err.Error())
		}

		var jobLevel entity.JobLevel
		if err := s.DB.Where("midsuit_id = ?", strconv.Itoa(record.JobLevelID)).First(&jobLevel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Log.Errorf("Job level with ID %d not found", record.JobLevelID)
				continue
			}
			s.Log.Error(err)
			return errors.New("[SyncMidsuitScheduler.SyncOrganizationStructure] Error when querying job level: " + err.Error())
		}

		var parentID *uuid.UUID
		if record.ParentID != 0 {
			var parentOrgStructure entity.OrganizationStructure
			if err := s.DB.Where("midsuit_id = ?", strconv.Itoa(record.ParentID)).First(&parentOrgStructure).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					s.Log.Errorf("Parent organization structure with ID %d not found", record.ParentID)
					continue
				}
				s.Log.Error(err)
				return errors.New("[SyncMidsuitScheduler.SyncOrganizationStructure] Error when querying parent organization structure: " + err.Error())
			}
			parentID = &parentOrgStructure.ID
		}

		orgStructure := &entity.OrganizationStructure{
			MidsuitID:      strconv.Itoa(record.ID),
			Name:           record.Name, // Assuming Name is a string, no need to convert to int
			OrganizationID: org.ID,
			JobLevelID:     jobLevel.ID,
			ParentID:       parentID,
		}

		// Check if the record already exists
		var existingOrgStructure entity.OrganizationStructure
		if err := s.DB.Where("midsuit_id = ?", orgStructure.MidsuitID).First(&existingOrgStructure).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create the record if it doesn't exist
				if err := s.DB.Create(orgStructure).Error; err != nil {
					s.Log.Error(err)
					return errors.New("[SyncMidsuitScheduler.SyncOrganizationStructure] Error when creating record: " + err.Error())
				}
				s.Log.Infof("Created record: %+v", orgStructure)
			} else {
				s.Log.Error(err)
				return errors.New("[SyncMidsuitScheduler.SyncOrganizationStructure] Error when querying record: " + err.Error())
			}
		} else {
			// Update the record if it exists
			if err := s.DB.Model(&existingOrgStructure).Updates(orgStructure).Error; err != nil {
				s.Log.Error(err)
				return errors.New("[SyncMidsuitScheduler.SyncOrganizationStructure] Error when updating record: " + err.Error())
			}
			s.Log.Infof("Updated record: %+v", existingOrgStructure)
		}
	}

	return nil
}

func (s *SyncMidsuitScheduler) SyncJob(jwtToken string) error {
	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + "/views/job"
	method := "GET"

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncJob] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncJob] Error when fetching response: " + err.Error())
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncJob] Error when fetching response: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var apiResponse JobMidsuitAPIResponse
	if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncJob] Error when unmarshalling response: " + err.Error())
	}

	// sort records: parents first, then children
	sort.Slice(apiResponse.Records, func(i, j int) bool {
		return apiResponse.Records[i].ParentID < apiResponse.Records[j].ParentID
	})

	// Process the records
	for _, record := range apiResponse.Records {
		// Process each record as needed
		s.Log.Infof("Processing record: %+v", record)

		var jobLevel entity.JobLevel
		if err := s.DB.Where("midsuit_id = ?", strconv.Itoa(record.JobLevelID)).First(&jobLevel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Log.Errorf("Job level with ID %d not found", record.JobLevelID)
				continue
			}
			s.Log.Error(err)
			return errors.New("[SyncMidsuitScheduler.SyncJob] Error when querying job level: " + err.Error())
		}

		var orgStructure entity.OrganizationStructure
		if err := s.DB.Where("midsuit_id = ?", strconv.Itoa(record.OrganizationStructureID)).First(&orgStructure).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Log.Errorf("Organization structure with ID %d not found", record.OrganizationStructureID)
				continue
			}
			s.Log.Error(err)
			return errors.New("[SyncMidsuitScheduler.SyncJob] Error when querying organization structure: " + err.Error())
		}

		var org entity.Organization
		if err := s.DB.Where("midsuit_id = ?", strconv.Itoa(record.OrganizationID)).First(&org).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Log.Errorf("Organization with ID %d not found", orgStructure.OrganizationID)
				continue
			}
			s.Log.Error(err)
			return errors.New("[SyncMidsuitScheduler.SyncJob] Error when querying organization: " + err.Error())
		}

		var parentID *uuid.UUID
		if record.ParentID != 0 {
			var parentJob entity.Job
			if err := s.DB.Where("midsuit_id = ?", strconv.Itoa(record.ParentID)).First(&parentJob).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					s.Log.Errorf("Parent job with ID %d not found", record.ParentID)
					continue
				}
				s.Log.Error(err)
				return errors.New("[SyncMidsuitScheduler.SyncJob] Error when querying parent job: " + err.Error())
			}
			parentID = &parentJob.ID
		}

		job := &entity.Job{
			MidsuitID:               strconv.Itoa(record.ID),
			Name:                    record.Name,
			OrganizationStructureID: orgStructure.ID,
			ParentID:                parentID,
			Existing:                record.Existing,
			Promotion:               record.Promotion,
			JobLevelID:              jobLevel.ID,
			OrganizationID:          org.ID,
		}

		// Check if the record already exists
		var existingJob entity.Job
		if err := s.DB.Where("midsuit_id = ?", job.MidsuitID).First(&existingJob).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create the record if it doesn't exist
				if err := s.DB.Create(job).Error; err != nil {
					s.Log.Error(err)
					return errors.New("[SyncMidsuitScheduler.SyncJob] Error when creating record: " + err.Error())
				}
				s.Log.Infof("Created record: %+v", job)
			} else {
				s.Log.Error(err)
				return errors.New("[SyncMidsuitScheduler.SyncJob] Error when querying record: " + err.Error())
			}
		} else {
			// Update the record if it exists
			if err := s.DB.Model(&existingJob).Updates(job).Error; err != nil {
				s.Log.Error(err)
				return errors.New("[SyncMidsuitScheduler.SyncJob] Error when updating record: " + err.Error())
			}
			s.Log.Infof("Updated record: %+v", existingJob)
		}
	}

	return nil
}

func (s *SyncMidsuitScheduler) SyncEmployee(jwtToken string) error {
	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + "/views/employee"
	method := "GET"

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncEmployee] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncEmployee] Error when fetching response: " + err.Error())
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncEmployee] Error when fetching response: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var apiResponse EmployeeMidsuitAPIResponse
	if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncEmployee] Error when unmarshalling response: " + err.Error())
	}

	// Process the records
	for _, record := range apiResponse.Records {
		// Process each record as needed
		s.Log.Infof("Processing record: %+v", record)

		var org entity.Organization
		if err := s.DB.Where("midsuit_id = ?", strconv.Itoa(record.OrganizationID)).First(&org).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Log.Errorf("Organization with ID %d not found", record.OrganizationID)
				continue
			}
			s.Log.Error(err)
			return errors.New("[SyncMidsuitScheduler.SyncEmployee] Error when querying organization: " + err.Error())
		}

		var endDate time.Time

		if record.EndDate != "" {
			parsedEndDate, err := time.Parse(time.RFC3339, record.EndDate)
			if err != nil {
				s.Log.Error(err)
				return errors.New("[SyncMidsuitScheduler.SyncEmployee] Error when parsing end date: " + err.Error())
			}
			endDate = parsedEndDate
		}

		employee := &entity.Employee{
			MidsuitID:      strconv.Itoa(record.ID),
			Email:          record.Email,
			EndDate:        endDate,
			MobilePhone:    record.MobilePhone,
			Name:           record.Name,
			OrganizationID: org.ID,
		}

		// Check if the record already exists
		var existingEmployee entity.Employee
		if err := s.DB.Where("midsuit_id = ?", employee.MidsuitID).First(&existingEmployee).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create the record if it doesn't exist
				if err := s.DB.Create(employee).Error; err != nil {
					s.Log.Error(err)
					return errors.New("[SyncMidsuitScheduler.SyncEmployee] Error when creating record: " + err.Error())
				}
				s.Log.Infof("Created record: %+v", employee)
			} else {
				s.Log.Error(err)
				return errors.New("[SyncMidsuitScheduler.SyncEmployee] Error when querying record: " + err.Error())
			}
		} else {
			// Update the record if it exists
			if err := s.DB.Model(&existingEmployee).Updates(employee).Error; err != nil {
				s.Log.Error(err)
				return errors.New("[SyncMidsuitScheduler.SyncEmployee] Error when updating record: " + err.Error())
			}
			s.Log.Infof("Updated record: %+v", existingEmployee)
		}

		// create or update user
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("changeme"), bcrypt.DefaultCost)
		if err != nil {
			s.Log.Error("[UserUseCase.Register] " + err.Error())
			return err
		}
		user := &entity.User{
			Email:    employee.Email,
			Username: employee.Email,
			Password: string(hashedPassword),
			Name:     employee.Name,
			EmployeeID: func() *uuid.UUID {
				if existingEmployee.ID != uuid.Nil {
					return &existingEmployee.ID
				}
				return nil
			}(),
		}

		var existingUser entity.User
		if err := s.DB.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := s.DB.Create(user).Error; err != nil {
					s.Log.Error(err)
					return errors.New("[SyncMidsuitScheduler.SyncEmployee] Error when creating user: " + err.Error())
				}
				s.Log.Infof("Created user: %+v", user)
			} else {
				s.Log.Error(err)
				return errors.New("[SyncMidsuitScheduler.SyncEmployee] Error when querying user: " + err.Error())
			}
		} else {
			if err := s.DB.Model(&existingUser).Updates(user).Error; err != nil {
				s.Log.Error(err)
				return errors.New("[SyncMidsuitScheduler.SyncEmployee] Error when updating user: " + err.Error())
			}
			s.Log.Infof("Updated user: %+v", existingUser)
		}
	}

	return nil
}
