package utils

import (
	"github.com/gin-gonic/gin"
)

type TemplateHelper struct {
	Ctx *gin.Context
}

func NewTemplateHelper(c *gin.Context) *TemplateHelper {
	return &TemplateHelper{
		Ctx: c,
	}
}

func (h *TemplateHelper) HasPermission(requiredPermission string) bool {
	permissions, err := GetUserPermissions(h.Ctx)
	if err != nil || permissions == nil {
		return false
	}

	for _, permission := range permissions {
		if permission.Name == requiredPermission {
			return true
		}
	}
	return false
}
