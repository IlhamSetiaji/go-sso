package usecase

import (
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IResignRoleFromPermissionUsecaseRequest struct {
	RoleID       string `form:"role_id" validate:"required"`
	PermissionID string `form:"permission_id" validate:"required"`
}

type IResignRoleFromPermissionUsecaseResponse struct {
	RoleID       string `json:"role_id"`
	PermissionID string `json:"permission_id"`
}

type ResignRoleFromPermissionUsecase struct {
	Log            *logrus.Logger
	RoleRepo       repository.IRoleRepository
	PermissionRepo repository.IPermissionRepository
}

type ResignRoleFromPermissionUsecaseInterface interface {
	Execute(req *IResignRoleFromPermissionUsecaseRequest) (*IResignRoleFromPermissionUsecaseResponse, error)
}

func ResignRoleFromPermissionUsecaseFactory(log *logrus.Logger) ResignRoleFromPermissionUsecaseInterface {
	roleRepo := repository.RoleRepositoryFactory(log)
	permissionRepo := repository.PermissionRepositoryFactory(log)
	return &ResignRoleFromPermissionUsecase{
		Log:            log,
		RoleRepo:       roleRepo,
		PermissionRepo: permissionRepo,
	}
}

func (u *ResignRoleFromPermissionUsecase) Execute(req *IResignRoleFromPermissionUsecaseRequest) (*IResignRoleFromPermissionUsecaseResponse, error) {
	roleID, err := u.RoleRepo.FindById(uuid.MustParse(req.RoleID))
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	permission, err := u.PermissionRepo.FindById(uuid.MustParse(req.PermissionID))
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	role, err := u.RoleRepo.ResignRoleFromPermission(roleID, permission.ID)
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	return &IResignRoleFromPermissionUsecaseResponse{
		RoleID:       role.ID.String(),
		PermissionID: permission.ID.String(),
	}, nil
}
