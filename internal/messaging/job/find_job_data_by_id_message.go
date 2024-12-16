package messaging

import (
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IFindJobDataByIdMessageRequest struct {
	JobID uuid.UUID `json:"job_id"`
}

type IFindJobDataByIdMessageResponse struct {
	JobID uuid.UUID             `json:"job_id"`
	Job   *response.JobResponse `json:"job"`
}

type IFindJobDataByIdMessage interface {
	Execute(request IFindJobDataByIdMessageRequest) (*IFindJobDataByIdMessageResponse, error)
}

type FindJobDataByIdMessage struct {
	Log           *logrus.Logger
	JobRepository repository.IJobRepository
}

func NewFindJobDataByIdMessage(log *logrus.Logger, jobRepository repository.IJobRepository) IFindJobDataByIdMessage {
	return &FindJobDataByIdMessage{
		Log:           log,
		JobRepository: jobRepository,
	}
}

func (m *FindJobDataByIdMessage) Execute(request IFindJobDataByIdMessageRequest) (*IFindJobDataByIdMessageResponse, error) {
	job, err := m.JobRepository.FindById(request.JobID)
	if err != nil {
		return nil, err
	}

	jobResponse := dto.ConvertToSingleJobResponse(job)

	return &IFindJobDataByIdMessageResponse{
		JobID: job.ID,
		Job:   jobResponse,
	}, nil
}

func FindJobDataByIdMessageFactory(log *logrus.Logger) IFindJobDataByIdMessage {
	repository := repository.JobRepositoryFactory(log)
	return NewFindJobDataByIdMessage(log, repository)
}
