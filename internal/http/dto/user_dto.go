package dto

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/http/response"
)

func ConvertToSingleUserResponse(user *entity.User) *response.UserResponse {
	return &response.UserResponse{
		ID:              user.ID,
		Username:        user.Username,
		Email:           user.Email,
		EmployeeID:      user.EmployeeID,
		Name:            user.Name,
		MobilePhone:     user.MobilePhone,
		OauthID:         user.OauthID,
		EmailVerifiedAt: user.EmailVerifiedAt,
		Gender:          user.Gender,
		Photo:           user.Photo,
		Status:          user.Status,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,

		Employee: response.EmployeeResponse{
			ID:             user.Employee.ID,
			OrganizationID: user.Employee.OrganizationID,
			Name:           user.Employee.Name,
			EndDate:        user.Employee.EndDate,
			RetirementDate: user.Employee.RetirementDate,
			Email:          user.Employee.Email,
			MobilePhone:    user.Employee.MobilePhone,
			Organization: response.OrganizationResponse{
				ID:                 user.Employee.Organization.ID,
				Name:               user.Employee.Organization.Name,
				OrganizationTypeID: user.Employee.Organization.OrganizationTypeID,
				OrganizationType: response.OrganizationTypeResponse{
					ID:   user.Employee.Organization.OrganizationType.ID,
					Name: user.Employee.Organization.OrganizationType.Name,
				},
				OrganizationStructureResponses: func() []response.OrganizationStructureMinimalResponse {
					var organizationStructures []response.OrganizationStructureMinimalResponse
					for _, organizationStructure := range user.Employee.Organization.OrganizationStructures {
						organizationStructures = append(organizationStructures, response.OrganizationStructureMinimalResponse{
							ID:             organizationStructure.ID,
							Name:           organizationStructure.Name,
							OrganizationID: organizationStructure.OrganizationID,
							ParentID:       organizationStructure.ParentID,
							Level:          organizationStructure.Level,
							Path:           organizationStructure.Path,
							JobLevelResponse: map[string]interface{}{
								"ID":    organizationStructure.JobLevel.ID,
								"Name":  organizationStructure.JobLevel.Name,
								"Level": organizationStructure.JobLevel.Level,
							},
						})
					}
					return organizationStructures
				}(),
			},
			EmployeeJob: map[string]interface{}{
				"id":                       user.Employee.EmployeeJob.ID,
				"name":                     user.Employee.EmployeeJob.Name,
				"emp_organization_id":      user.Employee.EmployeeJob.EmpOrganizationID,
				"job_id":                   user.Employee.EmployeeJob.JobID,
				"employee_id":              user.Employee.EmployeeJob.EmployeeID,
				"organization_location_id": user.Employee.EmployeeJob.OrganizationLocationID,
			},
		},
	}
}
