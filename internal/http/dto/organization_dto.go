package dto

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/http/response"
)

func ConvertToOrganizationStructureResponse(orgStructures *[]entity.OrganizationStructure) *[]response.OrganizationStructureResponse {
	var responseStructures []response.OrganizationStructureResponse
	for _, orgStructure := range *orgStructures {
		responseStructures = append(responseStructures, response.OrganizationStructureResponse{
			ID:             orgStructure.ID,
			Name:           orgStructure.Name,
			Level:          orgStructure.Level,
			Path:           orgStructure.Path,
			OrganizationID: orgStructure.OrganizationID,
			JobLevelID:     orgStructure.JobLevelID,
			ParentID:       orgStructure.ParentID,
			Children:       *ConvertToOrganizationStructureResponse(&orgStructure.Children),
			Organization: response.OrganizationResponse{
				ID:                 orgStructure.Organization.ID,
				OrganizationTypeID: orgStructure.Organization.OrganizationTypeID,
				Name:               orgStructure.Organization.Name,
				OrganizationType: response.OrganizationTypeResponse{
					ID:   orgStructure.Organization.OrganizationType.ID,
					Name: orgStructure.Organization.OrganizationType.Name,
				},
			},
			JobLevel: response.JobLevelResponse{
				ID:    orgStructure.JobLevel.ID,
				Name:  orgStructure.JobLevel.Name,
				Level: orgStructure.JobLevel.Level,
			},
		})
	}
	return &responseStructures
}

func ConvertToSingleOrganizationStructureResponse(orgStructure *entity.OrganizationStructure) *response.OrganizationStructureResponse {
	return &response.OrganizationStructureResponse{
		ID:             orgStructure.ID,
		Name:           orgStructure.Name,
		Level:          orgStructure.Level,
		Path:           orgStructure.Path,
		OrganizationID: orgStructure.OrganizationID,
		JobLevelID:     orgStructure.JobLevelID,
		ParentID:       orgStructure.ParentID,
		Children:       *ConvertToOrganizationStructureResponse(&orgStructure.Children),
		Organization: response.OrganizationResponse{
			ID:                 orgStructure.Organization.ID,
			OrganizationTypeID: orgStructure.Organization.OrganizationTypeID,
			Name:               orgStructure.Organization.Name,
			OrganizationType: response.OrganizationTypeResponse{
				ID:   orgStructure.Organization.OrganizationType.ID,
				Name: orgStructure.Organization.OrganizationType.Name,
			},
		},
		JobLevel: response.JobLevelResponse{
			ID:    orgStructure.JobLevel.ID,
			Name:  orgStructure.JobLevel.Name,
			Level: orgStructure.JobLevel.Level,
		},
	}
}
