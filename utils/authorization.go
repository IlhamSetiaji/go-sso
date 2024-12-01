package utils

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/entity"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetLoggedInUser(ctx *gin.Context) (*entity.User, error) {
	db := config.NewDatabase()
	session := sessions.Default(ctx)
	profile := session.Get("profile")
	if profile == nil {
		return nil, nil
	}
	userProfile, ok := profile.(entity.Profile)
	if !ok {
		return nil, nil
	}
	var user entity.User
	if err := db.Preload("Roles.Permissions").First(&user, userProfile.ID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserRoles(ctx *gin.Context) ([]entity.Role, error) {
	user, err := GetLoggedInUser(ctx)
	if err != nil || user == nil {
		return nil, err
	}
	return user.Roles, nil
}

func GetUserPermissions(ctx *gin.Context) ([]entity.Permission, error) {
	roles, err := GetUserRoles(ctx)
	if err != nil {
		return nil, err
	}
	permissions := make([]entity.Permission, 0)
	for _, role := range roles {
		permissions = append(permissions, role.Permissions...)
	}
	return permissions, nil
}
