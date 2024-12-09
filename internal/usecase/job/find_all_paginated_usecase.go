package usecase

import (
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IFindAllPaginatedUseCaseRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Search   string `json:"search"`
}

type IFindAllPaginatedUseCaseResponse struct {
	Jobs  *[]response.JobResponse `json:"jobs"`
	Total int64                   `json:"total"`
}

type IFindAllPaginatedUseCase interface {
	Execute(request *IFindAllPaginatedUseCaseRequest) (*IFindAllPaginatedUseCaseResponse, error)
}

type FindAllPaginatedUseCase struct {
	Log           *logrus.Logger
	JobRepository repository.IJobRepository
}

func NewFindAllPaginatedUseCase(
	log *logrus.Logger,
	jobRepository repository.IJobRepository,
) IFindAllPaginatedUseCase {
	return &FindAllPaginatedUseCase{
		Log:           log,
		JobRepository: jobRepository,
	}
}

func (uc *FindAllPaginatedUseCase) Execute(req *IFindAllPaginatedUseCaseRequest) (*IFindAllPaginatedUseCaseResponse, error) {
	jobs, total, err := uc.JobRepository.FindAllPaginated(req.Page, req.PageSize, req.Search)
	if err != nil {
		return nil, err
	}

	for i, job := range *jobs {
		children, err := uc.JobRepository.FindAllChildren(job.ID)
		if err != nil {
			return nil, err
		}
		(*jobs)[i].Children = children
	}

	return &IFindAllPaginatedUseCaseResponse{
		Jobs:  dto.ConvertToJobResponse(jobs),
		Total: total,
	}, nil
}

func FindAllPaginatedUseCaseFactory(log *logrus.Logger) IFindAllPaginatedUseCase {
	jobRepository := repository.JobRepositoryFactory(log)
	return NewFindAllPaginatedUseCase(log, jobRepository)
}
