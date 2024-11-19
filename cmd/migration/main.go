package main

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/entity"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	viper := config.NewViper()
	log := config.NewLogger(viper)
	db := config.NewDatabase()

	// Migrate the schema
	err := db.AutoMigrate(&entity.Client{}, &entity.Role{}, &entity.User{}, &entity.Permission{}, &entity.RolePermission{}, &entity.UserRole{})

	if err != nil {
		log.Fatalf("failed to migrate schema: %v", err)
	} else {
		log.Printf("migrate schema success")
	}

	// create multiple clients
	clients := []entity.Client{
		{
			Name:        "Authenticator",
			Secret:      "secret for authenticator",
			RedirectURI: "http://localhost:8080/auth/callback",
		},
		{
			Name:        "Web1 Client",
			Secret:      "secret for web1",
			RedirectURI: "http://localhost:8080/auth/callback",
		},
		{
			Name:        "Web2 Client",
			Secret:      "secret for web2",
			RedirectURI: "http://localhost:8080/auth/callback",
		},
	}

	for _, client := range clients {
		err = db.Create(&client).Error
		if err != nil {
			log.Fatalf("failed to create client: %v", err)
		} else {
			log.Printf("create client success")
		}
	}

	// create default roles and permissions
	var authClient entity.Client
	err = db.Where("name = ?", "Authenticator").First(&authClient).Error
	if err != nil {
		log.Fatalf("failed to find client: %v", err)
	}
	role := entity.Role{
		Name:      "superadmin",
		GuardName: "web",
		ClientID:  authClient.ID,
		Permissions: []entity.Permission{
			{
				Name:      "create-user",
				Label:     "Create User",
				GuardName: "web",
			},
			{
				Name:      "update-user",
				Label:     "Update User",
				GuardName: "web",
			},
			{
				Name:      "delete-user",
				Label:     "Delete User",
				GuardName: "web",
			},
			{
				Name:      "read-user",
				Label:     "Read User",
				GuardName: "web",
			},
			{
				Name:      "create-role",
				Label:     "Create Role",
				GuardName: "web",
			},
			{
				Name:      "update-role",
				Label:     "Update Role",
				GuardName: "web",
			},
			{
				Name:      "delete-role",
				Label:     "Delete Role",
				GuardName: "web",
			},
			{
				Name:      "read-role",
				Label:     "Read Role",
				GuardName: "web",
			},
			{
				Name:      "create-permission",
				Label:     "Create Permission",
				GuardName: "web",
			},
			{
				Name:      "update-permission",
				Label:     "Update Permission",
				GuardName: "web",
			},
			{
				Name:      "delete-permission",
				Label:     "Delete Permission",
				GuardName: "web",
			},
			{
				Name:      "read-permission",
				Label:     "Read Permission",
				GuardName: "web",
			},
			{
				Name:      "create-client",
				Label:     "Create Client",
				GuardName: "web",
			},
			{
				Name:      "update-client",
				Label:     "Update Client",
				GuardName: "web",
			},
			{
				Name:      "delete-client",
				Label:     "Delete Client",
				GuardName: "web",
			},
			{
				Name:      "read-client",
				Label:     "Read Client",
				GuardName: "web",
			},
			{
				Name:      "assign-role",
				Label:     "Assign Role",
				GuardName: "web",
			},
			{
				Name:      "assign-permission",
				Label:     "Assign Permission",
				GuardName: "web",
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

	// populate users and roles for web1 and web2 client
	var web1Client entity.Client
	err = db.Where("name = ?", "Web1 Client").First(&web1Client).Error
	if err != nil {
		log.Fatalf("failed to find client: %v", err)
	}

	var web2Client entity.Client
	err = db.Where("name = ?", "Web2 Client").First(&web2Client).Error
	if err != nil {
		log.Fatalf("failed to find client: %v", err)
	}

	roles := []entity.Role{
		{
			Name:      "admin",
			GuardName: "web",
			ClientID:  web1Client.ID,
			Permissions: []entity.Permission{
				{
					Name:      "example-permission-admin-1",
					Label:     "Example Permission Admin 1",
					GuardName: "web",
				},
			},
		},
		{
			Name:      "user",
			GuardName: "web",
			ClientID:  web1Client.ID,
			Permissions: []entity.Permission{
				{
					Name:      "example-permission-user-1",
					Label:     "Example Permission User 1",
					GuardName: "web",
				},
			},
		},
		{
			Name:      "admin",
			GuardName: "web",
			ClientID:  web2Client.ID,
			Permissions: []entity.Permission{
				{
					Name:      "example-permission-admin-2",
					Label:     "Example Permission Admin 2",
					GuardName: "web",
				},
			},
		},
		{
			Name:      "user",
			GuardName: "web",
			ClientID:  web2Client.ID,
			Permissions: []entity.Permission{
				{
					Name:      "example-permission-user-2",
					Label:     "Example Permission User 2",
					GuardName: "web",
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
			log.Fatalf("failed to create role: %v", err)
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
