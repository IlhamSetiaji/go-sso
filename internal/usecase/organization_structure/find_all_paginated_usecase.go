package usecase

import (
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IFindAllPaginatedUseCaseRequest struct {
	Page           int
	PageSize       int
	Search         string
	OrganizationID uuid.UUID
	Filter         map[string]interface{}
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
	organizationStructures, total, err := uc.Repository.FindAllPaginatedByOrganizationID(request.OrganizationID, request.Page, request.PageSize, request.Search, request.Filter)
	if err != nil {
		return nil, err
	}

	for i, orgStructure := range *organizationStructures {
		children, err := uc.Repository.FindAllChildren(orgStructure.ID)
		if err != nil {
			return nil, err
		}
		(*organizationStructures)[i].Children = children

		if (*organizationStructures)[i].ParentID != nil {
			parent, err := uc.Repository.FindParent(*(*organizationStructures)[i].ParentID)
			if err != nil {
				return nil, err
			}

			if parent != nil {
				(*organizationStructures)[i].Parent = parent
			} else {
				(*organizationStructures)[i].Parent = nil
			}
		}
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
