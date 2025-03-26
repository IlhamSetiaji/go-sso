package main

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/entity"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	viper := config.NewViper()
	log := config.NewLogrus(viper)
	db := config.NewDatabase()

	// Migrate the schema
	err := db.AutoMigrate(
		&entity.Application{},
		&entity.Role{},
		&entity.Permission{},
		&entity.RolePermission{},
		&entity.AuthToken{},
		&entity.OrganizationType{},
		&entity.Organization{},
		&entity.OrganizationLocation{},
		&entity.JobLevel{},
		&entity.OrganizationStructure{},
		&entity.Job{},
		&entity.Employee{},
		&entity.EmployeeJob{},
		&entity.User{},
		&entity.UserRole{},
		&entity.UserToken{},
		&entity.Grade{},
	)

	if err != nil {
		log.Fatalf("failed to migrate schema: %v", err)
	} else {
		log.Printf("migrate schema success")
	}

	// create organization types
	organizationTypes := []entity.OrganizationType{
		{
			Name: "Anjing",
		},
		{
			Name: "Kucing",
		},
	}

	err = db.Create(&organizationTypes).Error
	if err != nil {
		log.Fatalf("failed to create organization type: %v", err)
	} else {
		log.Printf("create organization type success")
	}

	// create organization
	organizations := []entity.Organization{
		{
			Name:               "Organization 1",
			OrganizationTypeID: organizationTypes[0].ID,
			OrganizationLocations: []entity.OrganizationLocation{
				{
					Name: "Location A1",
				},
				{
					Name: "Location A2",
				},
				{
					Name: "Location A3",
				},
			},
		},
		{
			Name:               "Organization 2",
			OrganizationTypeID: organizationTypes[1].ID,
			OrganizationLocations: []entity.OrganizationLocation{
				{
					Name: "Location B1",
				},
				{
					Name: "Location B2",
				},
				{
					Name: "Location B3",
				},
			},
		},
	}

	err = db.Create(&organizations).Error
	if err != nil {
		log.Fatalf("failed to create organization: %v", err)
	} else {
		log.Printf("create organization success")
	}

	// create job levels
	jobLevels := []entity.JobLevel{
		{
			Name:  "Kage",
			Level: "5",
		},
		{
			Name:  "Jounin",
			Level: "4",
		},
		{
			Name:  "Chunin",
			Level: "3",
		},
		{
			Name:  "Genin",
			Level: "2",
		},
		{
			Name:  "Academy",
			Level: "1",
		},
	}

	err = db.Create(&jobLevels).Error
	if err != nil {
		log.Fatalf("failed to create job level: %v", err)
	} else {
		log.Printf("create job level success")
	}

	// create organization structures
	organizationStructures := []entity.OrganizationStructure{
		{
			OrganizationID: organizations[0].ID,
			Name:           "Kage Department",
			JobLevelID:     jobLevels[0].ID,
		},
		{
			OrganizationID: organizations[0].ID,
			Name:           "Jounin Department",
			JobLevelID:     jobLevels[1].ID,
		},
		{
			OrganizationID: organizations[0].ID,
			Name:           "Chunin Department",
			JobLevelID:     jobLevels[2].ID,
		},
		{
			OrganizationID: organizations[0].ID,
			Name:           "Genin Department",
			JobLevelID:     jobLevels[3].ID,
		},
		{
			OrganizationID: organizations[0].ID,
			Name:           "Academy Department",
			JobLevelID:     jobLevels[4].ID,
		},
	}

	// Map to store created structures
	createdStructures := make(map[string]*entity.OrganizationStructure)

	// Create hierarchical structure
	for i, structure := range organizationStructures {
		currentStructure := structure // Create a copy to avoid pointer issues

		if i == 0 {
			// Create root level
			if err := db.Create(&currentStructure).Error; err != nil {
				log.Fatalf("failed to create root organization structure: %v", err)
			}
			createdStructures[currentStructure.Name] = &currentStructure
			log.Printf("created root organization structure: %s", currentStructure.Name)
		} else {
			// Set parent ID to previous level
			parentStructure := createdStructures[organizationStructures[i-1].Name]
			currentStructure.ParentID = &parentStructure.ID

			if err := db.Create(&currentStructure).Error; err != nil {
				log.Fatalf("failed to create organization structure: %v", err)
			}
			createdStructures[currentStructure.Name] = &currentStructure
			log.Printf("created organization structure: %s with parent: %s",
				currentStructure.Name, parentStructure.Name)
		}
	}

	// create jobs
	jobs := []entity.Job{
		{
			Name:                    "Kage",
			OrganizationStructureID: createdStructures["Kage Department"].ID,
			Existing:                30,
			JobLevelID:              jobLevels[0].ID,
			OrganizationID:          organizations[0].ID,
		},
		{
			Name:                    "Jounin",
			OrganizationStructureID: createdStructures["Jounin Department"].ID,
			Existing:                30,
			JobLevelID:              jobLevels[1].ID,
			OrganizationID:          organizations[0].ID,
		},
		{
			Name:                    "Chunin",
			OrganizationStructureID: createdStructures["Chunin Department"].ID,
			Existing:                30,
			JobLevelID:              jobLevels[2].ID,
			OrganizationID:          organizations[0].ID,
		},
		{
			Name:                    "Genin",
			OrganizationStructureID: createdStructures["Genin Department"].ID,
			Existing:                30,
			JobLevelID:              jobLevels[3].ID,
			OrganizationID:          organizations[0].ID,
		},
		{
			Name:                    "Academy",
			OrganizationStructureID: createdStructures["Academy Department"].ID,
			Existing:                30,
			JobLevelID:              jobLevels[4].ID,
			OrganizationID:          organizations[0].ID,
		},
	}

	createdJobs := make(map[string]*entity.Job)

	for i, job := range jobs {
		currentJob := job

		if i == 0 {
			// Create root level
			if err := db.Create(&currentJob).Error; err != nil {
				log.Fatalf("failed to create root job: %v", err)
			}
			createdJobs[currentJob.Name] = &currentJob
			log.Printf("created root job: %s", currentJob.Name)
		} else {
			// Set parent ID to previous level
			parentJob := createdJobs[jobs[i-1].Name]
			currentJob.ParentID = &parentJob.ID

			if err := db.Create(&currentJob).Error; err != nil {
				log.Fatalf("failed to create job: %v", err)
			}
			createdJobs[currentJob.Name] = &currentJob
			log.Printf("created job: %s with parent: %s",
				currentJob.Name, parentJob.Name)
		}
	}

	// create multiple applications
	applications := []entity.Application{
		{
			Name:        "authenticator",
			Label:       "Authenticator",
			Secret:      "secret for authenticator",
			RedirectURI: "http://localhost:3000",
			Domain:      "localhost",
		},
		{
			Name:        "manpower",
			Label:       "Julong Manpower Planning & Request",
			Secret:      "secret for web1",
			RedirectURI: "https://www.google.com",
			Domain:      "localhost",
		},
		{
			Name:        "recruitment",
			Label:       "Julong Recruitment",
			Secret:      "secret for web2",
			RedirectURI: "https://www.github.com",
			Domain:      "localhost",
		},
	}

	for _, application := range applications {
		err = db.Create(&application).Error
		if err != nil {
			log.Fatalf("failed to create application: %v", err)
		} else {
			log.Printf("create application success")
		}
	}

	// create default roles and permissions
	var authApplication entity.Application
	err = db.Where("name = ?", "authenticator").First(&authApplication).Error
	if err != nil {
		log.Fatalf("failed to find Application: %v", err)
	}
	role := entity.Role{
		Name:          "superadmin",
		GuardName:     "web",
		ApplicationID: authApplication.ID,
		Permissions: []entity.Permission{
			{
				Name:          "view-profile-applicant",
				Label:         "View Profile Applicant",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "update-user",
				Label:         "Update User",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "update-role",
				Label:         "Update Role",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "update-permission",
				Label:         "Update Permission",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "update-organization-structure",
				Label:         "Update Organization Structure",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "update-organization-location",
				Label:         "Update Organization Location",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "update-organization",
				Label:         "Update Organization",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "update-job-level",
				Label:         "Update Job Level",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "update-job",
				Label:         "Update Job",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "update-employee-job",
				Label:         "Update Employee Job",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "update-employee",
				Label:         "Update Employee",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "update-client",
				Label:         "Update Client",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "sync-job",
				Label:         "Sync Job",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "submit-verification-document",
				Label:         "Submit Verification Document",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "submit-result-test",
				Label:         "Submit Result Test",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "submit-mpr-staff",
				Label:         "Submit MPR Staff",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "submit-mpr-dept-head",
				Label:         "Submit MPR Dept Head",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "submit-mpp",
				Label:         "Submit MPP",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "submit-administrative-selection-setup",
				Label:         "Submit Administrative Selection Setup",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "send-offering-letter-document",
				Label:         "Send Offering Letter Document",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "send-contract-document-document",
				Label:         "Send Contract Document Document",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "reject-mpp",
				Label:         "Reject MPP",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-verification-profile",
				Label:         "Read Verification Profile",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-verification-document",
				Label:         "Read Verification Document",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-user",
				Label:         "Read User",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-university",
				Label:         "Read University",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-type-test",
				Label:         "Read Type Test",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-test-applicant-overview",
				Label:         "Read Test Applicant Overview",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-template-question",
				Label:         "Read Template Question",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-template-activity",
				Label:         "Read Template Activity",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-task-template",
				Label:         "Read Task Template",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-task",
				Label:         "Read Task",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-survey",
				Label:         "Read Survey",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-schedule-test",
				Label:         "Read Schedule Test",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-schedule-interview",
				Label:         "Read Schedule Interview",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-schedule-fgd",
				Label:         "Read Schedule FGD",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-role",
				Label:         "Read Role",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-result-test",
				Label:         "Read Result Test",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-result-interview",
				Label:         "Read Result Interview",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-result-fgd",
				Label:         "Read Result FGD",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-project-recruitment",
				Label:         "Read Project Recruitment",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-project-calender",
				Label:         "Read Project Calender",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-plafon",
				Label:         "Read Plafon",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-permission",
				Label:         "Read Permission",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-period",
				Label:         "Read Period",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-organization-structure",
				Label:         "Read Organization Structure",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-organization-location",
				Label:         "Read Organization Location",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-organization",
				Label:         "Read Organization",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-offering-letter-document",
				Label:         "Read Offering Letter Document",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-offering-letter-agreement",
				Label:         "Read Offering Letter Agreement",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-offering-letter",
				Label:         "Read Offering Letter",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-mpr-vp",
				Label:         "Read MPR VP",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-mpr-recruitment",
				Label:         "Read MPR Recruitment",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-mpr-ho",
				Label:         "Read MPR HO",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-mpr-dept-head",
				Label:         "Read MPR Dept Head",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-mpr-ceo",
				Label:         "Read MPR CEO",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-mpr",
				Label:         "Read MPR",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-mpp-rekruitmen",
				Label:         "Read MPP Rekruitmen",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-mpp-hrd-unit",
				Label:         "Read MPP HRD Unit",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-mpp-hrd-location",
				Label:         "Read MPP HRD Location",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-mpp-dir-unit",
				Label:         "Read MPP Dir Unit",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-mpp",
				Label:         "Read MPP",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-major",
				Label:         "Read Major",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-mail-template",
				Label:         "Read Mail Template",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-letterhead",
				Label:         "Read Letterhead",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-job-poting-onboarding",
				Label:         "Read Job Poting Onboarding",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-job-posting",
				Label:         "Read Job Posting",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-job-level",
				Label:         "Read Job Level",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-job",
				Label:         "Read Job",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-home-onboarding",
				Label:         "Read Home Onboarding",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-home-mpp",
				Label:         "Read Home MPP",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-final-result-interview",
				Label:         "Read Final Result Interview",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-final-result-fgd",
				Label:         "Read Final Result FGD",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-final-interview-result",
				Label:         "Read Final Interview Result",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-final-interview-recruitment",
				Label:         "Read Final Interview Recruitment",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-final-interview-calender",
				Label:         "Read Final Interview Calender",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-fgd-setup",
				Label:         "Read FGD Setup",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-fgd-schedule",
				Label:         "Read FGD Schedule",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-fgd-result",
				Label:         "Read FGD Result",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-events",
				Label:         "Read Events",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-employee-job",
				Label:         "Read Employee Job",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-employee",
				Label:         "Read Employee",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-document-setup",
				Label:         "Read Document Setup",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-document-checking",
				Label:         "Read Document Checking",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-dashboard",
				Label:         "Read Dashboard",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-contract-document-document",
				Label:         "Read Contract Document Document",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-contract-document-agreement",
				Label:         "Read Contract Document Agreement",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-contract-document",
				Label:         "Read Contract Document",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-client",
				Label:         "Read Client",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-batch",
				Label:         "Read Batch",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-applicant-result",
				Label:         "Read Applicant Result",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-applicant-document-cover-letter",
				Label:         "Read Applicant Document Cover Letter",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-applicant-document-checking",
				Label:         "Read Applicant Document Checking",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-applicant-document",
				Label:         "Read Applicant Document",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-administrative-selection-setup",
				Label:         "Read Administrative Selection Setup",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-administrative-applicant-overview",
				Label:         "Read Administrative Applicant Overview",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "read-activity",
				Label:         "Read Activity",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "process-mpp",
				Label:         "Process MPP",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "example-permission-user-2",
				Label:         "Example Permission User 2",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "example-permission-user-1",
				Label:         "Example Permission User 1",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "example-permission-admin-2",
				Label:         "Example Permission Admin 2",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "example-permission-admin-1",
				Label:         "Example Permission Admin 1",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-university",
				Label:         "Edit University",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-type-test",
				Label:         "Edit Type Test",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-template-question",
				Label:         "Edit Template Question",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-template-mail",
				Label:         "Edit Template Mail",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-template-activity",
				Label:         "Edit Template Activity",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-task-template",
				Label:         "Edit Task Template",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-task",
				Label:         "Edit Task",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-survey",
				Label:         "Edit Survey",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-schedule-test",
				Label:         "Edit Schedule Test",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-schedule-interview",
				Label:         "Edit Schedule Interview",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-schedule-fgd",
				Label:         "Edit Schedule FGD",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-project-recruitment",
				Label:         "Edit Project Recruitment",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-plafon",
				Label:         "Edit Plafon",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-period",
				Label:         "Edit Period",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-offering-letter-document",
				Label:         "Edit Offering Letter Document",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-mpr",
				Label:         "Edit MPR",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-mpp",
				Label:         "Edit MPP",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-major",
				Label:         "Edit Major",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-letterhead",
				Label:         "Edit Letterhead",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-job-posting",
				Label:         "Edit Job Posting",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-final-schedule-interview",
				Label:         "Edit Final Schedule Interview",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-final-result-interview",
				Label:         "Edit Final Result Interview",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-fgd-schedule",
				Label:         "Edit FGD Schedule",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-events",
				Label:         "Edit Events",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-document-setup",
				Label:         "Edit Document Setup",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-document-checking",
				Label:         "Edit Document Checking",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-contract-document-document",
				Label:         "Edit Contract Document Document",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-contract-document",
				Label:         "Edit Contract Document",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-applicant-overview",
				Label:         "Edit Applicant Overview",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-applicant-document-cover-letter",
				Label:         "Edit Applicant Document Cover Letter",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-applicant-document-checking",
				Label:         "Edit Applicant Document Checking",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "edit-activity",
				Label:         "Edit Activity",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-user",
				Label:         "Delete User",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-university",
				Label:         "Delete University",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-type-test",
				Label:         "Delete Type Test",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-template-question",
				Label:         "Delete Template Question",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-template-mail",
				Label:         "Delete Template Mail",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-template-activity",
				Label:         "Delete Template Activity",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-task-template",
				Label:         "Delete Task Template",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-task",
				Label:         "Delete Task",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-schedule-interview",
				Label:         "Delete Schedule Interview",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-role",
				Label:         "Delete Role",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-project-recruitment",
				Label:         "Delete Project Recruitment",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-permission",
				Label:         "Delete Permission",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-period",
				Label:         "Delete Period",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-organization-structure",
				Label:         "Delete Organization Structure",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-organization-location",
				Label:         "Delete Organization Location",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-organization",
				Label:         "Delete Organization",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-offering-letter-document",
				Label:         "Delete Offering Letter Document",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-mpr",
				Label:         "Delete MPR",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-mpp",
				Label:         "Delete MPP",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-major",
				Label:         "Delete Major",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-job-posting",
				Label:         "Delete Job Posting",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-job-level",
				Label:         "Delete Job Level",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-job",
				Label:         "Delete Job",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-final-schedule-interview",
				Label:         "Delete Final Schedule Interview",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-final-result-interview",
				Label:         "Delete Final Result Interview",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-fgd-schedule",
				Label:         "Delete FGD Schedule",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-events",
				Label:         "Delete Events",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-employee-job",
				Label:         "Delete Employee Job",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-employee",
				Label:         "Delete Employee",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-document-setup",
				Label:         "Delete Document Setup",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-document-checking",
				Label:         "Delete Document Checking",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-contract-document-document",
				Label:         "Delete Contract Document Document",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-client",
				Label:         "Delete Client",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-applicant-overview",
				Label:         "Delete Applicant Overview",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-activity",
				Label:         "Delete Activity",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-verification-document",
				Label:         "Create Verification Document",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-user",
				Label:         "Create User",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-university",
				Label:         "Create University",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-type-test",
				Label:         "Create Type Test",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-template-question",
				Label:         "Create Template Question",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-template-mail",
				Label:         "Create Template Mail",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-template-activity",
				Label:         "Create Template Activity",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-task-template",
				Label:         "Create Task Template",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-task",
				Label:         "Create Task",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-survey",
				Label:         "Create Survey",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-schedule-test",
				Label:         "Create Schedule Test",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-schedule-interview",
				Label:         "Create Schedule Interview",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-schedule-fgd",
				Label:         "Create Schedule FGD",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-role",
				Label:         "Create Role",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-result-fgd",
				Label:         "Create Result FGD",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-project-recruitment",
				Label:         "Create Project Recruitment",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-permission",
				Label:         "Create Permission",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-period",
				Label:         "Create Period",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-organization-structure",
				Label:         "Create Organization Structure",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-organization-location",
				Label:         "Create Organization Location",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-organization",
				Label:         "Create Organization",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-offering-letter-document",
				Label:         "Create Offering Letter Document",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-mpr",
				Label:         "Create MPR",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-mpp",
				Label:         "Create MPP",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-major",
				Label:         "Create Major",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-letterhead",
				Label:         "Create Letterhead",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-job-posting",
				Label:         "Create Job Posting",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-job-level",
				Label:         "Create Job Level",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-job",
				Label:         "Create Job",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-final-schedule-interview",
				Label:         "Create Final Schedule Interview",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-final-result-interview",
				Label:         "Create Final Result Interview",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-fgd-schedule",
				Label:         "Create FGD Schedule",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-events",
				Label:         "Create Events",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-employee-job",
				Label:         "Create Employee Job",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-employee",
				Label:         "Create Employee",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-document-setup",
				Label:         "Create Document Setup",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-document-checking",
				Label:         "Create Document Checking",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-contract-document-document",
				Label:         "Create Contract Document Document",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-contract-document",
				Label:         "Create Contract Document",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-client",
				Label:         "Create Client",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-batch-hrd-unit",
				Label:         "Create Batch HRD Unit",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-batch",
				Label:         "Create Batch",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-applicant-overview",
				Label:         "Create Applicant Overview",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-applicant-document-cover-letter",
				Label:         "Create Applicant Document Cover Letter",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-applicant-document-checking",
				Label:         "Create Applicant Document Checking",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "create-activity",
				Label:         "Create Activity",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "completed-offering-letter-document",
				Label:         "Completed Offering Letter Document",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "completed-offering-letter-agreement",
				Label:         "Completed Offering Letter Agreement",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "completed-contract-document-document",
				Label:         "Completed Contract Document Document",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "completed-contract-document-agreement",
				Label:         "Completed Contract Document Agreement",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "complete-batch",
				Label:         "Complete Batch",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "batch-ceo",
				Label:         "Batch CEO",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "assign-role",
				Label:         "Assign Role",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "assign-permission",
				Label:         "Assign Permission",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "approve-mpp",
				Label:         "Approve MPP",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "approval-verification-profile",
				Label:         "Approval Verification Profile",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "approval-result-fgd",
				Label:         "Approval Result FGD",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "approval-mpr-vp",
				Label:         "Approval MPR VP",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "approval-mpr-ho",
				Label:         "Approval MPR HO",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "approval-mpr-dept-head",
				Label:         "Approval MPR Dept Head",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "approval-mpr-ceo",
				Label:         "Approval MPR CEO",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "approval-mpp-direktur",
				Label:         "Approval MPP Direktur",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "approval-final-result-interview",
				Label:         "Approval Final Result Interview",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "approval-final-result-fgd",
				Label:         "Approval Final Result FGD",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "approval-contract-document-agreement",
				Label:         "Approval Contract Document Agreement",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "approval-batch-direktur",
				Label:         "Approval Batch Direktur",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "approval-batch-ceo",
				Label:         "Approval Batch CEO",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "approval-applicant-document-selection",
				Label:         "Approval Applicant Document Selection",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "approval-applicant-document-checking",
				Label:         "Approval Applicant Document Checking",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
		},
	}

	err = db.Create(&role).Error
	if err != nil {
		log.Fatalf("failed to create role: %v", err)
	} else {
		log.Printf("create role success")
	}

	// Create default user
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte("changeme"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	var superadminRole entity.Role
	err = db.Where("name = ?", "superadmin").First(&superadminRole).Error
	if err != nil {
		log.Fatalf("failed to find role: %v", err)
	}

	employeeAdmin := entity.Employee{
		OrganizationID: organizations[0].ID,
		Name:           "Employee Admin",
		Email:          "admin@test.test",
		MobilePhone:    "081234567890",
		EndDate:        time.Now().AddDate(5, 0, 0),
		RetirementDate: time.Now().AddDate(10, 0, 0),
	}

	user := entity.User{
		Username:        "admin",
		Email:           "admin@test.test",
		Password:        string(hashedPasswordBytes),
		Name:            "Admin",
		EmailVerifiedAt: time.Now(),
		Status:          entity.UserStatus("ACTIVE"),
		Gender:          entity.UserGender("MALE"),
		MobilePhone:     "081234567890",
		Employee:        &employeeAdmin,
		// Roles:           []entity.Role{role},
	}

	err = db.Create(&user).Error
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	} else {
		log.Printf("create user success")
	}

	employeeJob := &entity.EmployeeJob{
		Name:                    "Employee Admin Job",
		EmpOrganizationID:       &organizations[0].ID,
		OrganizationLocationID:  organizations[0].OrganizationLocations[0].ID,
		EmployeeID:              &user.Employee.ID,
		JobID:                   createdJobs["Kage"].ID,
		OrganizationStructureID: createdStructures["Kage Department"].ID,
	}

	err = db.Create(&employeeJob).Error
	if err != nil {
		log.Fatalf("failed to create employee job: %v", err)
	} else {
		log.Printf("create employee job success")
	}

	userRole := entity.UserRole{
		UserID: user.ID,
		RoleID: superadminRole.ID,
	}

	err = db.Create(&userRole).Error

	if err != nil {
		log.Fatalf("failed to create user role: %v", err)
	} else {
		log.Printf("create user role success")
	}

	employeeGoogle := entity.Employee{
		OrganizationID: organizations[0].ID,
		Name:           "Ilham Setiaji",
		Email:          "ilham.ahmadz18@gmail.com",
		MobilePhone:    "081234567891",
		EndDate:        time.Now().AddDate(5, 0, 0),
		RetirementDate: time.Now().AddDate(10, 0, 0),
	}

	// insert google account
	googleAccount := entity.User{
		Username:        "ilham",
		Email:           "ilham.ahmadz18@gmail.com",
		Password:        string(hashedPasswordBytes),
		Name:            "Ilham Setiaji",
		EmailVerifiedAt: time.Now(),
		Status:          entity.UserStatus("ACTIVE"),
		Gender:          entity.UserGender("MALE"),
		MobilePhone:     "081234567891",
		Employee:        &employeeGoogle,
	}

	err = db.Create(&googleAccount).Error
	if err != nil {
		log.Fatalf("failed to create google account: %v", err)
	} else {
		log.Printf("create google account success")
	}

	employeeGoogleJob := &entity.EmployeeJob{
		Name:                    "Ilham Setiaji Job",
		EmpOrganizationID:       &organizations[0].ID,
		OrganizationLocationID:  organizations[0].OrganizationLocations[0].ID,
		EmployeeID:              &googleAccount.Employee.ID,
		JobID:                   createdJobs["Kage"].ID,
		OrganizationStructureID: createdStructures["Kage Department"].ID,
	}

	err = db.Create(&employeeGoogleJob).Error
	if err != nil {
		log.Fatalf("failed to create employee google job: %v", err)
	} else {
		log.Printf("create employee google job success")
	}

	googleUserRole := entity.UserRole{
		UserID: googleAccount.ID,
		RoleID: superadminRole.ID,
	}

	err = db.Create(&googleUserRole).Error

	if err != nil {
		log.Fatalf("failed to create google user role: %v", err)
	} else {
		log.Printf("create google user role success")
	}

	// populate users and roles for web1 and web2 client
	var web1Application entity.Application
	err = db.Where("name = ?", "manpower").First(&web1Application).Error
	if err != nil {
		log.Fatalf("failed to find Application: %v", err)
	}

	var web2Application entity.Application
	err = db.Where("name = ?", "recruitment").First(&web2Application).Error
	if err != nil {
		log.Fatalf("failed to find Application: %v", err)
	}

	roles := []entity.Role{
		{
			Name:          "admin",
			GuardName:     "web",
			ApplicationID: web1Application.ID,
			Permissions: []entity.Permission{
				{
					Name:          "example-permission-admin-1",
					Label:         "Example Permission Admin 1",
					GuardName:     "web",
					ApplicationID: web1Application.ID,
				},
			},
		},
		{
			Name:          "user",
			GuardName:     "web",
			ApplicationID: web1Application.ID,
			Permissions: []entity.Permission{
				{
					Name:          "example-permission-user-1",
					Label:         "Example Permission User 1",
					GuardName:     "web",
					ApplicationID: web1Application.ID,
				},
			},
		},
		{
			Name:          "Tim Rekrutmen",
			GuardName:     "web",
			ApplicationID: web2Application.ID,
			// Permissions: []entity.Permission{
			// 	{Name: "example-permission-admin-2", Label: "Example Permission Admin 2", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization", Label: "Read Organization", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization-location", Label: "Read Organization Location", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization-structure", Label: "Read Organization Structure", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-job", Label: "Read Job", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-period", Label: "Read Period", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-period", Label: "Create Period", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-period", Label: "Edit Period", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-period", Label: "Delete Period", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-batch", Label: "Read Batch", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "complete-batch", Label: "Complete Batch", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-batch", Label: "Create Batch", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-mpr-ho", Label: "Read MPR HO", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "approval-mpr-ho", Label: "Approval MPR HO", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-user", Label: "Read User", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-user", Label: "Create User", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-user", Label: "Delete User", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "update-user", Label: "Update User", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-home-mpp", Label: "Read Home MPP", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "example-permission-user-2", Label: "Example Permission User 2", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-template-activity", Label: "Read Template Activity", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-template-activity", Label: "Create Template Activity", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-template-activity", Label: "Edit Template Activity", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-template-activity", Label: "Delete Template Activity", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-template-question", Label: "Read Template Question", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-template-question", Label: "Create Template Question", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-template-question", Label: "Edit Template Question", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-template-question", Label: "Delete Template Question", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-mail-template", Label: "Read Mail Template", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-template-mail", Label: "Create Template Mail", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-template-mail", Label: "Edit Template Mail", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-template-mail", Label: "Delete Template Mail", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-document-setup", Label: "Read Document Setup", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-document-setup", Label: "Create Document Setup", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-document-setup", Label: "Edit Document Setup", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-document-setup", Label: "Delete Document Setup", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-activity", Label: "Read Activity", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-activity", Label: "Create Activity", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-activity", Label: "Edit Activity", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-activity", Label: "Delete Activity", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-major", Label: "Read Major", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-major", Label: "Create Major", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-major", Label: "Edit Major", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-major", Label: "Delete Major", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-university", Label: "Read University", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-university", Label: "Create University", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-university", Label: "Edit University", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-university", Label: "Delete University", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-type-test", Label: "Read Type Test", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-type-test", Label: "Create Type Test", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-type-test", Label: "Edit Type Test", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-type-test", Label: "Delete Type Test", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-project-recruitment", Label: "Read Project Recruitment", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-job-posting", Label: "Create Job Posting", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-job-posting", Label: "Edit Job Posting", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-job-posting", Label: "Delete Job Posting", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-document-checking", Label: "Read Document Checking", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-document-checking", Label: "Create Document Checking", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-document-checking", Label: "Edit Document Checking", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-document-checking", Label: "Delete Document Checking", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-project-recruitment", Label: "Create Project Recruitment", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-project-recruitment", Label: "Edit Project Recruitment", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-project-recruitment", Label: "Delete Project Recruitment", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "sync-job", Label: "Sync Job", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-plafon", Label: "Read Plafon", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-plafon", Label: "Edit Plafon", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "approval-verification-profile", Label: "Approval Verification Profile", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-verification-profile", Label: "Read Verification Profile", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-administrative-applicant-overview", Label: "Read Administrative Applicant Overview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-administrative-selection-setup", Label: "Read Administrative Selection Setup", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-schedule-test", Label: "Read Schedule Test", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-result-test", Label: "Read Result Test", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-test-applicant-overview", Label: "Read Test Applicant Overview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-schedule-interview", Label: "Read Schedule Interview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-result-interview", Label: "Read Result Interview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-fgd-setup", Label: "Read FGD Setup", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-fgd-schedule", Label: "Read FGD Schedule", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-fgd-result", Label: "Read FGD Result", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-final-interview-calender", Label: "Read Final Interview Calender", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-final-interview-recruitment", Label: "Read Final Interview Recruitment", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-final-interview-result", Label: "Read Final Interview Result", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-offering-letter", Label: "Read Offering Letter", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-contract-document", Label: "Read Contract Document", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-applicant-document", Label: "Read Applicant Document", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-dashboard", Label: "Read Dashboard", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-applicant-result", Label: "Read Applicant Result", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-schedule-test", Label: "Create Schedule Test", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-schedule-test", Label: "Edit Schedule Test", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-schedule-interview", Label: "Create Schedule Interview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-schedule-interview", Label: "Edit Schedule Interview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-schedule-interview", Label: "Delete Schedule Interview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-fgd-schedule", Label: "Create FGD Schedule", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-fgd-schedule", Label: "Edit FGD Schedule", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-fgd-schedule", Label: "Delete FGD Schedule", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-final-schedule-interview", Label: "Create Final Schedule Interview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-final-schedule-interview", Label: "Edit Final Schedule Interview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-final-schedule-interview", Label: "Delete Final Schedule Interview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-final-result-interview", Label: "Create Final Result Interview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-final-result-interview", Label: "Edit Final Result Interview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-final-result-interview", Label: "Delete Final Result Interview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-applicant-overview", Label: "Create Applicant Overview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-applicant-overview", Label: "Edit Applicant Overview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-applicant-overview", Label: "Delete Applicant Overview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "approval-applicant-document-selection", Label: "Approval Applicant Document Selection", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "view-profile-applicant", Label: "View Profile Applicant", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "submit-administrative-selection-setup", Label: "Submit Administrative Selection Setup", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-offering-letter-document", Label: "Read Offering Letter Document", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-offering-letter-agreement", Label: "Read Offering Letter Agreement", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-offering-letter-document", Label: "Create Offering Letter Document", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-offering-letter-document", Label: "Edit Offering Letter Document", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-offering-letter-document", Label: "Delete Offering Letter Document", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "completed-offering-letter-agreement", Label: "Completed Offering Letter Agreement", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "completed-offering-letter-document", Label: "Completed Offering Letter Document", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "send-offering-letter-document", Label: "Send Offering Letter Document", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-contract-document-document", Label: "Read Contract Document Document", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-contract-document-agreement", Label: "Read Contract Document Agreement", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-contract-document-document", Label: "Create Contract Document Document", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-contract-document-document", Label: "Edit Contract Document Document", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-contract-document-document", Label: "Delete Contract Document Document", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "completed-contract-document-agreement", Label: "Completed Contract Document Agreement", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "completed-contract-document-document", Label: "Completed Contract Document Document", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "send-contract-document-document", Label: "Send Contract Document Document", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "submit-verification-document", Label: "Submit Verification Document", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-verification-document", Label: "Read Verification Document", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-verification-document", Label: "Create Verification Document", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-home-onboarding", Label: "Read Home Onboarding", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-events", Label: "Read Events", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-task", Label: "Read Task", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-job-poting-onboarding", Label: "Read Job Poting Onboarding", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-final-result-interview", Label: "Read Final Result Interview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "approval-final-result-interview", Label: "Approval Final Result Interview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-schedule-fgd", Label: "Read Schedule FGD", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-result-fgd", Label: "Read Result FGD", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-final-result-fgd", Label: "Read Final Result FGD", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "approval-final-result-fgd", Label: "Approval Final Result FGD", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-schedule-fgd", Label: "Create Schedule FGD", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-schedule-fgd", Label: "Edit Schedule FGD", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "approval-result-fgd", Label: "Approval Result FGD", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-result-fgd", Label: "Create Result FGD", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-task-template", Label: "Read Task Template", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-task-template", Label: "Create Task Template", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-task-template", Label: "Edit Task Template", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-task-template", Label: "Delete Task Template", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-events", Label: "Create Events", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-events", Label: "Edit Events", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-events", Label: "Delete Events", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-task", Label: "Create Task", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-task", Label: "Edit Task", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-task", Label: "Delete Task", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-job-posting", Label: "Read Job Posting", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-job-posting", Label: "Create Job Posting", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-job-posting", Label: "Edit Job Posting", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-contract-document", Label: "Create Contract Document", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-contract-document", Label: "Edit Contract Document", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "approval-contract-document-agreement", Label: "Approval Contract Document Agreement", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-applicant-document-checking", Label: "Read Applicant Document Checking", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "approval-applicant-document-checking", Label: "Approval Applicant Document Checking", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-applicant-document-cover-letter", Label: "Read Applicant Document Cover Letter", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-applicant-document-cover-letter", Label: "Edit Applicant Document Cover Letter", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-applicant-document-cover-letter", Label: "Create Applicant Document Cover Letter", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-applicant-document-checking", Label: "Create Applicant Document Checking", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-applicant-document-checking", Label: "Edit Applicant Document Checking", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "approval-applicant-document-checking", Label: "Approval Applicant Document Checking", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-project-calender", Label: "Read Project Calender", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "submit-result-test", Label: "Submit Result Test", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-mpr-recruitment", Label: "Read MPR Recruitment", GuardName: "web", ApplicationID: web2Application.ID},
			// },
		},
		{
			Name:          "Direktur Unit",
			GuardName:     "web",
			ApplicationID: web2Application.ID,
			// Permissions: []entity.Permission{
			// 	{Name: "read-user", Label: "Read User", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-role", Label: "Read Role", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-permission", Label: "Read Permission", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-client", Label: "Read Client", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "assign-role", Label: "Assign Role", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "assign-permission", Label: "Assign Permission", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization", Label: "Read Organization", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization-location", Label: "Read Organization Location", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-job-level", Label: "Read Job Level", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization-structure", Label: "Read Organization Structure", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-job", Label: "Read Job", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-employee", Label: "Read Employee", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-employee-job", Label: "Read Employee Job", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-mpp-dir-unit", Label: "Read MPP Dir Unit", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "approval-mpp-direktur", Label: "Approval MPP Direktur", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "approval-batch-direktur", Label: "Approval Batch Direktur", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "approve-mpp", Label: "Approve MPP", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "reject-mpp", Label: "Reject MPP", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "process-mpp", Label: "Process MPP", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-plafon", Label: "Read Plafon", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-plafon", Label: "Edit Plafon", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-home-mpp", Label: "Read Home MPP", GuardName: "web", ApplicationID: web2Application.ID},
			// },
		},
		{
			Name:          "HRD Site",
			GuardName:     "web",
			ApplicationID: web2Application.ID,
			Permissions: []entity.Permission{
				{
					Name:          "example-permission-user-2",
					Label:         "Example Permission User 2",
					GuardName:     "web",
					ApplicationID: web2Application.ID,
				},
			},
		},
		{
			Name:          "Staff",
			GuardName:     "web",
			ApplicationID: web2Application.ID,
			// Permissions: []entity.Permission{
			// 	{Name: "read-user", Label: "Read User", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-role", Label: "Read Role", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-permission", Label: "Read Permission", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-client", Label: "Read Client", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "assign-role", Label: "Assign Role", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "assign-permission", Label: "Assign Permission", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization", Label: "Read Organization", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization-location", Label: "Read Organization Location", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-job-level", Label: "Read Job Level", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization-structure", Label: "Read Organization Structure", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-job", Label: "Read Job", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-employee", Label: "Read Employee", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-employee-job", Label: "Read Employee Job", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-mpr", Label: "Read MPR", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-mpr", Label: "Create MPR", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-mpr", Label: "Edit MPR", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-mpr", Label: "Delete MPR", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "submit-mpr-staff", Label: "Submit MPR Staff", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-home-mpp", Label: "Read Home MPP", GuardName: "web", ApplicationID: web2Application.ID},
			// },
		},
		{
			Name:          "CEO",
			GuardName:     "web",
			ApplicationID: web2Application.ID,
			// Permissions: []entity.Permission{
			// 	{Name: "read-user", Label: "Read User", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-role", Label: "Read Role", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-permission", Label: "Read Permission", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-client", Label: "Read Client", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "assign-role", Label: "Assign Role", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "assign-permission", Label: "Assign Permission", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization", Label: "Read Organization", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization-location", Label: "Read Organization Location", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-job-level", Label: "Read Job Level", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization-structure", Label: "Read Organization Structure", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-job", Label: "Read Job", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-employee", Label: "Read Employee", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-employee-job", Label: "Read Employee Job", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "batch-ceo", Label: "Batch CEO", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "approval-batch-ceo", Label: "Approval Batch CEO", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-mpr-ceo", Label: "Read MPR CEO", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "approval-mpr-ceo", Label: "Approval MPR CEO", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-home-mpp", Label: "Read Home MPP", GuardName: "web", ApplicationID: web2Application.ID},
			// },
		},
		{
			Name:          "Interviewer",
			GuardName:     "web",
			ApplicationID: web2Application.ID,
			// Permissions: []entity.Permission{
			// 	{Name: "read-user", Label: "Read User", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-role", Label: "Read Role", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-permission", Label: "Read Permission", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-client", Label: "Read Client", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "assign-role", Label: "Assign Role", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "assign-permission", Label: "Assign Permission", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization", Label: "Read Organization", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization-location", Label: "Read Organization Location", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-job-level", Label: "Read Job Level", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization-structure", Label: "Read Organization Structure", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-job", Label: "Read Job", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-employee", Label: "Read Employee", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-employee-job", Label: "Read Employee Job", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-schedule-interview", Label: "Read Schedule Interview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-result-interview", Label: "Read Result Interview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-final-interview-calender", Label: "Read Final Interview Calendar", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-final-interview-recruitment", Label: "Read Final Interview Recruitment", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-final-interview-result", Label: "Read Final Interview Result", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-schedule-interview", Label: "Create Schedule Interview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-schedule-interview", Label: "Edit Schedule Interview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-schedule-interview", Label: "Delete Schedule Interview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-final-schedule-interview", Label: "Create Final Schedule Interview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-final-schedule-interview", Label: "Edit Final Schedule Interview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-final-schedule-interview", Label: "Delete Final Schedule Interview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-final-result-interview", Label: "Create Final Result Interview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-final-result-interview", Label: "Edit Final Result Interview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-final-result-interview", Label: "Delete Final Result Interview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-final-result-interview", Label: "Read Final Result Interview", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "approval-final-result-interview", Label: "Approval Final Result Interview", GuardName: "web", ApplicationID: web2Application.ID},
			// },
		},
		{
			Name:          "HRD Location",
			GuardName:     "web",
			ApplicationID: web2Application.ID,
			// Permissions: []entity.Permission{
			// 	{Name: "read-user", Label: "Read User", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-role", Label: "Read Role", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-permission", Label: "Read Permission", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-client", Label: "Read Client", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "assign-role", Label: "Assign Role", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "assign-permission", Label: "Assign Permission", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization", Label: "Read Organization", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization-location", Label: "Read Organization Location", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-job-level", Label: "Read Job Level", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization-structure", Label: "Read Organization Structure", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-job", Label: "Read Job", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-employee", Label: "Read Employee", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-employee-job", Label: "Read Employee Job", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "sync-job", Label: "Sync Job", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-plafon", Label: "Read Plafon", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-plafon", Label: "Edit Plafon", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-period", Label: "Read Period", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-mpp-hrd-location", Label: "Read MPP HRD Location", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-mpp", Label: "Create MPP", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-mpp", Label: "Edit MPP", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-mpp", Label: "Delete MPP", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "submit-mpp", Label: "Submit MPP", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-home-mpp", Label: "Read Home MPP", GuardName: "web", ApplicationID: web2Application.ID},
			// },
		},
		{
			Name:          "FGD Assessor",
			GuardName:     "web",
			ApplicationID: web2Application.ID,
			// Permissions: []entity.Permission{
			// 	{Name: "read-user", Label: "Read User", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-role", Label: "Read Role", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-permission", Label: "Read Permission", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-client", Label: "Read Client", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "assign-role", Label: "Assign Role", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "assign-permission", Label: "Assign Permission", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization", Label: "Read Organization", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization-location", Label: "Read Organization Location", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-job-level", Label: "Read Job Level", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization-structure", Label: "Read Organization Structure", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-job", Label: "Read Job", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-employee", Label: "Read Employee", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-employee-job", Label: "Read Employee Job", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-fgd-setup", Label: "Read FGD Setup", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-fgd-schedule", Label: "Read FGD Schedule", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-fgd-result", Label: "Read FGD Result", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-fgd-schedule", Label: "Create FGD Schedule", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-fgd-schedule", Label: "Edit FGD Schedule", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-fgd-schedule", Label: "Delete FGD Schedule", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-schedule-fgd", Label: "Read Schedule FGD", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-result-fgd", Label: "Read Result FGD", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-final-result-fgd", Label: "Read Final Result FGD", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "approval-final-result-fgd", Label: "Approval Final Result FGD", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-schedule-fgd", Label: "Create Schedule FGD", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-schedule-fgd", Label: "Edit Schedule FGD", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "approval-result-fgd", Label: "Approval Result FGD", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-result-fgd", Label: "Create Result FGD", GuardName: "web", ApplicationID: web2Application.ID},
			// },
		},
		{
			Name:          "Vice President",
			GuardName:     "web",
			ApplicationID: web2Application.ID,
			// Permissions: []entity.Permission{
			// 	{Name: "read-user", Label: "Read User", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-role", Label: "Read Role", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-permission", Label: "Read Permission", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-client", Label: "Read Client", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "assign-role", Label: "Assign Role", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "assign-permission", Label: "Assign Permission", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization", Label: "Read Organization", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization-location", Label: "Read Organization Location", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-job-level", Label: "Read Job Level", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization-structure", Label: "Read Organization Structure", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-job", Label: "Read Job", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-employee", Label: "Read Employee", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-employee-job", Label: "Read Employee Job", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-mpr-vp", Label: "Read MPR VP", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "approval-mpr-vp", Label: "Approval MPR VP", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-home-mpp", Label: "Read Home MPP", GuardName: "web", ApplicationID: web2Application.ID},
			// },
		},
		{
			Name:          "HRD Unit",
			GuardName:     "web",
			ApplicationID: web2Application.ID,
			// Permissions: []entity.Permission{
			// 	{Name: "read-user", Label: "Read User", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-role", Label: "Read Role", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-permission", Label: "Read Permission", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-client", Label: "Read Client", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "assign-role", Label: "Assign Role", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "assign-permission", Label: "Assign Permission", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization", Label: "Read Organization", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization-location", Label: "Read Organization Location", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-job-level", Label: "Read Job Level", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization-structure", Label: "Read Organization Structure", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-job", Label: "Read Job", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-employee", Label: "Read Employee", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-employee-job", Label: "Read Employee Job", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "sync-job", Label: "Sync Job", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-plafon", Label: "Read Plafon", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-plafon", Label: "Edit Plafon", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-period", Label: "Read Period", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-mpp-hrd-unit", Label: "Read MPP HRD Unit", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "process-mpp", Label: "Process MPP", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-batch-hrd-unit", Label: "Create Batch HRD Unit", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-home-mpp", Label: "Read Home MPP", GuardName: "web", ApplicationID: web2Application.ID},
			// },
		},
		{
			Name:          "Dept Head",
			GuardName:     "web",
			ApplicationID: web2Application.ID,
			// Permissions: []entity.Permission{
			// 	{Name: "read-user", Label: "Read User", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-role", Label: "Read Role", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-permission", Label: "Read Permission", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-client", Label: "Read Client", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "assign-role", Label: "Assign Role", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "assign-permission", Label: "Assign Permission", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization", Label: "Read Organization", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization-location", Label: "Read Organization Location", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-job-level", Label: "Read Job Level", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-organization-structure", Label: "Read Organization Structure", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-job", Label: "Read Job", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-employee", Label: "Read Employee", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-employee-job", Label: "Read Employee Job", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-mpr-dept-head", Label: "Read MPR Dept Head", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "create-mpr", Label: "Create MPR", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "edit-mpr", Label: "Edit MPR", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "delete-mpr", Label: "Delete MPR", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "submit-mpr-dept-head", Label: "Submit MPR Dept Head", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "approval-mpr-dept-head", Label: "Approval MPR Dept Head", GuardName: "web", ApplicationID: web2Application.ID},
			// 	{Name: "read-home-mpp", Label: "Read Home MPP", GuardName: "web", ApplicationID: web2Application.ID},
			// },
		},
	}

	user1 := entity.User{
		Username:        "timrekrutmen",
		Email:           "tr@test.test",
		Password:        string(hashedPasswordBytes),
		Name:            "Tim Rekrutmen",
		EmailVerifiedAt: time.Now(),
		Status:          entity.UserStatus("ACTIVE"),
		Gender:          entity.UserGender("MALE"),
	}

	err = db.Create(&user1).Error
	if err != nil {
		log.Fatalf("failed to create user 1: %v", err)
	} else {
		log.Printf("create user 1 success")
	}

	user2 := entity.User{
		Username:        "hrdsite",
		Email:           "hrd@test.test",
		Password:        string(hashedPasswordBytes),
		Name:            "HRD Site",
		EmailVerifiedAt: time.Now(),
		Status:          entity.UserStatus("ACTIVE"),
		Gender:          entity.UserGender("FEMALE"),
	}

	err = db.Create(&user2).Error
	if err != nil {
		log.Fatalf("failed to create user 2: %v", err)
	} else {
		log.Printf("create user 2 success")
	}

	for _, role := range roles {
		err = db.Create(&role).Error
		if err != nil {
			log.Fatalf("failed to create role (on for loop): %v", err)
		} else {
			log.Printf("create role success")
		}

		if role.Name == "admin" {
			err = db.Create(&entity.UserRole{
				UserID: user.ID,
				RoleID: role.ID,
			}).Error

			if err != nil {
				log.Fatalf("failed to create user role: %v", err)
			} else {
				log.Printf("create user role success")
			}
		}

		if role.Name == "user" {
			err = db.Create(&entity.UserRole{
				UserID: user1.ID,
				RoleID: role.ID,
			}).Error

			if err != nil {
				log.Fatalf("failed to create user role: %v", err)
			} else {
				log.Printf("create user role success")
			}

			err = db.Create(&entity.UserRole{
				UserID: user2.ID,
				RoleID: role.ID,
			}).Error

			if err != nil {
				log.Fatalf("failed to create user role: %v", err)
			} else {
				log.Printf("create user role success")
			}
		}

		if role.Name == "Dept Head" {
			permissionNames := []string{
				"read-user",
				"read-role",
				"read-permission",
				"read-client",
				"assign-role",
				"assign-permission",
				"read-organization",
				"read-organization-location",
				"read-job-level",
				"read-organization-structure",
				"read-job",
				"read-employee",
				"read-employee-job",
				"read-mpr-dept-head",
				"create-mpr",
				"edit-mpr",
				"delete-mpr",
				"submit-mpr-dept-head",
				"approval-mpr-dept-head",
				"read-home-mpp",
			}

			for _, permissionName := range permissionNames {
				permission := entity.Permission{}
				err = db.Where("name = ?", permissionName).First(&permission).Error
				if err != nil {
					log.Fatalf("failed to find permission: %v", err)
				}

				err = db.Create(&entity.RolePermission{
					RoleID:       role.ID,
					PermissionID: permission.ID,
				}).Error
				if err != nil {
					log.Fatalf("failed to create role permission: %v", err)
				} else {
					log.Printf("create role permission success")
				}
			}
		}

		if role.Name == "HRD Unit" {
			permissionNames := []string{
				"read-user",
				"read-role",
				"read-permission",
				"read-client",
				"assign-role",
				"assign-permission",
				"read-organization",
				"read-organization-location",
				"read-job-level",
				"read-organization-structure",
				"read-job",
				"read-employee",
				"read-employee-job",
				"sync-job",
				"read-plafon",
				"edit-plafon",
				"read-period",
				"read-mpp-hrd-unit",
				"process-mpp",
				"create-batch-hrd-unit",
				"read-home-mpp",
			}

			for _, permissionName := range permissionNames {
				permission := entity.Permission{}
				err = db.Where("name = ?", permissionName).First(&permission).Error
				if err != nil {
					log.Fatalf("failed to find permission: %v", err)
				}

				err = db.Create(&entity.RolePermission{
					RoleID:       role.ID,
					PermissionID: permission.ID,
				}).Error
				if err != nil {
					log.Fatalf("failed to create role permission: %v", err)
				} else {
					log.Printf("create role permission success")
				}
			}
		}

		if role.Name == "Vice President" {
			permissionNames := []string{
				"read-user",
				"read-role",
				"read-permission",
				"read-client",
				"assign-role",
				"assign-permission",
				"read-organization",
				"read-organization-location",
				"read-job-level",
				"read-organization-structure",
				"read-job",
				"read-employee",
				"read-employee-job",
				"read-mpr-vp",
				"approval-mpr-vp",
				"read-home-mpp",
			}

			for _, permissionName := range permissionNames {
				permission := entity.Permission{}
				err = db.Where("name = ?", permissionName).First(&permission).Error
				if err != nil {
					log.Fatalf("failed to find permission: %v", err)
				}

				err = db.Create(&entity.RolePermission{
					RoleID:       role.ID,
					PermissionID: permission.ID,
				}).Error
				if err != nil {
					log.Fatalf("failed to create role permission: %v", err)
				} else {
					log.Printf("create role permission success")
				}
			}
		}

		if role.Name == "FGD Assessor" {
			permissionNames := []string{
				"read-user",
				"read-role",
				"read-permission",
				"read-client",
				"assign-role",
				"assign-permission",
				"read-organization",
				"read-organization-location",
				"read-job-level",
				"read-organization-structure",
				"read-job",
				"read-employee",
				"read-employee-job",
				"read-fgd-setup",
				"read-fgd-schedule",
				"read-fgd-result",
				"create-fgd-schedule",
				"edit-fgd-schedule",
				"delete-fgd-schedule",
				"read-schedule-fgd",
				"read-result-fgd",
				"read-final-result-fgd",
				"approval-final-result-fgd",
				"create-schedule-fgd",
				"edit-schedule-fgd",
				"approval-result-fgd",
				"create-result-fgd",
			}

			for _, permissionName := range permissionNames {
				permission := entity.Permission{}
				err = db.Where("name = ?", permissionName).First(&permission).Error
				if err != nil {
					log.Fatalf("failed to find permission: %v", err)
				}

				err = db.Create(&entity.RolePermission{
					RoleID:       role.ID,
					PermissionID: permission.ID,
				}).Error
				if err != nil {
					log.Fatalf("failed to create role permission: %v", err)
				} else {
					log.Printf("create role permission success")
				}
			}
		}

		if role.Name == "HRD Location" {
			permissionNames := []string{
				"read-user",
				"read-role",
				"read-permission",
				"read-client",
				"assign-role",
				"assign-permission",
				"read-organization",
				"read-organization-location",
				"read-job-level",
				"read-organization-structure",
				"read-job",
				"read-employee",
				"read-employee-job",
				"sync-job",
				"read-plafon",
				"edit-plafon",
				"read-period",
				"read-mpp-hrd-location",
				"create-mpp",
				"edit-mpp",
				"delete-mpp",
				"submit-mpp",
				"read-home-mpp",
			}

			for _, permissionName := range permissionNames {
				permission := entity.Permission{}
				err = db.Where("name = ?", permissionName).First(&permission).Error
				if err != nil {
					log.Fatalf("failed to find permission: %v", err)
				}

				err = db.Create(&entity.RolePermission{
					RoleID:       role.ID,
					PermissionID: permission.ID,
				}).Error
				if err != nil {
					log.Fatalf("failed to create role permission: %v", err)
				} else {
					log.Printf("create role permission success")
				}
			}
		}

		if role.Name == "Interviewer" {
			permissionNames := []string{
				"read-user",
				"read-role",
				"read-permission",
				"read-client",
				"assign-role",
				"assign-permission",
				"read-organization",
				"read-organization-location",
				"read-job-level",
				"read-organization-structure",
				"read-job",
				"read-employee",
				"read-employee-job",
				"read-schedule-interview",
				"read-result-interview",
				"read-final-interview-calender",
				"read-final-interview-recruitment",
				"read-final-interview-result",
				"create-schedule-interview",
				"edit-schedule-interview",
				"delete-schedule-interview",
				"create-final-schedule-interview",
				"edit-final-schedule-interview",
				"delete-final-schedule-interview",
				"create-final-result-interview",
				"edit-final-result-interview",
				"delete-final-result-interview",
				"read-final-result-interview",
				"approval-final-result-interview",
			}

			for _, permissionName := range permissionNames {
				permission := entity.Permission{}
				err = db.Where("name = ?", permissionName).First(&permission).Error
				if err != nil {
					log.Fatalf("failed to find permission: %v", err)
				}

				err = db.Create(&entity.RolePermission{
					RoleID:       role.ID,
					PermissionID: permission.ID,
				}).Error
				if err != nil {
					log.Fatalf("failed to create role permission: %v", err)
				} else {
					log.Printf("create role permission success")
				}
			}
		}

		if role.Name == "CEO" {
			permissionNames := []string{
				"read-user",
				"read-role",
				"read-permission",
				"read-client",
				"assign-role",
				"assign-permission",
				"read-organization",
				"read-organization-location",
				"read-job-level",
				"read-organization-structure",
				"read-job",
				"read-employee",
				"read-employee-job",
				"batch-ceo",
				"approval-batch-ceo",
				"read-mpr-ceo",
				"approval-mpr-ceo",
				"read-home-mpp",
			}

			for _, permissionName := range permissionNames {
				permission := entity.Permission{}
				err = db.Where("name = ?", permissionName).First(&permission).Error
				if err != nil {
					log.Fatalf("failed to find permission: %v", err)
				}

				err = db.Create(&entity.RolePermission{
					RoleID:       role.ID,
					PermissionID: permission.ID,
				}).Error
				if err != nil {
					log.Fatalf("failed to create role permission: %v", err)
				} else {
					log.Printf("create role permission success")
				}
			}
		}

		if role.Name == "Staff" {
			permissionsNames := []string{
				"read-user",
				"read-role",
				"read-permission",
				"read-client",
				"assign-role",
				"assign-permission",
				"read-organization",
				"read-organization-location",
				"read-job-level",
				"read-organization-structure",
				"read-job",
				"read-employee",
				"read-employee-job",
				"read-mpr",
				"create-mpr",
				"edit-mpr",
				"delete-mpr",
				"submit-mpr-staff",
				"read-home-mpp",
			}

			for _, permissionName := range permissionsNames {
				permission := entity.Permission{}
				err = db.Where("name = ?", permissionName).First(&permission).Error
				if err != nil {
					log.Fatalf("failed to find permission: %v", err)
				}

				err = db.Create(&entity.RolePermission{
					RoleID:       role.ID,
					PermissionID: permission.ID,
				}).Error
				if err != nil {
					log.Fatalf("failed to create role permission: %v", err)
				} else {
					log.Printf("create role permission success")
				}
			}
		}

		if role.Name == "Direktur Unit" {
			permissionNames := []string{
				"read-user",
				"read-role",
				"read-permission",
				"read-client",
				"assign-role",
				"assign-permission",
				"read-organization",
				"read-organization-location",
				"read-job-level",
				"read-organization-structure",
				"read-job",
				"read-employee",
				"read-employee-job",
				"read-mpp-dir-unit",
				"approval-mpp-direktur",
				"approval-batch-direktur",
				"approve-mpp",
				"reject-mpp",
				"process-mpp",
				"read-plafon",
				"edit-plafon",
				"read-home-mpp",
			}

			for _, permissionName := range permissionNames {
				permission := entity.Permission{}
				err = db.Where("name = ?", permissionName).First(&permission).Error
				if err != nil {
					log.Fatalf("failed to find permission: %v", err)
				}

				err = db.Create(&entity.RolePermission{
					RoleID:       role.ID,
					PermissionID: permission.ID,
				}).Error
				if err != nil {
					log.Fatalf("failed to create role permission: %v", err)
				} else {
					log.Printf("create role permission success")
				}
			}
		}

		if role.Name == "Tim Rekrutmen" {
			err = db.Create(&entity.UserRole{
				UserID: user1.ID,
				RoleID: role.ID,
			}).Error

			if err != nil {
				log.Fatalf("failed to create user role: %v", err)
			} else {
				log.Printf("create user role success")
			}

			permissionNames := []string{
				"example-permission-admin-2",
				"read-organization",
				"read-organization-location",
				"read-organization-structure",
				"read-job",
				"read-period",
				"create-period",
				"edit-period",
				"delete-period",
				"read-batch",
				"complete-batch",
				"create-batch",
				"read-mpr-ho",
				"approval-mpr-ho",
				"read-user",
				"create-user",
				"delete-user",
				"update-user",
				"read-home-mpp",
				"example-permission-user-2",
				"read-template-activity",
				"create-template-activity",
				"edit-template-activity",
				"delete-template-activity",
				"read-template-question",
				"create-template-question",
				"edit-template-question",
				"delete-template-question",
				"read-mail-template",
				"create-template-mail",
				"edit-template-mail",
				"delete-template-mail",
				"read-document-setup",
				"create-document-setup",
				"edit-document-setup",
				"delete-document-setup",
				"read-activity",
				"create-activity",
				"edit-activity",
				"delete-activity",
				"read-major",
				"create-major",
				"edit-major",
				"delete-major",
				"read-university",
				"create-university",
				"edit-university",
				"delete-university",
				"read-type-test",
				"create-type-test",
				"edit-type-test",
				"delete-type-test",
				"read-project-recruitment",
				"create-job-posting",
				"edit-job-posting",
				"delete-job-posting",
				"read-document-checking",
				"create-document-checking",
				"edit-document-checking",
				"delete-document-checking",
				"create-project-recruitment",
				"edit-project-recruitment",
				"delete-project-recruitment",
				"sync-job",
				"read-plafon",
				"edit-plafon",
				"approval-verification-profile",
				"read-verification-profile",
				"read-administrative-applicant-overview",
				"read-administrative-selection-setup",
				"read-schedule-test",
				"read-result-test",
				"read-test-applicant-overview",
				"read-schedule-interview",
				"read-result-interview",
				"read-fgd-setup",
				"read-fgd-schedule",
				"read-fgd-result",
				"read-final-interview-calender",
				"read-final-interview-recruitment",
				"read-final-interview-result",
				"read-offering-letter",
				"read-contract-document",
				"read-applicant-document",
				"read-dashboard",
				"read-applicant-result",
				"create-schedule-test",
				"edit-schedule-test",
				"create-schedule-interview",
				"edit-schedule-interview",
				"delete-schedule-interview",
				"create-fgd-schedule",
				"edit-fgd-schedule",
				"delete-fgd-schedule",
				"create-final-schedule-interview",
				"edit-final-schedule-interview",
				"delete-final-schedule-interview",
				"create-final-result-interview",
				"edit-final-result-interview",
				"delete-final-result-interview",
				"create-applicant-overview",
				"edit-applicant-overview",
				"delete-applicant-overview",
				"approval-applicant-document-selection",
				"view-profile-applicant",
				"submit-administrative-selection-setup",
				"read-offering-letter-document",
				"read-offering-letter-agreement",
				"create-offering-letter-document",
				"edit-offering-letter-document",
				"delete-offering-letter-document",
				"completed-offering-letter-agreement",
				"completed-offering-letter-document",
				"send-offering-letter-document",
				"read-contract-document-document",
				"read-contract-document-agreement",
				"create-contract-document-document",
				"edit-contract-document-document",
				"delete-contract-document-document",
				"completed-contract-document-agreement",
				"completed-contract-document-document",
				"send-contract-document-document",
				"submit-verification-document",
				"read-verification-document",
				"create-verification-document",
				"read-home-onboarding",
				"read-events",
				"read-task",
				"read-job-poting-onboarding",
				"read-final-result-interview",
				"approval-final-result-interview",
				"read-schedule-fgd",
				"read-result-fgd",
				"read-final-result-fgd",
				"approval-final-result-fgd",
				"create-schedule-fgd",
				"edit-schedule-fgd",
				"approval-result-fgd",
				"create-result-fgd",
				"read-task-template",
				"create-task-template",
				"edit-task-template",
				"delete-task-template",
				"create-events",
				"edit-events",
				"delete-events",
				"create-task",
				"edit-task",
				"delete-task",
				"read-job-posting",
				"create-contract-document",
				"edit-contract-document",
				"approval-contract-document-agreement",
				"read-applicant-document-checking",
				"approval-applicant-document-checking",
				"read-applicant-document-cover-letter",
				"edit-applicant-document-cover-letter",
				"create-applicant-document-cover-letter",
				"create-applicant-document-checking",
				"edit-applicant-document-checking",
				"read-project-calender",
				"submit-result-test",
				"read-mpr-recruitment",
			}

			for _, permissionName := range permissionNames {
				permission := entity.Permission{}
				err = db.Where("name = ?", permissionName).First(&permission).Error
				if err != nil {
					log.Fatalf("failed to get permission: %v", err)
				}

				err = db.Create(&entity.RolePermission{
					RoleID:       role.ID,
					PermissionID: permission.ID,
				}).Error
				if err != nil {
					log.Fatalf("failed to create role permission tim rekrutmen: %v", err)
				}
			}
		}

		if role.Name == "HRD Site" {
			err = db.Create(&entity.UserRole{
				UserID: user2.ID,
				RoleID: role.ID,
			}).Error

			if err != nil {
				log.Fatalf("failed to create user role: %v", err)
			} else {
				log.Printf("create user role success")
			}
		}
	}
}
