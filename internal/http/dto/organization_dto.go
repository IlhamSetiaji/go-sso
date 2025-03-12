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
			Logo:               organization.Logo,
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
			Region:             organization.Region,
			Address:            organization.Address,
			Logo:               organization.Logo,
			OrganizationType: response.OrganizationTypeResponse{
				ID:   organization.OrganizationType.ID,
				Name: organization.OrganizationType.Name,
			},
			OrganizationStructureResponses: *ConvertToOrganizationStructureMinimalResponse(&organization.OrganizationStructures),
		})
	}
	return &responseOrganizations
}

func ConvertToSingleOrganizationResponse(organization *entity.Organization) *response.OrganizationResponse {
	return &response.OrganizationResponse{
		ID:                 organization.ID,
		OrganizationTypeID: organization.OrganizationTypeID,
		Name:               organization.Name,
		Region:             organization.Region,
		Address:            organization.Address,
		Logo:               organization.Logo,
		OrganizationType: response.OrganizationTypeResponse{
			ID:   organization.OrganizationType.ID,
			Name: organization.OrganizationType.Name,
		},
		OrganizationStructureResponses: *ConvertToOrganizationStructureMinimalResponse(&organization.OrganizationStructures),
	}
}
