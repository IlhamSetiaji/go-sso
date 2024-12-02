package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IGetAllPermissionsUseCaseResponse struct {
	Permissions *[]entity.Permission `json:"permissions"`
}

type IGetAllPermissionsUseCase interface {
	Execute() (*IGetAllPermissionsUseCaseResponse, error)
}

type GetAllPermissionsUseCase struct {
	Log                  *logrus.Logger
	PermissionRepository repository.IPermissionRepository
}

func NewGetAllPermissionsUseCase(log *logrus.Logger, permissionRepository repository.IPermissionRepository) IGetAllPermissionsUseCase {
	return &GetAllPermissionsUseCase{
		Log:                  log,
		PermissionRepository: permissionRepository,
	}
}

func (uc *GetAllPermissionsUseCase) Execute() (*IGetAllPermissionsUseCaseResponse, error) {
	permissions, err := uc.PermissionRepository.GetAllPermissions()
	if err != nil {
		return nil, err
	}

	return &IGetAllPermissionsUseCaseResponse{
		Permissions: permissions,
	}, nil
}

func GetAllPermissionsUseCaseFactory(log *logrus.Logger) IGetAllPermissionsUseCase {
	permissionRepository := repository.PermissionRepositoryFactory(log)
	return &GetAllPermissionsUseCase{
		Log:                  log,
		PermissionRepository: permissionRepository,
	}
}
