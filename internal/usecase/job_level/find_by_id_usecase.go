package usecase

import (
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IFindByIdUseCaseRequest struct {
	ID uuid.UUID `json:"id"`
}

type IFindByIdUseCaseResponse struct {
	Job *response.JobLevelResponse `json:"job"`
}

type IFindByIdUseCase interface {
	Execute(request *IFindByIdUseCaseRequest) (*IFindByIdUseCaseResponse, error)
}

type FindByIdUseCase struct {
	Log           *logrus.Logger
	JobRepository repository.IJobLevelRepository
}

func NewFindByIdUseCase(
	log *logrus.Logger,
	jobRepository repository.IJobLevelRepository,
) IFindByIdUseCase {
	return &FindByIdUseCase{
		Log:           log,
		JobRepository: jobRepository,
	}
}

func (uc *FindByIdUseCase) Execute(req *IFindByIdUseCaseRequest) (*IFindByIdUseCaseResponse, error) {
	job, err := uc.JobRepository.FindById(req.ID)
	if err != nil {
		return nil, err
	}

	return &IFindByIdUseCaseResponse{
		Job: dto.ConvertToSingleJobLevelResponse(job),
	}, nil
}

func FindByIdUseCaseFactory(log *logrus.Logger) IFindByIdUseCase {
	jobRepository := repository.JobLevelRepositoryFactory(log)
	return NewFindByIdUseCase(log, jobRepository)
}
