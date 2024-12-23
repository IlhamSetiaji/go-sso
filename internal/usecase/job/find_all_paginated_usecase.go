package usecase

import (
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IFindAllPaginatedUseCaseRequest struct {
	Page           int    `json:"page"`
	PageSize       int    `json:"page_size"`
	Search         string `json:"search"`
	OrganizationID string `json:"organization_id"`
}

type IFindAllPaginatedUseCaseResponse struct {
	Jobs  *[]response.JobResponse `json:"jobs"`
	Total int64                   `json:"total"`
}

type IFindAllPaginatedUseCase interface {
	Execute(request *IFindAllPaginatedUseCaseRequest) (*IFindAllPaginatedUseCaseResponse, error)
}

type FindAllPaginatedUseCase struct {
	Log              *logrus.Logger
	JobRepository    repository.IJobRepository
	OrgStructureRepo repository.IOrganizationStructureRepository
}

func NewFindAllPaginatedUseCase(
	log *logrus.Logger,
	jobRepository repository.IJobRepository,
	orgStructureRepo repository.IOrganizationStructureRepository,
) IFindAllPaginatedUseCase {
	return &FindAllPaginatedUseCase{
		Log:              log,
		JobRepository:    jobRepository,
		OrgStructureRepo: orgStructureRepo,
	}
}

func (uc *FindAllPaginatedUseCase) Execute(req *IFindAllPaginatedUseCaseRequest) (*IFindAllPaginatedUseCaseResponse, error) {
	var includedIDs []string
	if req.OrganizationID != "" {
		orgStructures, err := uc.OrgStructureRepo.FindAllOrgStructuresByOrganizationID(uuid.MustParse(req.OrganizationID))
		if err != nil {
			return nil, err
		}

		for _, orgStructure := range *orgStructures {
			includedIDs = append(includedIDs, orgStructure.ID.String())
		}

	} else {
		includedIDs = []string{}
	}

	jobs, total, err := uc.JobRepository.FindAllPaginated(req.Page, req.PageSize, req.Search, includedIDs)
	if err != nil {
		return nil, err
	}

	for i, job := range *jobs {
		children, err := uc.JobRepository.FindAllChildren(job.ID)
		if err != nil {
			return nil, err
		}
		(*jobs)[i].Children = children

		if (*jobs)[i].ParentID != nil {
			parent, err := uc.JobRepository.FindParent(*(*jobs)[i].ParentID)
			if err != nil {
				return nil, err
			}

			if parent != nil {
				uc.Log.Info("parent", parent.ID)
				(*jobs)[i].Parent = parent
			} else {
				(*jobs)[i].Parent = nil
			}
		}
	}

	return &IFindAllPaginatedUseCaseResponse{
		Jobs:  dto.ConvertToJobResponse(jobs),
		Total: total,
	}, nil
}

func FindAllPaginatedUseCaseFactory(log *logrus.Logger) IFindAllPaginatedUseCase {
	jobRepository := repository.JobRepositoryFactory(log)
	orgStructureRepo := repository.OrganizationStructureRepositoryFactory(log)
	return NewFindAllPaginatedUseCase(log, jobRepository, orgStructureRepo)
}
