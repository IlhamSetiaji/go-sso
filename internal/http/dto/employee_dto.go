package dto

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/http/response"
	"fmt"

	gt "github.com/bas24/googletranslatefree"
)

func ConvertToSingleEmployeeResponse(employee *entity.Employee) *response.EmployeeResponse {
	chinese, err := gt.Translate(employee.EmployeeJob.Job.Name, "en", "zh")
	if err != nil {
		fmt.Println(err)
		chinese = ""
	}
	return &response.EmployeeResponse{
		ID:             employee.ID,
		Name:           employee.Name,
		OrganizationID: employee.OrganizationID,
		Email:          employee.Email,
		MobilePhone:    employee.MobilePhone,
		EndDate:        employee.EndDate,
		RetirementDate: employee.RetirementDate,
		Organization: response.OrganizationResponse{
			ID: employee.Organization.ID,
			OrganizationType: response.OrganizationTypeResponse{
				ID:   employee.Organization.OrganizationType.ID,
				Name: employee.Organization.OrganizationType.Name,
			},
			Name: employee.Organization.Name,
		},
		EmployeeJob: map[string]interface{}{
			"id":                       employee.EmployeeJob.ID,
			"name":                     employee.EmployeeJob.Name,
			"emp_organization_id":      employee.EmployeeJob.EmpOrganizationID,
			"job_id":                   employee.EmployeeJob.JobID,
			"job_name":                 employee.EmployeeJob.Job.Name,
			"job_name_chinese":         chinese,
			"employee_id":              employee.EmployeeJob.EmployeeID,
			"organization_location_id": employee.EmployeeJob.OrganizationLocationID,
		},
	}
}
