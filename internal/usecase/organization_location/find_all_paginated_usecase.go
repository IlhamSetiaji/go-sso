package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IFindAllPaginatedUseCaseRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Search   string `json:"search"`
}

type IFindAllPaginatedUseCaseResponse struct {
	OrganizationLocations *[]entity.OrganizationLocation `json:"organization_locations"`
	Total                 int64                          `json:"total"`
}

type IFindAllPaginatedUseCase interface {
	Execute(request *IFindAllPaginatedUseCaseRequest) (*IFindAllPaginatedUseCaseResponse, error)
}

type FindAllPaginatedUseCase struct {
	Log        *logrus.Logger
	Repository repository.IOrganizationLocationRepository
}

func NewFindAllPaginatedUseCase(
	log *logrus.Logger,
	repository repository.IOrganizationLocationRepository,
) IFindAllPaginatedUseCase {
	return &FindAllPaginatedUseCase{
		Log:        log,
		Repository: repository,
	}
}

func (uc *FindAllPaginatedUseCase) Execute(req *IFindAllPaginatedUseCaseRequest) (*IFindAllPaginatedUseCaseResponse, error) {
	organizationLocations, total, err := uc.Repository.FindAllPaginated(req.Page, req.PageSize, req.Search)
	if err != nil {
		return nil, err
	}

	return &IFindAllPaginatedUseCaseResponse{
		OrganizationLocations: organizationLocations,
		Total:                 total,
	}, nil
}

func FindAllPaginatedUseCaseFactory(log *logrus.Logger) IFindAllPaginatedUseCase {
	repository := repository.OrganizationLocationRepositoryFactory(log)
	return NewFindAllPaginatedUseCase(log, repository)
}
