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
	err := db.AutoMigrate(&entity.Application{}, &entity.Role{}, &entity.Permission{}, &entity.RolePermission{}, &entity.AuthToken{}, &entity.Organization{},
		&entity.OrganizationLocation{}, &entity.JobLevel{}, &entity.OrganizationStructure{}, &entity.Job{}, &entity.Employee{}, &entity.EmployeeJob{}, &entity.User{}, &entity.UserRole{})

	if err != nil {
		log.Fatalf("failed to migrate schema: %v", err)
	} else {
		log.Printf("migrate schema success")
	}

	// create organization
	organizations := []entity.Organization{
		{
			Name: "Organization 1",
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
			Name: "Organization 2",
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

	// create multiple applications
	applications := []entity.Application{
		{
			Name:        "Authenticator",
			Secret:      "secret for authenticator",
			RedirectURI: "http://localhost:8080/auth/callback",
		},
		{
			Name:        "Web1 Application",
			Secret:      "secret for web1",
			RedirectURI: "http://localhost:8080/auth/callback",
		},
		{
			Name:        "Web2 Application",
			Secret:      "secret for web2",
			RedirectURI: "http://localhost:8080/auth/callback",
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
	err = db.Where("name = ?", "Authenticator").First(&authApplication).Error
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

	user := entity.User{
		Username:        "admin",
		Email:           "admin@test.test",
		Password:        string(hashedPasswordBytes),
		Name:            "Admin",
		EmailVerifiedAt: time.Now(),
		Status:          entity.UserStatus("ACTIVE"),
		Gender:          entity.UserGender("MALE"),
		// Roles:           []entity.Role{role},
	}

	err = db.Create(&user).Error
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	} else {
		log.Printf("create user success")
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

	// insert google account
	googleAccount := entity.User{
		Username:        "ilham",
		Email:           "ilham.ahmadz18@gmail.com",
		Password:        string(hashedPasswordBytes),
		Name:            "Ilham Setiaji",
		EmailVerifiedAt: time.Now(),
		Status:          entity.UserStatus("ACTIVE"),
		Gender:          entity.UserGender("MALE"),
	}

	err = db.Create(&googleAccount).Error
	if err != nil {
		log.Fatalf("failed to create google account: %v", err)
	} else {
		log.Printf("create google account success")
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
	err = db.Where("name = ?", "Web1 Application").First(&web1Application).Error
	if err != nil {
		log.Fatalf("failed to find Application: %v", err)
	}

	var web2Application entity.Application
	err = db.Where("name = ?", "Web2 Application").First(&web2Application).Error
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
			Name:          "admin",
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
			Name:          "user",
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
		Username:        "user1",
		Email:           "user1@test.test",
		Password:        string(hashedPasswordBytes),
		Name:            "User 1",
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
		Username:        "user2",
		Email:           "user2@test.test",
		Password:        string(hashedPasswordBytes),
		Name:            "User 2",
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

	user3 := entity.User{
		Username:        "user3",
		Email:           "user3@test.test",
		Password:        string(hashedPasswordBytes),
		Name:            "User 3",
		EmailVerifiedAt: time.Now(),
		Status:          entity.UserStatus("ACTIVE"),
		Gender:          entity.UserGender("FEMALE"),
	}

	err = db.Create(&user3).Error
	if err != nil {
		log.Fatalf("failed to create user 3: %v", err)
	} else {
		log.Printf("create user 3 success")
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
		} else {
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

			err = db.Create(&entity.UserRole{
				UserID: user3.ID,
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
