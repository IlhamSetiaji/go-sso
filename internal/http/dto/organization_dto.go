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

func ConvertToOrganizationResponse(organizations *[]entity.Organization) *[]response.OrganizationResponse {
	var responseOrganizations []response.OrganizationResponse
	for _, organization := range *organizations {
		responseOrganizations = append(responseOrganizations, response.OrganizationResponse{
			ID:                 organization.ID,
			OrganizationTypeID: organization.OrganizationTypeID,
			Name:               organization.Name,
			OrganizationType: response.OrganizationTypeResponse{
				ID:   organization.OrganizationType.ID,
				Name: organization.OrganizationType.Name,
			},
			OrganizationStructureResponses: *ConvertToOrganizationStructureMinimalResponse(&organization.OrganizationStructures),
		})
	}
	return &responseOrganizations
}

