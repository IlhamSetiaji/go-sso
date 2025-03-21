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
	err := db.AutoMigrate(&entity.Application{}, &entity.Role{}, &entity.Permission{}, &entity.RolePermission{}, &entity.AuthToken{}, &entity.OrganizationType{}, &entity.Organization{},
		&entity.OrganizationLocation{}, &entity.JobLevel{}, &entity.OrganizationStructure{}, &entity.Job{}, &entity.Employee{}, &entity.EmployeeJob{}, &entity.User{}, &entity.UserRole{}, &entity.UserToken{})

	if err != nil {
		log.Fatalf("failed to migrate schema: %v", err)
	} else {
		log.Printf("migrate schema success")
	}

	// create organization types
	organizationTypes := []entity.OrganizationType{
		{
			Name: "Office",
		},
		{
			Name: "Gatau Ntar Isi Sendiri",
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
				Name:          "create-user",
				Label:         "Create User",
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
				Name:          "delete-user",
				Label:         "Delete User",
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
				Name:          "create-role",
				Label:         "Create Role",
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
				Name:          "delete-role",
				Label:         "Delete Role",
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
				Name:          "create-permission",
				Label:         "Create Permission",
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
				Name:          "delete-permission",
				Label:         "Delete Permission",
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
				Name:          "create-client",
				Label:         "Create Client",
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
				Name:          "delete-client",
				Label:         "Delete Client",
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
				Name:          "read-organization",
				Label:         "Read Organization",
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
				Name:          "update-organization",
				Label:         "Update Organization",
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
				Name:          "read-organization-location",
				Label:         "Read Organization Location",
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
				Name:          "update-organization-location",
				Label:         "Update Organization Location",
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
				Name:          "read-job-level",
				Label:         "Read Job Level",
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
				Name:          "update-job-level",
				Label:         "Update Job Level",
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
				Name:          "read-organization-structure",
				Label:         "Read Organization Structure",
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
				Name:          "update-organization-structure",
				Label:         "Update Organization Structure",
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
				Name:          "read-job",
				Label:         "Read Job",
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
				Name:          "update-job",
				Label:         "Update Job",
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
				Name:          "read-employee",
				Label:         "Read Employee",
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
				Name:          "update-employee",
				Label:         "Update Employee",
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
				Name:          "read-employee-job",
				Label:         "Read Employee Job",
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
				Name:          "update-employee-job",
				Label:         "Update Employee Job",
				GuardName:     "web",
				ApplicationID: authApplication.ID,
			},
			{
				Name:          "delete-employee-job",
				Label:         "Delete Employee Job",
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
			Permissions: []entity.Permission{
				{
					Name:          "example-permission-admin-2",
					Label:         "Example Permission Admin 2",
					GuardName:     "web",
					ApplicationID: web2Application.ID,
				},
			},
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
