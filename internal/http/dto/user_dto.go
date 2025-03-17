package dto

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/http/response"

	"github.com/google/uuid"
)

func ConvertToSingleUserResponse(user *entity.User) *response.UserResponse {
	return &response.UserResponse{
		ID:                  user.ID,
		ChoosedRole:         user.ChoosedRole,
		VerifiedUserProfile: user.VerifiedUserProfile,
		Username:            user.Username,
		Email:               user.Email,
		EmployeeID: func() uuid.UUID {
			if user.EmployeeID == nil {
				return uuid.UUID{}
			}
			return *user.EmployeeID
		}(),
		Name:            user.Name,
		MobilePhone:     user.MobilePhone,
		OauthID:         user.OauthID,
		EmailVerifiedAt: user.EmailVerifiedAt,
		Gender:          user.Gender,
		Photo:           user.Photo,
		Status:          user.Status,
		Address:         user.Address,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
		Roles: func() []response.RoleResponse {
			var roles []response.RoleResponse
			for _, role := range user.Roles {
				roles = append(roles, response.RoleResponse{
					ID:              role.ID,
					ApplicationID:   role.ApplicationID,
					ApplicationName: role.Application.Name,
					Name:            role.Name,
					GuardName:       role.GuardName,
					Status:          string(role.Status),
					CreatedAt:       role.CreatedAt,
					UpdatedAt:       role.UpdatedAt,
					Permissions: func() []response.PermissionResponse {
						var permissions []response.PermissionResponse
						for _, permission := range role.Permissions {
							permissions = append(permissions, response.PermissionResponse{
								ID:            permission.ID,
								ApplicationID: permission.ApplicationID,
								Name:          permission.Name,
								Label:         permission.Label,
								GuardName:     permission.GuardName,
							})
						}
						return permissions
					}(),
				})
			}
			return roles
		}(),
		UserProfile: func() map[string]interface{} {
			if user.UserProfile == nil {
				return map[string]interface{}{}
			} else {
				return user.UserProfile
			}
		}(),
		Employee: func() response.EmployeeResponse {
			if user.Employee == nil {
				return response.EmployeeResponse{}
			}
			return response.EmployeeResponse{
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
				EmployeeJob: func() *response.EmployeeJobResponse {
					if user.Employee.EmployeeJob == nil {
						return nil
					}
					return &response.EmployeeJobResponse{
						ID:                     user.Employee.EmployeeJob.ID,
						EmpOrganizationID:      *user.Employee.EmployeeJob.EmpOrganizationID,
						OrganizationLocationID: user.Employee.EmployeeJob.OrganizationLocationID,
						EmployeeID:             user.Employee.ID,
						JobID:                  user.Employee.EmployeeJob.JobID,
						Name:                   user.Employee.EmployeeJob.Name,
						Job: func() *response.JobResponse {
							if user.Employee.EmployeeJob.Job == nil {
								return nil
							}
							return ConvertToSingleJobResponse(user.Employee.EmployeeJob.Job)
						}(),
						EmpOrganization: func() *response.OrganizationMinimalResponse {
							if user.Employee.EmployeeJob.EmpOrganization == nil {
								return nil
							}
							return &response.OrganizationMinimalResponse{
								ID:   user.Employee.EmployeeJob.EmpOrganization.ID,
								Name: user.Employee.EmployeeJob.EmpOrganization.Name,
							}
						}(),
						OrganizationLocation: func() *response.OrganizationLocationResponse {
							if user.Employee.EmployeeJob.OrganizationLocation == nil {
								return nil
							}
							return ConvertToSingleOrganizationLocationResponse(user.Employee.EmployeeJob.OrganizationLocation)
						}(),
						OrganizationStructure: func() *response.OrganizationStructureResponse {
							if user.Employee.EmployeeJob.OrganizationStructure == nil {
								return nil
							}
							return ConvertToSingleOrganizationStructureResponse(user.Employee.EmployeeJob.OrganizationStructure)
						}(),
						CreatedAt: user.Employee.EmployeeJob.CreatedAt,
						UpdatedAt: user.Employee.EmployeeJob.UpdatedAt,
					}
				}(),
			}
		}(),
	}
}
