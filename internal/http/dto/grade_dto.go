package dto

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/http/response"
)

func ConvertToSingleGradeResponse(grade *entity.Grade) *response.GradeResponse {
	return &response.GradeResponse{
		ID:         grade.ID,
		JobLevelID: grade.JobLevelID,
		Name:       grade.Name,
		MidsuitID:  grade.MidsuitID,
		CreatedAt:  grade.CreatedAt,
		UpdatedAt:  grade.UpdatedAt,
		JobLevel: func() *response.JobLevelResponse {
			if grade.JobLevel == nil {
				return nil
			}
			resp := ConvertToSingleJobLevelResponse(grade.JobLevel)
			return resp
		}(),
	}
}

func ConvertToGradeResponse(grades *[]entity.Grade) *[]response.GradeResponse {
	var responseGrades []response.GradeResponse
	for _, grade := range *grades {
		responseGrades = append(responseGrades, *ConvertToSingleGradeResponse(&grade))
	}
	return &responseGrades
}
