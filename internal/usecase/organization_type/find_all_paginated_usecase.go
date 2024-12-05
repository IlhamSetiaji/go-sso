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
	OrganizationTypes *[]entity.OrganizationType `json:"organization_types"`
	Total             int64                      `json:"total"`
}

type IFindAllPaginated interface {
	Execute(request *IFindAllPaginatedRequest) (*IFindAllPaginatedResponse, error)
}

type FindAllPaginated struct {
	Log                        *logrus.Logger
	OrganizationTypeRepository repository.IOrganizationTypeRepository
}

func NewFindAllPaginated(
	log *logrus.Logger,
	organizationTypeRepository repository.IOrganizationTypeRepository,
) IFindAllPaginated {
	return &FindAllPaginated{
		Log:                        log,
		OrganizationTypeRepository: organizationTypeRepository,
	}
}

func (uc *FindAllPaginated) Execute(req *IFindAllPaginatedRequest) (*IFindAllPaginatedResponse, error) {
	organizationTypes, total, err := uc.OrganizationTypeRepository.FindAllPaginated(req.Page, req.PageSize, req.Search)
	if err != nil {
		return nil, err
	}

	return &IFindAllPaginatedResponse{
		OrganizationTypes: organizationTypes,
		Total:             total,
	}, nil
}

func FindAllPaginatedUseCaseFactory(log *logrus.Logger) IFindAllPaginated {
	organizationTypeRepository := repository.OrganizationTypeRepositoryFactory(log)
	return NewFindAllPaginated(log, organizationTypeRepository)
}
