package middleware

import (
	"app/go-sso/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, err := utils.GetUserRoles(c)
		if err != nil || roles == nil {
			c.String(http.StatusForbidden, "You don't have permission to access this resource")
			c.Abort()
			return
		}
		for _, role := range roles {
			if role.Name == requiredRole {
				c.Next()
				return
			}
		}
		c.String(http.StatusForbidden, "You don't have permission to access this resource")
		c.Abort()
	}
}

func ManyRolesMiddleware(requiredRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, err := utils.GetUserRoles(c)
		if err != nil || roles == nil {
			c.String(http.StatusForbidden, "You don't have permission to access this resource")
			c.Abort()
			return
		}
		for _, role := range roles {
			for _, requiredRole := range requiredRoles {
				if role.Name == requiredRole {
					c.Next()
					return
				}
			}
		}
		c.String(http.StatusForbidden, "You don't have permission to access this resource")
		c.Abort()
	}
}

func ExceptRoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, err := utils.GetUserRoles(c)
		if err != nil || roles == nil {
			c.String(http.StatusForbidden, "You don't have permission to access this resource")
			c.Abort()
			return
		}
		for _, role := range roles {
			if role.Name == requiredRole {
				c.String(http.StatusForbidden, "You don't have permission to access this resource")
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
