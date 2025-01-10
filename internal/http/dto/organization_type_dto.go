package dto

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/http/response"
)

func ConvertToOrganizationTypeResponse(organizationType *entity.OrganizationType) *response.OrganizationTypeResponse {
	return &response.OrganizationTypeResponse{
		ID:   organizationType.ID,
		Name: organizationType.Name,
	}
}
