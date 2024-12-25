package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IGetAllPermissionsByRoleIDUsecaseRequest struct {
	RoleID string `json:"role_id" validate:"required"`
}

type IGetAllPermissionsByRoleIDUsecaseResponse struct {
	Permissions []entity.Permission `json:"permissions"`
}

type GetAllPermissionsByRoleIDUsecase struct {
	Log            *logrus.Logger
	PermissionRepo repository.IPermissionRepository
}

type GetAllPermissionsByRoleIDUsecaseInterface interface {
	Execute(request *IGetAllPermissionsByRoleIDUsecaseRequest) (*IGetAllPermissionsByRoleIDUsecaseResponse, error)
}

func GetAllPermissionsByRoleIDUsecaseFactory(log *logrus.Logger) GetAllPermissionsByRoleIDUsecaseInterface {
	permissionRepo := repository.PermissionRepositoryFactory(log)
	return &GetAllPermissionsByRoleIDUsecase{
		Log:            log,
		PermissionRepo: permissionRepo,
	}
}

func (u *GetAllPermissionsByRoleIDUsecase) Execute(request *IGetAllPermissionsByRoleIDUsecaseRequest) (*IGetAllPermissionsByRoleIDUsecaseResponse, error) {
	roleID, err := uuid.Parse(request.RoleID)
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	permissions, err := u.PermissionRepo.GetAllPermissionsByRoleID(roleID)
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	return &IGetAllPermissionsByRoleIDUsecaseResponse{
		Permissions: *permissions,
	}, nil
}
