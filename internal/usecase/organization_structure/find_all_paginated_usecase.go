package usecase

import (
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IFindAllPaginatedUseCaseRequest struct {
	Page     int
	PageSize int
	Search   string
}

type IFindAllPaginatedResponse struct {
	OrganizationStructures *[]response.OrganizationStructureResponse
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

	for i, orgStructure := range *organizationStructures {
		children, err := uc.Repository.FindAllChildren(orgStructure.ID)
		if err != nil {
			return nil, err
		}
		(*organizationStructures)[i].Children = children
	}

	return &IFindAllPaginatedResponse{
		OrganizationStructures: dto.ConvertToOrganizationStructureResponse(organizationStructures),
		Total:                  total,
	}, nil
}

func FindAllPaginatedUseCaseFactory(log *logrus.Logger) IFindAllPaginatedUseCase {
	repository := repository.OrganizationStructureRepositoryFactory(log)
	return NewFindAllPaginatedUseCase(log, repository)
}
