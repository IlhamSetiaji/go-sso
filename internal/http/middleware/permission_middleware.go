package middleware

import (
	"app/go-sso/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PermissionMiddleware(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions, err := utils.GetUserPermissions(c)
		if err != nil || permissions == nil {
			c.String(http.StatusForbidden, "You don't have permission to access this resource")
			c.Abort()
			return
		}
		for _, permission := range permissions {
			if permission.Name == requiredPermission {
				c.Next()
				return
			}
		}
		c.String(http.StatusForbidden, "You don't have permission to access this resource")
		c.Abort()
	}
}

func ManyPermissionMiddleware(requiredPermissions []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions, err := utils.GetUserPermissions(c)
		if err != nil || permissions == nil {
			c.String(http.StatusForbidden, "You don't have permission to access this resource")
			c.Abort()
			return
		}
		for _, permission := range permissions {
			for _, requiredPermission := range requiredPermissions {
				if permission.Name == requiredPermission {
					c.Next()
					return
				}
			}
		}
		c.String(http.StatusForbidden, "You don't have permission to access this resource")
		c.Abort()
	}
}
