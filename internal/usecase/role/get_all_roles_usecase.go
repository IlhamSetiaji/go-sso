package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IGetAllRolesUseCaseResponse struct {
	Roles *[]entity.Role `json:"roles"`
}

type IGetAllRolesUseCase interface {
	Execute() (*IGetAllRolesUseCaseResponse, error)
}

type GetAllRolesUseCase struct {
	Log            *logrus.Logger
	RoleRepository repository.IRoleRepository
}

func NewGetAllRolesUseCase(log *logrus.Logger, roleRepository repository.IRoleRepository) IGetAllRolesUseCase {
	return &GetAllRolesUseCase{
		Log:            log,
		RoleRepository: roleRepository,
	}
}

func (uc *GetAllRolesUseCase) Execute() (*IGetAllRolesUseCaseResponse, error) {
	roles, err := uc.RoleRepository.GetAllRoles()
	if err != nil {
		return nil, err
	}

	return &IGetAllRolesUseCaseResponse{
		Roles: roles,
	}, nil
}

func GetAllRolesUseCaseFactory(log *logrus.Logger) IGetAllRolesUseCase {
	roleRepository := repository.RoleRepositoryFactory(log)
	return &GetAllRolesUseCase{
		Log:            log,
		RoleRepository: roleRepository,
	}
}
