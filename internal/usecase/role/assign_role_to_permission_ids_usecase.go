package usecase

import (
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IAssignRoleToPermissionIDsUsecaseRequest struct {
	RoleID        string   `form:"role_id" validate:"required"`
	PermissionIDs []string `form:"permission_ids[]" validate:"required,dive"`
}

type IAssignRoleToPermissionIDsUsecaseResponse struct {
	RoleID string `json:"role_id"`
}

type AssignRoleToPermissionIDsUsecase struct {
	Log            *logrus.Logger
	RoleRepo       repository.IRoleRepository
	PermissionRepo repository.IPermissionRepository
}

type AssignRoleToPermissionIDsUsecaseInterface interface {
	Execute(request *IAssignRoleToPermissionIDsUsecaseRequest) (*IAssignRoleToPermissionIDsUsecaseResponse, error)
}

func AssignRoleToPermissionIDsUsecaseFactory(log *logrus.Logger) AssignRoleToPermissionIDsUsecaseInterface {
	roleRepo := repository.RoleRepositoryFactory(log)
	permissionRepo := repository.PermissionRepositoryFactory(log)
	return &AssignRoleToPermissionIDsUsecase{
		Log:            log,
		RoleRepo:       roleRepo,
		PermissionRepo: permissionRepo,
	}
}

func (u *AssignRoleToPermissionIDsUsecase) Execute(request *IAssignRoleToPermissionIDsUsecaseRequest) (*IAssignRoleToPermissionIDsUsecaseResponse, error) {
	roleID, err := uuid.Parse(request.RoleID)
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	role, err := u.RoleRepo.FindById(roleID)
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	role, err = u.RoleRepo.AssignRoleToPermissions(role, request.PermissionIDs)
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	return &IAssignRoleToPermissionIDsUsecaseResponse{
		RoleID: role.ID.String(),
	}, nil
}
