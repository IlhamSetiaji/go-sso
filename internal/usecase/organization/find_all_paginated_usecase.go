package usecase

import (
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IFindAllPaginatedRequest struct {
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
	Search   string                 `json:"search"`
	Filter   map[string]interface{} `json:"filter"`
}

type IFindAllPaginatedResponse struct {
	Organizations *[]response.OrganizationResponse `json:"organizations"`
	Total         int64                            `json:"total"`
}

type IFindAllPaginated interface {
	Execute(request *IFindAllPaginatedRequest) (*IFindAllPaginatedResponse, error)
}

type FindAllPaginated struct {
	Log                    *logrus.Logger
	OrganizationRepository repository.IOrganizationRepository
}

func NewFindAllPaginated(
	log *logrus.Logger,
	organizationRepository repository.IOrganizationRepository,
) IFindAllPaginated {
	return &FindAllPaginated{
		Log:                    log,
		OrganizationRepository: organizationRepository,
	}
}

func (uc *FindAllPaginated) Execute(req *IFindAllPaginatedRequest) (*IFindAllPaginatedResponse, error) {
	organizations, total, err := uc.OrganizationRepository.FindAllPaginated(req.Page, req.PageSize, req.Search, req.Filter)
	if err != nil {
		return nil, err
	}

	return &IFindAllPaginatedResponse{
		Organizations: dto.ConvertToOrganizationResponse(organizations),
		Total:         total,
	}, nil
}

func FindAllPaginatedUseCaseFactory(log *logrus.Logger) IFindAllPaginated {
	organizationRepository := repository.OrganizationRepositoryFactory(log)
	return NewFindAllPaginated(log, organizationRepository)
}
