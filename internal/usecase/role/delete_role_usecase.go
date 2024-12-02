package usecase

import (
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IDeleteRoleUseCaseRequest struct {
	ID uuid.UUID `json:"id"`
}

type IDeleteRoleUseCase interface {
	Execute(request *IDeleteRoleUseCaseRequest) error
}

type DeleteRoleUseCase struct {
	Log            *logrus.Logger
	RoleRepository repository.IRoleRepository
}

func NewDeleteRoleUseCase(log *logrus.Logger, roleRepository repository.IRoleRepository) IDeleteRoleUseCase {
	return &DeleteRoleUseCase{
		Log:            log,
		RoleRepository: roleRepository,
	}
}

func (uc *DeleteRoleUseCase) Execute(request *IDeleteRoleUseCaseRequest) error {
	uc.Log.Info("DeleteRoleUseCase.Execute")

	err := uc.RoleRepository.DeleteRole(request.ID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteRoleUseCaseFactory(log *logrus.Logger) IDeleteRoleUseCase {
	roleRepository := repository.RoleRepositoryFactory(log)
	return &DeleteRoleUseCase{
		Log:            log,
		RoleRepository: roleRepository,
	}
}
