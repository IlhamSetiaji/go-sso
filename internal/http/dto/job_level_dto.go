package dto

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/http/response"
)

func ConvertToJobLevelResponse(jobLevels *[]entity.JobLevel) *[]response.JobLevelResponse {
	var responseJobLevels []response.JobLevelResponse
	for _, jobLevel := range *jobLevels {
		responseJobLevels = append(responseJobLevels, response.JobLevelResponse{
			ID:    jobLevel.ID,
			Name:  jobLevel.Name,
			Level: jobLevel.Level,
			OrganizationStructures: func() *[]response.OrganizationStructureMinimalResponse {
				var responseOrganizationStructures []response.OrganizationStructureMinimalResponse
				for _, organizationStructure := range jobLevel.OrganizationStructures {
					responseOrganizationStructures = append(responseOrganizationStructures, response.OrganizationStructureMinimalResponse{
						ID:             organizationStructure.ID,
						Name:           organizationStructure.Name,
						Level:          organizationStructure.Level,
						OrganizationID: organizationStructure.OrganizationID,
						Path:           organizationStructure.Path,
					})
				}
				return &responseOrganizationStructures
			}(),
		})
	}
	return &responseJobLevels
}

func ConvertToSingleJobLevelResponse(jobLevel *entity.JobLevel) *response.JobLevelResponse {
	return &response.JobLevelResponse{
		ID:    jobLevel.ID,
		Name:  jobLevel.Name,
		Level: jobLevel.Level,
		OrganizationStructures: func() *[]response.OrganizationStructureMinimalResponse {
			var responseOrganizationStructures []response.OrganizationStructureMinimalResponse
			for _, organizationStructure := range jobLevel.OrganizationStructures {
				responseOrganizationStructures = append(responseOrganizationStructures, response.OrganizationStructureMinimalResponse{
					ID:             organizationStructure.ID,
					Name:           organizationStructure.Name,
					Level:          organizationStructure.Level,
					OrganizationID: organizationStructure.OrganizationID,
					Path:           organizationStructure.Path,
				})
			}
			return &responseOrganizationStructures
		}(),
	}
}
