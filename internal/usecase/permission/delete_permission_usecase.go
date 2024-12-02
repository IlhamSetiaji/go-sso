package usecase

import (
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IDeletePermissionUseCaseRequest struct {
	ID uuid.UUID `json:"id"`
}

type IDeletePermissionUseCase interface {
	Execute(request *IDeletePermissionUseCaseRequest) error
}

type DeletePermissionUseCase struct {
	Log                  *logrus.Logger
	PermissionRepository repository.IPermissionRepository
}

func NewDeletePermissionUseCase(log *logrus.Logger, permissionRepository repository.IPermissionRepository) IDeletePermissionUseCase {
	return &DeletePermissionUseCase{
		Log:                  log,
		PermissionRepository: permissionRepository,
	}
}

func (uc *DeletePermissionUseCase) Execute(request *IDeletePermissionUseCaseRequest) error {
	uc.Log.Info("DeletePermissionUseCase.Execute")

	err := uc.PermissionRepository.DeletePermission(request.ID)
	if err != nil {
		return err
	}

	return nil
}

func DeletePermissionUseCaseFactory(log *logrus.Logger) IDeletePermissionUseCase {
	permissionRepository := repository.PermissionRepositoryFactory(log)
	return &DeletePermissionUseCase{
		Log:                  log,
		PermissionRepository: permissionRepository,
	}
}
