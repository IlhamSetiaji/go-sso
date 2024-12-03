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
	OrganizationLocation *entity.OrganizationLocation `json:"organization_location"`
}

type IFindByIdUseCase interface {
	Execute(request *IFindByIdUseCaseRequest) (*IFindByIdUseCaseResponse, error)
}

type FindByIdUseCase struct {
	Log        *logrus.Logger
	Repository repository.IOrganizationLocationRepository
}

func NewFindByIdUseCase(
	log *logrus.Logger,
	repository repository.IOrganizationLocationRepository,
) IFindByIdUseCase {
	return &FindByIdUseCase{
		Log:        log,
		Repository: repository,
	}
}

func (uc *FindByIdUseCase) Execute(req *IFindByIdUseCaseRequest) (*IFindByIdUseCaseResponse, error) {
	organizationLocation, err := uc.Repository.FindById(req.ID)
	if err != nil {
		return nil, err
	}

	return &IFindByIdUseCaseResponse{
		OrganizationLocation: organizationLocation,
	}, nil
}

func FindByIdUseCaseFactory(log *logrus.Logger) IFindByIdUseCase {
	repository := repository.OrganizationLocationRepositoryFactory(log)
	return NewFindByIdUseCase(log, repository)
}
