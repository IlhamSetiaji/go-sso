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
	Organization *entity.Organization `json:"organization"`
}

type IFindByIdUseCase interface {
	Execute(request *IFindByIdUseCaseRequest) (*IFindByIdUseCaseResponse, error)
}

type FindByIdUseCase struct {
	Log                    *logrus.Logger
	OrganizationRepository repository.IOrganizationRepository
}

func NewFindByIdUseCase(
	log *logrus.Logger,
	organizationRepository repository.IOrganizationRepository,
) IFindByIdUseCase {
	return &FindByIdUseCase{
		Log:                    log,
		OrganizationRepository: organizationRepository,
	}
}

func (uc *FindByIdUseCase) Execute(req *IFindByIdUseCaseRequest) (*IFindByIdUseCaseResponse, error) {
	organization, err := uc.OrganizationRepository.FindById(req.ID)
	if err != nil {
		return nil, err
	}

	return &IFindByIdUseCaseResponse{
		Organization: organization,
	}, nil
}

func FindByIdUseCaseFactory(log *logrus.Logger) IFindByIdUseCase {
	organizationRepository := repository.OrganizationRepositoryFactory(log)
	return NewFindByIdUseCase(log, organizationRepository)
}
