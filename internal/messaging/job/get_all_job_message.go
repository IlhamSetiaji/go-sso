package messaging

import (
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IGetAllJobMessageResponse struct {
	Jobs *[]response.JobResponse `json:"jobs"`
}

type IGetAllJobMessage interface {
	Execute() (*IGetAllJobMessageResponse, error)
}

type GetAllJobMessage struct {
	Log           *logrus.Logger
	JobRepository repository.IJobRepository
}

func NewGetAllJobMessage(
	log *logrus.Logger,
	jobRepository repository.IJobRepository,
) IGetAllJobMessage {
	return &GetAllJobMessage{
		Log:           log,
		JobRepository: jobRepository,
	}
}

func (uc *GetAllJobMessage) Execute() (*IGetAllJobMessageResponse, error) {
	jobs, err := uc.JobRepository.GetAll()
	if err != nil {
		return nil, err
	}

	return &IGetAllJobMessageResponse{
		Jobs: dto.ConvertToJobResponse(jobs),
	}, nil
}

func GetAllJobMessageFactory(log *logrus.Logger) IGetAllJobMessage {
	jobRepository := repository.JobRepositoryFactory(log)
	return NewGetAllJobMessage(log, jobRepository)
}
