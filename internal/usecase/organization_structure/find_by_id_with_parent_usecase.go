package usecase

import (
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IFindByIDWithParentUseCaseRequest struct {
	ID uuid.UUID `json:"id"`
}

type IFindByIDWithParentUseCaseResponse struct {
	OrganizationStructure *response.OrganizationStructureParentResponse
}

type IFindByIDWithParentUseCase interface {
	Execute(request *IFindByIDWithParentUseCaseRequest) (*IFindByIDWithParentUseCaseResponse, error)
}

type FindByIDWithParentUseCase struct {
	Log                             *logrus.Logger
	OrganizationStructureRepository repository.IOrganizationStructureRepository
}

func NewFindByIDWithParentUseCase(log *logrus.Logger, organizationStructureRepository repository.IOrganizationStructureRepository) IFindByIDWithParentUseCase {
	return &FindByIDWithParentUseCase{
		Log:                             log,
		OrganizationStructureRepository: organizationStructureRepository,
	}
}

func (u *FindByIDWithParentUseCase) Execute(request *IFindByIDWithParentUseCaseRequest) (*IFindByIDWithParentUseCaseResponse, error) {
	organizationStructure, err := u.OrganizationStructureRepository.FindById(request.ID)
	if err != nil {
		return nil, err
	}

	parents, err := u.OrganizationStructureRepository.FindAllParents(organizationStructure.ID)
	if err != nil {
		return nil, err
	}

	(*organizationStructure).Parents = parents

	return &IFindByIDWithParentUseCaseResponse{
		OrganizationStructure: dto.ConvertToSingleOrganizationStructureParentResponse(organizationStructure),
	}, nil
}

func FindByIDWithParentUseCaseFactory(log *logrus.Logger) IFindByIDWithParentUseCase {
	organizationStructureRepository := repository.OrganizationStructureRepositoryFactory(log)
	return NewFindByIDWithParentUseCase(log, organizationStructureRepository)
}
