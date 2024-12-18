package dto

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/http/response"
)

func ConvertToOrganizationLocationResponse(orgLocations *[]entity.OrganizationLocation) *[]response.OrganizationLocationResponse {
	var responseLocations []response.OrganizationLocationResponse
	for _, orgLocation := range *orgLocations {
		responseLocations = append(responseLocations, response.OrganizationLocationResponse{
			ID:               orgLocation.ID,
			OrganizationID:   orgLocation.OrganizationID,
			OrganizationName: orgLocation.Organization.Name,
			Name:             orgLocation.Name,
			CreatedAt:        orgLocation.CreatedAt,
			UpdatedAt:        orgLocation.UpdatedAt,
		})
	}
	return &responseLocations
}

func ConvertToSingleOrganizationLocationResponse(orgLocation *entity.OrganizationLocation) *response.OrganizationLocationResponse {
	return &response.OrganizationLocationResponse{
		ID:             orgLocation.ID,
		OrganizationID: orgLocation.OrganizationID,
		Name:           orgLocation.Name,
		CreatedAt:      orgLocation.CreatedAt,
		UpdatedAt:      orgLocation.UpdatedAt,
	}
}
