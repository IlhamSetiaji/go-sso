package middleware

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/entity"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetApiLoggedInUser(ctx *gin.Context) (*entity.User, error) {
	db := config.NewDatabase()
	user, err := GetUser(ctx)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("User not found")
	}

	var userEntity entity.User
	userID, ok := user["id"].(string)
	if !ok {
		return nil, errors.New("Invalid user ID")
	}

	fmt.Print(userID)
	if err := db.Preload("Roles.Permissions").First(&userEntity, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &userEntity, nil
}

func getApiUserRoles(ctx *gin.Context) ([]entity.Role, error) {
	user, err := GetApiLoggedInUser(ctx)
	if err != nil || user == nil {
		return nil, err
	}
	return user.Roles, nil
}

func getApiUserPermissions(ctx *gin.Context) ([]entity.Permission, error) {
	roles, err := getApiUserRoles(ctx)
	if err != nil {
		return nil, err
	}
	permissions := make([]entity.Permission, 0)
	for _, role := range roles {
		permissions = append(permissions, role.Permissions...)
	}
	return permissions, nil
}

// func PermissionApiMiddleware(requiredPermission string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		permissions, err := getApiUserPermissions(c)
// 		if err != nil || permissions == nil {
// 			c.String(http.StatusForbidden, "You don't have permission to access this resource")
// 			c.Abort()
// 			return
// 		}
// 		for _, permission := range permissions {
// 			if permission.Name == requiredPermission {
// 				c.Next()
// 				return
// 			}
// 		}
// 		c.String(http.StatusForbidden, "You don't have permission to access this resource")
// 		c.Abort()
// 	}
// }

func PermissionApiMiddleware(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions, err := getApiUserPermissions(c)
		if err != nil || permissions == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this resource"})
			c.Set("permission_denied", true)
			return
		}
		for _, permission := range permissions {
			if permission.Name == requiredPermission {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this resource"})
		c.Set("permission_denied", true)
	}
}
