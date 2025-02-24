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
	JobLevels *[]response.JobLevelResponse `json:"job_levels"`
	Total     int64                        `json:"total"`
}

type IFindAllPaginatedUseCase interface {
	Execute(request *IFindAllPaginatedUseCaseRequest) (*IFindAllPaginatedUseCaseResponse, error)
}

type FindAllPaginatedUseCase struct {
	Log                *logrus.Logger
	JobLevelRepository repository.IJobLevelRepository
}

func NewFindAllPaginatedUseCase(
	log *logrus.Logger,
	jobLevelRepository repository.IJobLevelRepository,
) IFindAllPaginatedUseCase {
	return &FindAllPaginatedUseCase{
		Log:                log,
		JobLevelRepository: jobLevelRepository,
	}
}

func (uc *FindAllPaginatedUseCase) Execute(req *IFindAllPaginatedUseCaseRequest) (*IFindAllPaginatedUseCaseResponse, error) {
	jobs, total, err := uc.JobLevelRepository.FindAllPaginated(req.Page, req.PageSize, req.Search)
	if err != nil {
		return nil, err
	}

	return &IFindAllPaginatedUseCaseResponse{
		JobLevels: dto.ConvertToJobLevelResponse(jobs),
		Total:     total,
	}, nil
}

func FindAllPaginatedUseCaseFactory(log *logrus.Logger) IFindAllPaginatedUseCase {
	jobLevelRepository := repository.JobLevelRepositoryFactory(log)
	return NewFindAllPaginatedUseCase(log, jobLevelRepository)
}
