package dto

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/http/response"
)

func ConvertToJobResponse(jobs *[]entity.Job) *[]response.JobResponse {
	var responseJobs []response.JobResponse
	for _, job := range *jobs {
		var parentResponse *response.ParentJobResponse
		if job.Parent != nil {
			parentResponse = &response.ParentJobResponse{ID: job.Parent.ID, Name: job.Parent.Name}
		}
		responseJobs = append(responseJobs, response.JobResponse{
			ID:                        job.ID,
			Name:                      job.Name,
			OrganizationStructureID:   job.OrganizationStructureID,
			OrganizationStructureName: job.OrganizationStructure.Name,
			OrganizationID:            job.OrganizationStructure.OrganizationID,
			OrganizationName:          job.OrganizationStructure.Organization.Name,
			Level:                     job.Level,
			ParentID:                  job.ParentID,
			Path:                      job.Path,
			Existing:                  job.Existing,
			Promotion:                 job.Promotion,
			JobPlafon:                 job.Plafon,
			// JobLevel:                  *ConvertToSingleJobLevelResponse(&job.OrganizationStructure.JobLevel),
			JobLevel: *ConvertToSingleJobLevelResponse(&job.JobLevel),
			Parent:   parentResponse,
			Children: *ConvertToJobResponse(&job.Children),
		})
	}
	return &responseJobs
}

func ConvertToSingleJobResponse(job *entity.Job) *response.JobResponse {
	var parentResponse *response.ParentJobResponse
	if job.Parent != nil {
		parentResponse = &response.ParentJobResponse{ID: job.Parent.ID, Name: job.Parent.Name}
	}
	return &response.JobResponse{
		ID:                        job.ID,
		Name:                      job.Name,
		OrganizationStructureID:   job.OrganizationStructureID,
		OrganizationStructureName: job.OrganizationStructure.Name,
		OrganizationID:            job.OrganizationStructure.OrganizationID,
		OrganizationName:          job.OrganizationStructure.Organization.Name,
		Level:                     job.Level,
		ParentID:                  job.ParentID,
		Path:                      job.Path,
		Existing:                  job.Existing,
		Promotion:                 job.Promotion,
		JobPlafon:                 job.Plafon,
		JobLevel:                  *ConvertToSingleJobLevelResponse(&job.OrganizationStructure.JobLevel),
		Parent:                    parentResponse,
		Children:                  *ConvertToJobResponse(&job.Children),
	}
}
