package usecase

import (
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IFindAllPaginatedUseCaseRequest struct {
	Page        int      `json:"page"`
	PageSize    int      `json:"page_size"`
	Search      string   `json:"search"`
	IncludedIDs []string `json:"included_ids"`
}

type IFindAllPaginatedUseCaseResponse struct {
	OrganizationLocations *[]response.OrganizationLocationResponse `json:"organization_locations"`
	Total                 int64                                    `json:"total"`
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
	organizationLocations, total, err := uc.Repository.FindAllPaginated(req.Page, req.PageSize, req.Search, req.IncludedIDs)
	if err != nil {
		return nil, err
	}

	return &IFindAllPaginatedUseCaseResponse{
		OrganizationLocations: dto.ConvertToOrganizationLocationResponse(organizationLocations),
		Total:                 total,
	}, nil
}

func FindAllPaginatedUseCaseFactory(log *logrus.Logger) IFindAllPaginatedUseCase {
	repository := repository.OrganizationLocationRepositoryFactory(log)
	return NewFindAllPaginatedUseCase(log, repository)
}
