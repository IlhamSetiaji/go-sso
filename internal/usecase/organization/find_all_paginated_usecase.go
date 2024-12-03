package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IFindAllPaginatedRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Search   string `json:"search"`
}

type IFindAllPaginatedResponse struct {
	Organizations *[]entity.Organization `json:"organizations"`
	Total         int64                  `json:"total"`
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
	organizations, total, err := uc.OrganizationRepository.FindAllPaginated(req.Page, req.PageSize, req.Search)
	if err != nil {
		return nil, err
	}

	return &IFindAllPaginatedResponse{
		Organizations: organizations,
		Total:         total,
	}, nil
}

func FindAllPaginatedUseCaseFactory(log *logrus.Logger) IFindAllPaginated {
	organizationRepository := repository.OrganizationRepositoryFactory(log)
	return NewFindAllPaginated(log, organizationRepository)
}
