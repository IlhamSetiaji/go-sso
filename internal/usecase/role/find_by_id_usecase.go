package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IFindByIdUseCaseRequest struct {
	ID uuid.UUID `json:"id"`
}

type IFindByIdUseCaseResponse struct {
	Role *entity.Role `json:"role"`
}

type IFindByIdUseCase interface {
	Execute(request *IFindByIdUseCaseRequest) (*IFindByIdUseCaseResponse, error)
}

type FindByIdUseCase struct {
	Log            *logrus.Logger
	RoleRepository repository.IRoleRepository
}

func NewFindByIdUseCase(log *logrus.Logger, roleRepository repository.IRoleRepository) IFindByIdUseCase {
	return &FindByIdUseCase{
		Log:            log,
		RoleRepository: roleRepository,
	}
}

func (u *FindByIdUseCase) Execute(request *IFindByIdUseCaseRequest) (*IFindByIdUseCaseResponse, error) {
	role, err := u.RoleRepository.FindById(request.ID)
	if err != nil {
		u.Log.Error("[FindByIdUseCase.Execute] " + err.Error())
		return nil, err
	}

	return &IFindByIdUseCaseResponse{
		Role: role,
	}, nil
}

func FindByIdUseCaseFactory(log *logrus.Logger) IFindByIdUseCase {
	roleRepository := repository.RoleRepositoryFactory(log)
	return NewFindByIdUseCase(log, roleRepository)
}
