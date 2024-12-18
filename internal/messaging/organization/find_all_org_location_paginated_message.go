package messaging

import (
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IFindAllOrgLocationPaginatedRequest struct {
	Page        int      `json:"page"`
	PageSize    int      `json:"page_size"`
	Search      string   `json:"search"`
	IncludedIDs []string `json:"included_ids"`
}

type IFindAllOrgLocationPaginatedResponse struct {
	OrganizationLocations *[]response.OrganizationLocationResponse `json:"organization_locations"`
	Total                 int64                                    `json:"total"`
}

type IFindAllOrgLocationPaginated interface {
	Execute(request *IFindAllOrgLocationPaginatedRequest) (*IFindAllOrgLocationPaginatedResponse, error)
}

type FindAllOrgLocationPaginated struct {
	Log                            *logrus.Logger
	OrganizationLocationRepository repository.IOrganizationLocationRepository
}

func NewFindAllOrgLocationPaginated(
	log *logrus.Logger,
	organizationLocationRepository repository.IOrganizationLocationRepository,
) IFindAllOrgLocationPaginated {
	return &FindAllOrgLocationPaginated{
		Log:                            log,
		OrganizationLocationRepository: organizationLocationRepository,
	}
}

func (uc *FindAllOrgLocationPaginated) Execute(req *IFindAllOrgLocationPaginatedRequest) (*IFindAllOrgLocationPaginatedResponse, error) {
	organizationLocations, total, err := uc.OrganizationLocationRepository.FindAllPaginated(req.Page, req.PageSize, req.Search, req.IncludedIDs)
	if err != nil {
		return nil, err
	}

	return &IFindAllOrgLocationPaginatedResponse{
		OrganizationLocations: dto.ConvertToOrganizationLocationResponse(organizationLocations),
		Total:                 total,
	}, nil
}

func FindAllOrgLocationPaginatedUseCaseFactory(log *logrus.Logger) IFindAllOrgLocationPaginated {
	organizationLocationRepository := repository.OrganizationLocationRepositoryFactory(log)
	return NewFindAllOrgLocationPaginated(log, organizationLocationRepository)
}
