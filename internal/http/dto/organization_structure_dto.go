package dto

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/http/response"
)

func ConvertToOrganizationStructureResponse(orgStructures *[]entity.OrganizationStructure) *[]response.OrganizationStructureResponse {
	var responseStructures []response.OrganizationStructureResponse
	for _, orgStructure := range *orgStructures {
		var parentResponse *response.ParentOrganizationStructureResponse
		if orgStructure.Parent != nil {
			parentResponse = &response.ParentOrganizationStructureResponse{ID: orgStructure.Parent.ID, Name: orgStructure.Parent.Name}
		}
		responseStructures = append(responseStructures, response.OrganizationStructureResponse{
			ID:             orgStructure.ID,
			Name:           orgStructure.Name,
			Level:          orgStructure.Level,
			Path:           orgStructure.Path,
			OrganizationID: orgStructure.OrganizationID,
			JobLevelID:     orgStructure.JobLevelID,
			ParentID:       orgStructure.ParentID,
			Parent:         parentResponse,
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

func ConvertToOrganizationStructureParentResponse(orgStructures *[]entity.OrganizationStructure) *[]response.OrganizationStructureParentResponse {
	var responseStructures []response.OrganizationStructureParentResponse
	for _, orgStructure := range *orgStructures {
		responseStructures = append(responseStructures, response.OrganizationStructureParentResponse{
			ID:             orgStructure.ID,
			Name:           orgStructure.Name,
			Level:          orgStructure.Level,
			Path:           orgStructure.Path,
			OrganizationID: orgStructure.OrganizationID,
			JobLevelID:     orgStructure.JobLevelID,
			ParentID:       orgStructure.ParentID,
			Parents:        *ConvertToOrganizationStructureParentResponse(&orgStructure.Parents),
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

func ConvertToSingleOrganizationStructureParentResponse(orgStructure *entity.OrganizationStructure) *response.OrganizationStructureParentResponse {
	return &response.OrganizationStructureParentResponse{
		ID:             orgStructure.ID,
		Name:           orgStructure.Name,
		Level:          orgStructure.Level,
		Path:           orgStructure.Path,
		OrganizationID: orgStructure.OrganizationID,
		JobLevelID:     orgStructure.JobLevelID,
		ParentID:       orgStructure.ParentID,
		Parents:        *ConvertToOrganizationStructureParentResponse(&orgStructure.Parents),
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

func ConvertToOrganizationStructureMinimalResponse(orgStructures *[]entity.OrganizationStructure) *[]response.OrganizationStructureMinimalResponse {
	var responseStructures []response.OrganizationStructureMinimalResponse
	for _, orgStructure := range *orgStructures {
		responseStructures = append(responseStructures, response.OrganizationStructureMinimalResponse{
			ID:             orgStructure.ID,
			Name:           orgStructure.Name,
			Level:          orgStructure.Level,
			OrganizationID: orgStructure.OrganizationID,
			ParentID:       orgStructure.ParentID,
		})
	}
	return &responseStructures
}
