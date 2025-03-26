package dto

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/http/response"
	"fmt"

	gt "github.com/bas24/googletranslatefree"
)

func ConvertToSingleEmployeeResponse(employee *entity.Employee) *response.EmployeeResponse {
	var chinese string
	var err error
	if employee.EmployeeJob != nil {
		if employee.EmployeeJob.Job != nil {
			chinese, err = gt.Translate(employee.EmployeeJob.Job.Name, "en", "zh-CN")
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return &response.EmployeeResponse{
		ID:             employee.ID,
		Name:           employee.Name,
		OrganizationID: employee.OrganizationID,
		Email:          employee.Email,
		MobilePhone:    employee.MobilePhone,
		EndDate:        employee.EndDate,
		RetirementDate: employee.RetirementDate,
		NIK:            employee.NIK,
		Organization: response.OrganizationResponse{
			ID: employee.Organization.ID,
			OrganizationType: response.OrganizationTypeResponse{
				ID:   employee.Organization.OrganizationType.ID,
				Name: employee.Organization.OrganizationType.Name,
			},
			Name: employee.Organization.Name,
		},
		EmployeeJob: func() *response.EmployeeJobResponse {
			if employee.EmployeeJob == nil {
				return nil
			}
			return &response.EmployeeJobResponse{
				ID:                     employee.EmployeeJob.ID,
				EmpOrganizationID:      *employee.EmployeeJob.EmpOrganizationID,
				OrganizationLocationID: employee.EmployeeJob.OrganizationLocationID,
				EmployeeID:             employee.ID,
				JobID:                  employee.EmployeeJob.JobID,
				JobName: func() string {
					if employee.EmployeeJob.Job == nil {
						return ""
					}
					return employee.EmployeeJob.Job.Name
				}(),
				JobNameChinese: chinese,
				Name:           employee.EmployeeJob.Name,
				Job: func() *response.JobResponse {
					if employee.EmployeeJob.Job == nil {
						return nil
					}
					return ConvertToSingleJobResponse(employee.EmployeeJob.Job)
				}(),
				EmpOrganization: func() *response.OrganizationMinimalResponse {
					if employee.EmployeeJob.EmpOrganization == nil {
						return nil
					}
					return &response.OrganizationMinimalResponse{
						ID:   employee.EmployeeJob.EmpOrganization.ID,
						Name: employee.EmployeeJob.EmpOrganization.Name,
					}
				}(),
				OrganizationLocation: func() *response.OrganizationLocationResponse {
					if employee.EmployeeJob.OrganizationLocation == nil {
						return nil
					}
					return ConvertToSingleOrganizationLocationResponse(employee.EmployeeJob.OrganizationLocation)
				}(),
				OrganizationStructure: func() *response.OrganizationStructureResponse {
					if employee.EmployeeJob.OrganizationStructure == nil {
						return nil
					}
					return ConvertToSingleOrganizationStructureResponse(employee.EmployeeJob.OrganizationStructure)
				}(),
				Grade: func() *response.GradeResponse {
					if employee.EmployeeJob.Grade == nil {
						return nil
					}
					return ConvertToSingleGradeResponse(employee.EmployeeJob.Grade)
				}(),
				CreatedAt: employee.EmployeeJob.CreatedAt,
				UpdatedAt: employee.EmployeeJob.UpdatedAt,
			}
		}(),
		KanbanProgress: func() *response.EmployeeKanbanProgressResponse {
			if employee.EmployeeKanbanProgress == nil {
				return nil
			}
			return &response.EmployeeKanbanProgressResponse{
				TotalTask:  employee.EmployeeKanbanProgress.TotalTask,
				ToDo:       employee.EmployeeKanbanProgress.ToDo,
				InProgress: employee.EmployeeKanbanProgress.InProgress,
				NeedReview: employee.EmployeeKanbanProgress.NeedReview,
				Completed:  employee.EmployeeKanbanProgress.Completed,
			}
		}(),
	}
}
