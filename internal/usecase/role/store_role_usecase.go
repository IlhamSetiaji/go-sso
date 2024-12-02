package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IStoreRoleUseCaseRequest struct {
	Role          *entity.Role `json:"role"`
	ApplicationID uuid.UUID    `json:"application_id"`
}

type IStoreRoleUseCaseResponse struct {
	Role *entity.Role `json:"role"`
}

type IStoreRoleUseCase interface {
	Execute(request *IStoreRoleUseCaseRequest) (*IStoreRoleUseCaseResponse, error)
}

type StoreRoleUseCase struct {
	Log            *logrus.Logger
	RoleRepository repository.IRoleRepository
}

func NewStoreRoleUseCase(log *logrus.Logger, roleRepository repository.IRoleRepository) IStoreRoleUseCase {
	return &StoreRoleUseCase{
		Log:            log,
		RoleRepository: roleRepository,
	}
}

func (uc *StoreRoleUseCase) Execute(request *IStoreRoleUseCaseRequest) (*IStoreRoleUseCaseResponse, error) {
	uc.Log.Info("StoreRoleUseCase.Execute")

	role, err := uc.RoleRepository.StoreRole(&entity.Role{
		Name:          request.Role.Name,
		ApplicationID: request.ApplicationID,
		GuardName:     request.Role.GuardName,
		Status:        request.Role.Status,
	})
	if err != nil {
		return nil, err
	}

	return &IStoreRoleUseCaseResponse{
		Role: role,
	}, nil
}

func StoreRoleUseCaseFactory(log *logrus.Logger) IStoreRoleUseCase {
	roleRepository := repository.RoleRepositoryFactory(log)
	return &StoreRoleUseCase{
		Log:            log,
		RoleRepository: roleRepository,
	}
}
