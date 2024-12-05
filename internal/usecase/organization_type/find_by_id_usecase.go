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
	OrganizationType *entity.OrganizationType `json:"organization_type"`
}

type IFindByIdUseCase interface {
	Execute(request *IFindByIdUseCaseRequest) (*IFindByIdUseCaseResponse, error)
}

type FindByIdUseCase struct {
	Log                        *logrus.Logger
	OrganizationTypeRepository repository.IOrganizationTypeRepository
}

func NewFindByIdUseCase(
	log *logrus.Logger,
	organizationTypeRepository repository.IOrganizationTypeRepository,
) IFindByIdUseCase {
	return &FindByIdUseCase{
		Log:                        log,
		OrganizationTypeRepository: organizationTypeRepository,
	}
}

func (uc *FindByIdUseCase) Execute(req *IFindByIdUseCaseRequest) (*IFindByIdUseCaseResponse, error) {
	organizationType, err := uc.OrganizationTypeRepository.FindById(req.ID)
	if err != nil {
		return nil, err
	}

	return &IFindByIdUseCaseResponse{
		OrganizationType: organizationType,
	}, nil
}

func FindByIdUseCaseFactory(log *logrus.Logger) IFindByIdUseCase {
	organizationTypeRepository := repository.OrganizationTypeRepositoryFactory(log)
	return NewFindByIdUseCase(log, organizationTypeRepository)
}
