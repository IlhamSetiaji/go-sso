package dto

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/http/response"
)

func ConvertToJobResponse(jobs *[]entity.Job) *[]response.JobResponse {
	var responseJobs []response.JobResponse
	for _, job := range *jobs {
		responseJobs = append(responseJobs, response.JobResponse{
			ID:                      job.ID,
			Name:                    job.Name,
			OrganizationStructureID: job.OrganizationStructureID,
			Level:                   job.Level,
			ParentID:                job.ParentID,
			Path:                    job.Path,
			Existing:                job.Existing,
			Children:                *ConvertToJobResponse(&job.Children),
		})
	}
	return &responseJobs
}

func ConvertToSingleJobResponse(job *entity.Job) *response.JobResponse {
	return &response.JobResponse{
		ID:                      job.ID,
		Name:                    job.Name,
		OrganizationStructureID: job.OrganizationStructureID,
		Level:                   job.Level,
		ParentID:                job.ParentID,
		Path:                    job.Path,
		Existing:                job.Existing,
		Children:                *ConvertToJobResponse(&job.Children),
	}
}
