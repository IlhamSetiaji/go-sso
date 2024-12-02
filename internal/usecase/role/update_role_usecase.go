package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IUpdateRoleUseCaseRequest struct {
	ID            uuid.UUID    `json:"id"`
	Role          *entity.Role `json:"role"`
	ApplicationID uuid.UUID    `json:"application_id"`
}

type IUpdateRoleUseCaseResponse struct {
	Role *entity.Role `json:"role"`
}

type IUpdateRoleUseCase interface {
	Execute(request *IUpdateRoleUseCaseRequest) (*IUpdateRoleUseCaseResponse, error)
}

type UpdateRoleUseCase struct {
	Log            *logrus.Logger
	RoleRepository repository.IRoleRepository
}

func NewUpdateRoleUseCase(log *logrus.Logger, roleRepository repository.IRoleRepository) IUpdateRoleUseCase {
	return &UpdateRoleUseCase{
		Log:            log,
		RoleRepository: roleRepository,
	}
}

func (uc *UpdateRoleUseCase) Execute(request *IUpdateRoleUseCaseRequest) (*IUpdateRoleUseCaseResponse, error) {
	uc.Log.Info("UpdateRoleUseCase.Execute")

	roleExist, err := uc.RoleRepository.FindById(request.ID)
	if err != nil {
		return nil, errors.New("[UpdateRoleUseCase.Execute] " + err.Error())
	}

	if roleExist == nil {
		return nil, errors.New("[UpdateRoleUseCase.Execute] Role not found")
	}

	role, err := uc.RoleRepository.UpdateRole(&entity.Role{
		ID:            request.ID,
		Name:          request.Role.Name,
		ApplicationID: request.ApplicationID,
		GuardName:     request.Role.GuardName,
		Status:        request.Role.Status,
	})
	if err != nil {
		return nil, err
	}

	return &IUpdateRoleUseCaseResponse{
		Role: role,
	}, nil
}

func UpdateRoleUseCaseFactory(log *logrus.Logger) IUpdateRoleUseCase {
	roleRepository := repository.RoleRepositoryFactory(log)
	return &UpdateRoleUseCase{
		Log:            log,
		RoleRepository: roleRepository,
	}
}
