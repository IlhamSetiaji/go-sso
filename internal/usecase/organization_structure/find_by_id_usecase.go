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

type IFindByIdUseCase interface {
	Execute(request *IFindByIdUseCaseRequest) (IFindByIdUseCaseResponse, error)
}

type IFindByIdUseCaseResponse struct {
	OrganizationStructure *entity.OrganizationStructure `json:"organization_structure"`
}

type FindByIdUseCase struct {
	Log                             *logrus.Logger
	OrganizationStructureRepository repository.IOrganizationStructureRepository
}

func NewFindByIdUseCase(log *logrus.Logger, organizationStructureRepository repository.IOrganizationStructureRepository) IFindByIdUseCase {
	return &FindByIdUseCase{
		Log:                             log,
		OrganizationStructureRepository: organizationStructureRepository,
	}
}

func (u *FindByIdUseCase) Execute(request *IFindByIdUseCaseRequest) (IFindByIdUseCaseResponse, error) {
	organizationStructure, err := u.OrganizationStructureRepository.FindById(request.ID)
	if err != nil {
		return IFindByIdUseCaseResponse{}, err
	}

	return IFindByIdUseCaseResponse{
		OrganizationStructure: organizationStructure,
	}, nil
}

func FindByIdUseCaseFactory(log *logrus.Logger) IFindByIdUseCase {
	organizationStructureRepository := repository.OrganizationStructureRepositoryFactory(log)
	return NewFindByIdUseCase(log, organizationStructureRepository)
}
