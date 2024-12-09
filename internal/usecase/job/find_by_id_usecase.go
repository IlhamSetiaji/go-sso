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
	Job *response.JobResponse `json:"job"`
}

type IFindByIdUseCase interface {
	Execute(request *IFindByIdUseCaseRequest) (*IFindByIdUseCaseResponse, error)
}

type FindByIdUseCase struct {
	Log           *logrus.Logger
	JobRepository repository.IJobRepository
}

func NewFindByIdUseCase(
	log *logrus.Logger,
	jobRepository repository.IJobRepository,
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

	children, err := uc.JobRepository.FindAllChildren(job.ID)
	if err != nil {
		return nil, err
	}
	job.Children = children

	return &IFindByIdUseCaseResponse{
		Job: dto.ConvertToSingleJobResponse(job),
	}, nil
}

func FindByIdUseCaseFactory(log *logrus.Logger) IFindByIdUseCase {
	jobRepository := repository.JobRepositoryFactory(log)
	return NewFindByIdUseCase(log, jobRepository)
}
