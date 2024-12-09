package usecase

import (
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IFindByOrganizationIdUseCaseRequest struct {
	OrganizationID uuid.UUID `json:"organization_id"`
}

type IFindByOrganizationIdUseCaseResponse struct {
	OrganizationLocations *[]response.OrganizationLocationResponse `json:"organization_locations"`
}

type IFindByOrganizationIdUseCase interface {
	Execute(request *IFindByOrganizationIdUseCaseRequest) (*IFindByOrganizationIdUseCaseResponse, error)
}

type FindByOrganizationIdUseCase struct {
	Log        *logrus.Logger
	Repository repository.IOrganizationLocationRepository
}

func NewFindByOrganizationIdUseCase(
	log *logrus.Logger,
	repository repository.IOrganizationLocationRepository,
) IFindByOrganizationIdUseCase {
	return &FindByOrganizationIdUseCase{
		Log:        log,
		Repository: repository,
	}
}

func (uc *FindByOrganizationIdUseCase) Execute(req *IFindByOrganizationIdUseCaseRequest) (*IFindByOrganizationIdUseCaseResponse, error) {
	organizationLocations, err := uc.Repository.FindByOrganizationID(req.OrganizationID)
	if err != nil {
		return nil, err
	}

	return &IFindByOrganizationIdUseCaseResponse{
		OrganizationLocations: dto.ConvertToOrganizationLocationResponse(organizationLocations),
	}, nil
}

func FindByOrganizationIdUseCaseFactory(log *logrus.Logger) IFindByOrganizationIdUseCase {
	repository := repository.OrganizationLocationRepositoryFactory(log)
	return NewFindByOrganizationIdUseCase(log, repository)
}
