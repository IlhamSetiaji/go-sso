package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IFindAllPaginatedUseCaseRequest struct {
	Page     int
	PageSize int
	Search   string
}

type IFindAllPaginatedResponse struct {
	OrganizationStructures *[]entity.OrganizationStructure
	Total                  int64
}

type FindAllPaginatedUseCase struct {
	Log        *logrus.Logger
	Repository repository.IOrganizationStructureRepository
}

type IFindAllPaginatedUseCase interface {
	Execute(request *IFindAllPaginatedUseCaseRequest) (*IFindAllPaginatedResponse, error)
}

func NewFindAllPaginatedUseCase(log *logrus.Logger, repository repository.IOrganizationStructureRepository) IFindAllPaginatedUseCase {
	return &FindAllPaginatedUseCase{
		Log:        log,
		Repository: repository,
	}
}

func (uc *FindAllPaginatedUseCase) Execute(request *IFindAllPaginatedUseCaseRequest) (*IFindAllPaginatedResponse, error) {
	organizationStructures, total, err := uc.Repository.FindAllPaginated(request.Page, request.PageSize, request.Search)
	if err != nil {
		return nil, err
	}

	return &IFindAllPaginatedResponse{
		OrganizationStructures: organizationStructures,
		Total:                  total,
	}, nil
}

func FindAllPaginatedUseCaseFactory(log *logrus.Logger) IFindAllPaginatedUseCase {
	repository := repository.OrganizationStructureRepositoryFactory(log)
	return NewFindAllPaginatedUseCase(log, repository)
}
