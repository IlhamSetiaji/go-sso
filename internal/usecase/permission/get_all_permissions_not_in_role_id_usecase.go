package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IGetAllPermissionsNotInRoleIDUsecaseRequest struct {
	RoleID string `json:"role_id" validate:"required"`
}

type IGetAllPermissionsNotInRoleIDUsecaseResponse struct {
	Permissions []entity.Permission `json:"permissions"`
}

type GetAllPermissionsNotInRoleIDUsecase struct {
	Log            *logrus.Logger
	PermissionRepo repository.IPermissionRepository
}

type GetAllPermissionsNotInRoleIDUsecaseInterface interface {
	Execute(request *IGetAllPermissionsNotInRoleIDUsecaseRequest) (*IGetAllPermissionsNotInRoleIDUsecaseResponse, error)
}

func GetAllPermissionsNotInRoleIDUsecaseFactory(log *logrus.Logger) GetAllPermissionsNotInRoleIDUsecaseInterface {
	permissionRepo := repository.PermissionRepositoryFactory(log)
	return &GetAllPermissionsNotInRoleIDUsecase{
		Log:            log,
		PermissionRepo: permissionRepo,
	}
}

func (u *GetAllPermissionsNotInRoleIDUsecase) Execute(request *IGetAllPermissionsNotInRoleIDUsecaseRequest) (*IGetAllPermissionsNotInRoleIDUsecaseResponse, error) {
	roleID, err := uuid.Parse(request.RoleID)
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	permissions, err := u.PermissionRepo.GetAllPermissionsNotInRoleID(roleID)
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	return &IGetAllPermissionsNotInRoleIDUsecaseResponse{
		Permissions: *permissions,
	}, nil
}
