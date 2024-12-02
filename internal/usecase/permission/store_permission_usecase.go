package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IStorePermissionUseCaseRequest struct {
	Permission    *entity.Permission `json:"permission"`
	ApplicationID uuid.UUID          `json:"application_id"`
}

type IStorePermissionUseCaseResponse struct {
	Permission *entity.Permission `json:"permission"`
}

type IStorePermissionUseCase interface {
	Execute(request *IStorePermissionUseCaseRequest) (*IStorePermissionUseCaseResponse, error)
}

type StorePermissionUseCase struct {
	Log                  *logrus.Logger
	PermissionRepository repository.IPermissionRepository
}

func NewStorePermissionUseCase(log *logrus.Logger, permissionRepository repository.IPermissionRepository) IStorePermissionUseCase {
	return &StorePermissionUseCase{
		Log:                  log,
		PermissionRepository: permissionRepository,
	}
}

func (uc *StorePermissionUseCase) Execute(request *IStorePermissionUseCaseRequest) (*IStorePermissionUseCaseResponse, error) {
	uc.Log.Info("StorePermissionUseCase.Execute")

	permission, err := uc.PermissionRepository.StorePermission(&entity.Permission{
		Name:          request.Permission.Name,
		ApplicationID: request.ApplicationID,
		GuardName:     request.Permission.GuardName,
		Label:         request.Permission.Label,
	})
	if err != nil {
		return nil, err
	}

	return &IStorePermissionUseCaseResponse{
		Permission: permission,
	}, nil
}

func StorePermissionUseCaseFactory(log *logrus.Logger) IStorePermissionUseCase {
	permissionRepository := repository.PermissionRepositoryFactory(log)
	return &StorePermissionUseCase{
		Log:                  log,
		PermissionRepository: permissionRepository,
	}
}
