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
	db := config.NewDatabase(log)

	// Migrate the schema
	err := db.AutoMigrate(&entity.Role{}, &entity.User{}, &entity.Permission{}, &entity.RolePermission{}, &entity.UserRole{})

	if err != nil {
		log.Fatalf("failed to migrate schema: %v", err)
	} else {
		log.Printf("migrate schema success")
	}

	// create default roles and permissions
	role := entity.Role{
		Name:      "superadmin",
		GuardName: "web",
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
}
