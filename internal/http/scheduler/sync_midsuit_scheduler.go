package scheduler

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/entity"
	messaging "app/go-sso/internal/messaging/user"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
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
	SyncEmployeeJob(jwtToken string) error
	SyncUserProfile(jwtToken string) error
	SyncGrade(jwtToken string) error
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

type CategoryMidsuitResponse struct {
	PropertyLabel string `json:"propertyLabel"`
	ID            string `json:"id"`
}

type OrganizationTypeMidsuitResponse struct {
	ID int `json:"id"`
	// Category  string `json:"category"`
	Category  CategoryMidsuitResponse `json:"category"`
	Name      string                  `json:"name"`
	ModelName string                  `json:"model-name"`
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

type AdOrgMidsuitResponse struct {
	ID            int    `json:"id"`
	PropertyLabel string `json:"propertyLabel"`
}

type HcEmployeeMidSuitResponse struct {
	ID            int    `json:"id"`
	PropertyLabel string `json:"propertyLabel"`
}

type OrganizationLocationMidsuitResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	// OrganizationID int    `json:"organization_id"`
	AdOrgID AdOrgMidsuitResponse `json:"AD_Org_ID"`
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
	ID         int    `json:"id"`
	JobLevelID int    `json:"job_level_id"`
	Name       string `json:"name"`
	// OrganizationID int    `json:"organization_id"`
	AdOrgID  AdOrgMidsuitResponse `json:"Ad_Org_ID"`
	ParentID int                  `json:"Parent_ID"`
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
	Existing                string `json:"existing"`
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
	RetirementDate string `json:"retirement_date"`
	MobilePhone    string `json:"mobile_phone"`
	Name           string `json:"name"`
	// OrganizationID int    `json:"organization_id"`
	ADOrgID AdOrgMidsuitResponse `json:"AD_Org_ID"`
}

type EmployeeMidsuitAPIResponse struct {
	PageCount   int                       `json:"page-count"`
	RecordsSize int                       `json:"records-size"`
	SkipRecords int                       `json:"skip-records"`
	RowCount    int                       `json:"row-count"`
	ArrayCount  int                       `json:"array-count"`
	Records     []EmployeeMidsuitResponse `json:"records"`
}

type HcEmployeeGradeIDResponse struct {
	PropertyLabel string `json:"propertyLabel"`
	ID            int    `json:"id"`
	Identifier    string `json:"identifier"`
	ModelName     string `json:"model-name"`
}

type EmployeeJobMidsuitResponse struct {
	ID                      int                       `json:"id"`
	EmployeeID              int                       `json:"employee_id"`
	EmpOrganizationID       int                       `json:"emp_organization_id"`
	JobID                   int                       `json:"job_id"`
	JobLevelID              int                       `json:"job_level_id"`
	OrganizationID          int                       `json:"organization_id"`
	OrganizationStructureID int                       `json:"organization_structure_id"`
	OrganizationLocationID  int                       `json:"organization_location_id"`
	HCEmployeeGradeID       HcEmployeeGradeIDResponse `json:"HC_EmployeeGrade_ID"`
}

type EmployeeJobMidsuitAPIResponse struct {
	PageCount   int                          `json:"page-count"`
	RecordsSize int                          `json:"records-size"`
	SkipRecords int                          `json:"skip-records"`
	RowCount    int                          `json:"row-count"`
	ArrayCount  int                          `json:"array-count"`
	Records     []EmployeeJobMidsuitResponse `json:"records"`
}

type UserProfileMidsuitResponse struct {
	ID            int                       `json:"id"`
	ADOrgID       AdOrgMidsuitResponse      `json:"AD_Org_ID"`
	HCEmployeeID  HcEmployeeMidSuitResponse `json:"HC_Employee_ID"`
	Age           int                       `json:"age"`
	Name          string                    `json:"name"`
	BirthDate     string                    `json:"birth_date"`
	BirthPlace    string                    `json:"birth_place"`
	MaritalStatus string                    `json:"marital_status"`
	PhoneNumber   string                    `json:"phone_number"`
}

type UserProfileMidsuitAPIResponse struct {
	PageCount   int                          `json:"page-count"`
	RecordsSize int                          `json:"records-size"`
	SkipRecords int                          `json:"skip-records"`
	RowCount    int                          `json:"row-count"`
	ArrayCount  int                          `json:"array-count"`
	Records     []UserProfileMidsuitResponse `json:"records"`
}

type HcJobLevelIDResponse struct {
	PropertyLabel string `json:"propertyLabel"`
	ID            int    `json:"id"`
	Identifier    string `json:"identifier"`
	ModelName     string `json:"model-name"`
}

type GradeMidsuitResponse struct {
	ID           int                  `json:"id"`
	UID          string               `json:"uid"`
	HcJobLevelID HcJobLevelIDResponse `json:"HC_JobLevel_ID"`
	Name         string               `json:"name"`
	ModelName    string               `json:"model-name"`
}

type GradeMidsuitAPIResponse struct {
	PageCount   int                    `json:"page-count"`
	RecordsSize int                    `json:"records-size"`
	SkipRecords int                    `json:"skip-records"`
	RowCount    int                    `json:"row-count"`
	ArrayCount  int                    `json:"array-count"`
	Records     []GradeMidsuitResponse `json:"records"`
}

func (s *SyncMidsuitScheduler) AuthOneStep() (*AuthOneStepResponse, error) {
	payload := map[string]interface{}{
		"userName": s.Viper.GetString("midsuit.username"),
		// "password": s.Viper.GetString("midsuit.username") + "321!",
		"password": "JgiMidsuit123!",
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

	var orgTypeMidsuitIDs []string

	// Process the records
	for _, record := range apiResponse.Records {
		// Process each record as needed
		orgTypeMidsuitIDs = append(orgTypeMidsuitIDs, strconv.Itoa(record.ID))

		s.Log.Infof("Processing record: %+v", record)
		orgType := &entity.OrganizationType{
			MidsuitID: strconv.Itoa(record.ID),
			Name:      record.Name,
			Category:  record.Category.ID,
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

	// Delete organization types that are not in the response
	var existingOrgTypes []entity.OrganizationType
	if err := s.DB.Where("midsuit_id NOT IN ?", orgTypeMidsuitIDs).Find(&existingOrgTypes).Error; err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncOrganizationType] Error when querying organization types: " + err.Error())
	}

	for _, existingOrgType := range existingOrgTypes {
		if err := s.DB.Delete(&existingOrgType).Error; err != nil {
			s.Log.Error(err)
			return errors.New("[SyncMidsuitScheduler.SyncOrganizationType] Error when deleting organization type: " + err.Error())
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

	var orgMidsuitIDs []string
	// Process the records
	for _, record := range apiResponse.Records {
		// Process each record as needed
		orgMidsuitIDs = append(orgMidsuitIDs, strconv.Itoa(record.ID))

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

	// Delete organizations that are not in the response
	var existingOrgs []entity.Organization
	if err := s.DB.Where("midsuit_id NOT IN ?", orgMidsuitIDs).Find(&existingOrgs).Error; err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncOrganization] Error when querying organizations: " + err.Error())
	}

	for _, existingOrg := range existingOrgs {
		if err := s.DB.Delete(&existingOrg).Error; err != nil {
			s.Log.Error(err)
			return errors.New("[SyncMidsuitScheduler.SyncOrganization] Error when deleting organization: " + err.Error())
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

	var jobLevelMidsuitIDs []string
	// Process the records
	for _, record := range apiResponse.Records {
		// Process each record as needed
		jobLevelMidsuitIDs = append(jobLevelMidsuitIDs, strconv.Itoa(record.ID))

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

	// Delete job levels that are not in the response
	var existingJobLevels []entity.JobLevel
	if err := s.DB.Where("midsuit_id NOT IN ?", jobLevelMidsuitIDs).Find(&existingJobLevels).Error; err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncJobLevel] Error when querying job levels: " + err.Error())
	}

	for _, existingJobLevel := range existingJobLevels {
		if err := s.DB.Delete(&existingJobLevel).Error; err != nil {
			s.Log.Error(err)
			return errors.New("[SyncMidsuitScheduler.SyncJobLevel] Error when deleting job level: " + err.Error())
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

	var orgLocationMidsuitIDs []string
	// Process the records
	for _, record := range apiResponse.Records {
		// Process each record as needed
		orgLocationMidsuitIDs = append(orgLocationMidsuitIDs, strconv.Itoa(record.ID))

		s.Log.Infof("Processing record: %+v", record)

		var org entity.Organization
		if err := s.DB.Where("midsuit_id = ?", strconv.Itoa(record.AdOrgID.ID)).First(&org).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Log.Errorf("Organization with ID %d not found", record.AdOrgID.ID)
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

	// Delete organization locations that are not in the response
	var existingOrgLocations []entity.OrganizationLocation
	if err := s.DB.Where("midsuit_id NOT IN ?", orgLocationMidsuitIDs).Find(&existingOrgLocations).Error; err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncOrganizationLocation] Error when querying organization locations: " + err.Error())
	}

	for _, existingOrgLocation := range existingOrgLocations {
		if err := s.DB.Delete(&existingOrgLocation).Error; err != nil {
			s.Log.Error(err)
			return errors.New("[SyncMidsuitScheduler.SyncOrganizationLocation] Error when deleting organization location: " + err.Error())
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

	var orgStructureMidsuitIDs []string
	// Process the records
	for _, record := range apiResponse.Records {
		// Process each record as needed
		orgStructureMidsuitIDs = append(orgStructureMidsuitIDs, strconv.Itoa(record.ID))

		s.Log.Infof("Processing record: %+v", record)

		var org entity.Organization
		if err := s.DB.Where("midsuit_id = ?", strconv.Itoa(record.AdOrgID.ID)).First(&org).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Log.Errorf("Organization with ID %d not found", record.AdOrgID.ID)
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

	// Delete organization structures that are not in the response
	var existingOrgStructures []entity.OrganizationStructure
	if err := s.DB.Where("midsuit_id NOT IN ?", orgStructureMidsuitIDs).Find(&existingOrgStructures).Error; err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncOrganizationStructure] Error when querying organization structures: " + err.Error())
	}

	for _, existingOrgStructure := range existingOrgStructures {
		if err := s.DB.Delete(&existingOrgStructure).Error; err != nil {
			s.Log.Error(err)
			return errors.New("[SyncMidsuitScheduler.SyncOrganizationStructure] Error when deleting organization structure: " + err.Error())
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

	var jobMidsuitIDs []string
	// Process the records
	for _, record := range apiResponse.Records {
		// Process each record as needed
		jobMidsuitIDs = append(jobMidsuitIDs, strconv.Itoa(record.ID))

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
			Existing: func() int {
				existing, err := strconv.Atoi(record.Existing)
				if err != nil {
					s.Log.Error(err)
					return 0
				}
				return existing
			}(),
			Promotion:      record.Promotion,
			JobLevelID:     jobLevel.ID,
			OrganizationID: org.ID,
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

	// Delete jobs that are not in the response
	var existingJobs []entity.Job
	if err := s.DB.Where("midsuit_id NOT IN ?", jobMidsuitIDs).Find(&existingJobs).Error; err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncJob] Error when querying jobs: " + err.Error())
	}

	for _, existingJob := range existingJobs {
		if err := s.DB.Delete(&existingJob).Error; err != nil {
			s.Log.Error(err)
			return errors.New("[SyncMidsuitScheduler.SyncJob] Error when deleting job: " + err.Error())
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

	var employeeMidsuitIDs []string
	// Process the records
	for _, record := range apiResponse.Records {
		// Process each record as needed
		employeeMidsuitIDs = append(employeeMidsuitIDs, strconv.Itoa(record.ID))

		s.Log.Infof("Processing record: %+v", record)

		var org entity.Organization
		if err := s.DB.Where("midsuit_id = ?", strconv.Itoa(record.ADOrgID.ID)).First(&org).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Log.Errorf("Organization with ID %d not found", record.ADOrgID.ID)
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

		var retirementDate time.Time

		if record.RetirementDate != "" {
			parsedRetirementDate, err := time.Parse(time.RFC3339, record.RetirementDate)
			if err != nil {
				s.Log.Error(err)
				return errors.New("[SyncMidsuitScheduler.SyncEmployee] Error when parsing retirement date: " + err.Error())
			}
			retirementDate = parsedRetirementDate
		}

		employee := &entity.Employee{
			MidsuitID:      strconv.Itoa(record.ID),
			Email:          record.Email,
			EndDate:        endDate,
			RetirementDate: retirementDate,
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

		var employeeExist entity.Employee
		if err := s.DB.Where("midsuit_id = ?", employee.MidsuitID).First(&employeeExist).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Log.Errorf("Employee with ID %d not found", record.ID)
				employeeExist.ID = uuid.Nil
			} else {
				s.Log.Error(err)
				return errors.New("[SyncMidsuitScheduler.SyncEmployee] Error when querying employee: " + err.Error())
			}
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
				if employeeExist.ID != uuid.Nil {
					return &employeeExist.ID
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

	// Delete employees that are not in the response
	var existingEmployees []entity.Employee
	if err := s.DB.Where("midsuit_id NOT IN ?", employeeMidsuitIDs).Find(&existingEmployees).Error; err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncEmployee] Error when querying employees: " + err.Error())
	}

	for _, existingEmployee := range existingEmployees {
		if err := s.DB.Delete(&existingEmployee).Error; err != nil {
			s.Log.Error(err)
			return errors.New("[SyncMidsuitScheduler.SyncEmployee] Error when deleting employee: " + err.Error())
		}
	}

	return nil
}

func (s *SyncMidsuitScheduler) SyncEmployeeJob(jwtToken string) error {
	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + "/views/employee-job"
	method := "GET"

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncEmployeeJob] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncEmployeeJob] Error when fetching response: " + err.Error())
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncEmployeeJob] Error when fetching response: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var apiResponse EmployeeJobMidsuitAPIResponse
	if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncEmployeeJob] Error when unmarshalling response: " + err.Error())
	}

	var employeeJobMidsuitIDs []string
	// Process the records
	for _, record := range apiResponse.Records {
		// Process each record as needed
		employeeJobMidsuitIDs = append(employeeJobMidsuitIDs, strconv.Itoa(record.ID))

		s.Log.Infof("Processing record: %+v", record)

		var employee entity.Employee
		if err := s.DB.Where("midsuit_id = ?", strconv.Itoa(record.EmployeeID)).First(&employee).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Log.Errorf("Employee with ID %d not found", record.EmployeeID)
				continue
			}
			s.Log.Error(err)
			return errors.New("[SyncMidsuitScheduler.SyncEmployeeJob] Error when querying employee: " + err.Error())
		}

		var job entity.Job
		if err := s.DB.Where("midsuit_id = ?", strconv.Itoa(record.JobID)).First(&job).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Log.Errorf("Job with ID %d not found", record.JobID)
				continue
			}
			s.Log.Error(err)
			return errors.New("[SyncMidsuitScheduler.SyncEmployeeJob] Error when querying job: " + err.Error())
		}

		var empOrganizationID *uuid.UUID
		if record.EmpOrganizationID != 0 {
			var org entity.Organization
			if err := s.DB.Where("midsuit_id = ?", strconv.Itoa(record.EmpOrganizationID)).First(&org).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					s.Log.Errorf("Organization with ID %d not found", record.EmpOrganizationID)
					empOrganizationID = nil
				}
				s.Log.Error(err)
				empOrganizationID = nil
			} else {
				empOrganizationID = &org.ID
			}
		}

		var orgLocationID *uuid.UUID
		if record.OrganizationLocationID != 0 {
			var orgLocation entity.OrganizationLocation
			if err := s.DB.Where("midsuit_id = ?", strconv.Itoa(record.OrganizationLocationID)).First(&orgLocation).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					s.Log.Errorf("Organization location with ID %d not found", record.OrganizationLocationID)
					orgLocationID = nil
				}
				s.Log.Error(err)
				orgLocationID = nil
			}
			orgLocationID = &orgLocation.ID
		}

		var orgStructureID *uuid.UUID
		if record.OrganizationStructureID != 0 {
			var orgStructure entity.OrganizationStructure
			if err := s.DB.Where("midsuit_id = ?", strconv.Itoa(record.OrganizationStructureID)).First(&orgStructure).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					s.Log.Errorf("Organization structure with ID %d not found", record.OrganizationStructureID)
					orgStructureID = nil
				}
				s.Log.Error(err)
				orgStructureID = nil
			} else {
				orgStructureID = &orgStructure.ID
			}
		}

		var gradeID *uuid.UUID
		if record.HCEmployeeGradeID.ID != 0 {
			var grade entity.Grade
			if err := s.DB.Where("midsuit_id = ?", strconv.Itoa(record.HCEmployeeGradeID.ID)).First(&grade).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					s.Log.Errorf("Grade with ID %d not found", record.HCEmployeeGradeID.ID)
					gradeID = nil
				}
				s.Log.Error(err)
				gradeID = nil
			} else {
				gradeID = &grade.ID
			}
		}

		employeeJob := &entity.EmployeeJob{
			MidsuitID:  strconv.Itoa(record.ID),
			EmployeeID: &employee.ID,
			JobID:      job.ID,
			EmpOrganizationID: func() *uuid.UUID {
				if empOrganizationID != nil {
					return empOrganizationID
				}
				return nil
			}(),
			OrganizationLocationID:  *orgLocationID,
			OrganizationStructureID: *orgStructureID,
			GradeID:                 gradeID,
		}

		// Check if the record already exists
		var existingEmployeeJob entity.EmployeeJob
		if err := s.DB.Where("midsuit_id = ?", employeeJob.MidsuitID).First(&existingEmployeeJob).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create the record if it doesn't exist
				if err := s.DB.Create(employeeJob).Error; err != nil {
					s.Log.Error(err)
					return errors.New("[SyncMidsuitScheduler.SyncEmployeeJob] Error when creating record: " + err.Error())
				}
				s.Log.Infof("Created record: %+v", employeeJob)
			} else {
				s.Log.Error(err)
				return errors.New("[SyncMidsuitScheduler.SyncEmployeeJob] Error when querying record: " + err.Error())
			}
		} else {
			// Update the record if it exists
			if err := s.DB.Model(&existingEmployeeJob).Updates(employeeJob).Error; err != nil {
				s.Log.Error(err)
				return errors.New("[SyncMidsuitScheduler.SyncEmployeeJob] Error when updating record: " + err.Error())
			}
			s.Log.Infof("Updated record: %+v", existingEmployeeJob)
		}
	}

	// Delete employee jobs that are not in the response
	var existingEmployeeJobs []entity.EmployeeJob
	if err := s.DB.Where("midsuit_id NOT IN ?", employeeJobMidsuitIDs).Find(&existingEmployeeJobs).Error; err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncEmployeeJob] Error when querying employee jobs: " + err.Error())
	}

	for _, existingEmployeeJob := range existingEmployeeJobs {
		if err := s.DB.Delete(&existingEmployeeJob).Error; err != nil {
			s.Log.Error(err)
			return errors.New("[SyncMidsuitScheduler.SyncEmployeeJob] Error when deleting employee job: " + err.Error())
		}
	}

	return nil
}

func (s *SyncMidsuitScheduler) SyncUserProfile(jwtToken string) error {
	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + "/views/user-profile"
	method := "GET"

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncEmployeeJobHistory] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncEmployeeJobHistory] Error when fetching response: " + err.Error())
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncEmployeeJobHistory] Error when fetching response: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var apiResponse UserProfileMidsuitAPIResponse
	if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncEmployeeJobHistory] Error when unmarshalling response: " + err.Error())
	}

	// Process the records
	for _, record := range apiResponse.Records {
		// Process each record as needed
		s.Log.Infof("Processing record: %+v", record)

		var employee entity.Employee
		if err := s.DB.Where("midsuit_id = ?", strconv.Itoa(record.HCEmployeeID.ID)).First(&employee).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Log.Errorf("Employee with ID %d not found", record.HCEmployeeID.ID)
				continue
			}
			s.Log.Error(err)
			return errors.New("[SyncMidsuitScheduler.SyncEmployeeJobHistory] Error when querying employee: " + err.Error())
		}

		var user entity.User
		if err := s.DB.Where("employee_id = ?", employee.ID).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Log.Errorf("User with employee_id %s not found", employee.ID)
				continue
			}
			s.Log.Error(err)
			return errors.New("[SyncMidsuitScheduler.SyncEmployeeJobHistory] Error when querying user: " + err.Error())
		}

		parsedBirthDate, err := time.Parse(time.RFC3339, record.BirthDate)
		if err != nil {
			s.Log.Error(err)
			return errors.New("[SyncMidsuitScheduler.SyncEmployeeJobHistory] Error when parsing birth date: " + err.Error())
		}

		userProfileMessageFactory := messaging.SendSyncUserProfileMessageFactory(s.Log)
		_, err = userProfileMessageFactory.Execute(&messaging.ISendSyncUserProfileMessageRequest{
			UserID:        user.ID.String(),
			Name:          record.Name,
			Age:           record.Age,
			BirthDate:     parsedBirthDate.String(),
			BirthPlace:    record.BirthPlace,
			MaritalStatus: strings.ToLower(record.MaritalStatus),
			PhoneNumber:   record.PhoneNumber,
		})
		if err != nil {
			s.Log.Error(err)
			return errors.New("[SyncMidsuitScheduler.SyncEmployeeJobHistory] Error when sending message: " + err.Error())
		}
	}

	return nil
}

func (s *SyncMidsuitScheduler) SyncGrade(jwtToken string) error {
	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + "/views/employee-grade"
	method := "GET"

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncGrade] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncGrade] Error when fetching response: " + err.Error())
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncGrade] Error when fetching response: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var apiResponse GradeMidsuitAPIResponse
	if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncGrade] Error when unmarshalling response: " + err.Error())
	}

	var gradeMidsuitIDs []string
	// Process the records
	for _, record := range apiResponse.Records {
		// Process each record as needed
		gradeMidsuitIDs = append(gradeMidsuitIDs, strconv.Itoa(record.ID))

		s.Log.Infof("Processing record: %+v", record)

		// find job level by midsuit id
		var jobLevel entity.JobLevel
		if err := s.DB.Where("midsuit_id = ?", strconv.Itoa(record.HcJobLevelID.ID)).First(&jobLevel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Log.Errorf("Job level with ID %d not found", record.HcJobLevelID.ID)
				continue
			}

			s.Log.Error(err)
			return errors.New("[SyncMidsuitScheduler.SyncGrade] Error when querying job level: " + err.Error())
		}

		grade := &entity.Grade{
			MidsuitID:  strconv.Itoa(record.ID),
			Name:       record.Name,
			JobLevelID: jobLevel.ID,
		}

		// Check if the record already exists
		var existingGrade entity.Grade
		if err := s.DB.Where("midsuit_id = ?", grade.MidsuitID).First(&existingGrade).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create the record if it doesn't exist
				if err := s.DB.Create(grade).Error; err != nil {
					s.Log.Error(err)
					return errors.New("[SyncMidsuitScheduler.SyncGrade] Error when creating record: " + err.Error())
				}
				s.Log.Infof("Created record: %+v", grade)
			} else {
				s.Log.Error(err)
				return errors.New("[SyncMidsuitScheduler.SyncGrade] Error when querying record: " + err.Error())
			}
		} else {
			// Update the record if it exists
			if err := s.DB.Model(&existingGrade).Updates(grade).Error; err != nil {
				s.Log.Error(err)
				return errors.New("[SyncMidsuitScheduler.SyncGrade] Error when updating record: " + err.Error())
			}
			s.Log.Infof("Updated record: %+v", existingGrade)
		}
	}

	// Delete grades that are not in the response
	var existingGrades []entity.Grade
	if err := s.DB.Where("midsuit_id NOT IN ?", gradeMidsuitIDs).Find(&existingGrades).Error; err != nil {
		s.Log.Error(err)
		return errors.New("[SyncMidsuitScheduler.SyncGrade] Error when querying grades: " + err.Error())
	}

	for _, existingGrade := range existingGrades {
		if err := s.DB.Delete(&existingGrade).Error; err != nil {
			s.Log.Error(err)
			return errors.New("[SyncMidsuitScheduler.SyncGrade] Error when deleting grade: " + err.Error())
		}
	}

	return nil
}
