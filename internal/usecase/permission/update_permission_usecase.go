package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IUpdatePermissionUseCaseRequest struct {
	ID            uuid.UUID          `json:"id"`
	Permission    *entity.Permission `json:"permission"`
	ApplicationID uuid.UUID          `json:"application_id"`
}

type IUpdatePermissionUseCaseResponse struct {
	Permission *entity.Permission `json:"permission"`
}

type IUpdatePermissionUseCase interface {
	Execute(request *IUpdatePermissionUseCaseRequest) (*IUpdatePermissionUseCaseResponse, error)
}

type UpdatePermissionUseCase struct {
	Log                  *logrus.Logger
	PermissionRepository repository.IPermissionRepository
}

func NewUpdatePermissionUseCase(log *logrus.Logger, permissionRepository repository.IPermissionRepository) IUpdatePermissionUseCase {
	return &UpdatePermissionUseCase{
		Log:                  log,
		PermissionRepository: permissionRepository,
	}
}

func (uc *UpdatePermissionUseCase) Execute(request *IUpdatePermissionUseCaseRequest) (*IUpdatePermissionUseCaseResponse, error) {
	uc.Log.Info("UpdatePermissionUseCase.Execute")

	permissionExist, err := uc.PermissionRepository.FindById(request.ID)
	if err != nil {
		return nil, err
	}

	if permissionExist == nil {
		return nil, err
	}

	permission, err := uc.PermissionRepository.UpdatePermission(&entity.Permission{
		ID:            request.ID,
		Name:          request.Permission.Name,
		ApplicationID: request.ApplicationID,
		GuardName:     request.Permission.GuardName,
		Label:         request.Permission.Label,
		Description:   request.Permission.Description,
	})
	if err != nil {
		return nil, err
	}

	return &IUpdatePermissionUseCaseResponse{
		Permission: permission,
	}, nil
}

func UpdatePermissionUseCaseFactory(log *logrus.Logger) IUpdatePermissionUseCase {
	permissionRepository := repository.PermissionRepositoryFactory(log)
	return &UpdatePermissionUseCase{
		Log:                  log,
		PermissionRepository: permissionRepository,
	}
}
