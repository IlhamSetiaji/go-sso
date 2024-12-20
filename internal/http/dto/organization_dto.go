package dto

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/http/response"
)

func ConvertToOrganizationForMessageResponse(organizations *[]entity.Organization) *[]response.OrganizationForMessageResponse {
	var responseOrganizations []response.OrganizationForMessageResponse
	for _, organization := range *organizations {
		responseOrganizations = append(responseOrganizations, response.OrganizationForMessageResponse{
			ID:                 organization.ID,
			OrganizationTypeID: organization.OrganizationTypeID,
			Name:               organization.Name,
		})
	}
	return &responseOrganizations
}
