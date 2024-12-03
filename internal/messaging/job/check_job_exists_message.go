package messaging

import (
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ICheckJobExistMessageRequest struct {
	JobID uuid.UUID `json:"job_id"`
}

type ICheckJobExistMessageResponse struct {
	JobID  uuid.UUID `json:"job_id"`
	Exists bool      `json:"exists"`
}

type ICheckJobExistMessage interface {
	Execute(request ICheckJobExistMessageRequest) (*ICheckJobExistMessageResponse, error)
}

type CheckJobExistMessage struct {
	Log        *logrus.Logger
	Repository repository.IJobRepository
}

func NewCheckJobExistMessage(log *logrus.Logger, repository repository.IJobRepository) ICheckJobExistMessage {
	return &CheckJobExistMessage{
		Log:        log,
		Repository: repository,
	}
}

func (m *CheckJobExistMessage) Execute(request ICheckJobExistMessageRequest) (*ICheckJobExistMessageResponse, error) {
	job, err := m.Repository.FindById(request.JobID)
	if err != nil {
		return nil, err
	}

	exists := job != nil

	return &ICheckJobExistMessageResponse{
		JobID:  request.JobID,
		Exists: exists,
	}, nil
}

func CheckJobExistMessageFactory(log *logrus.Logger) ICheckJobExistMessage {
	repository := repository.JobRepositoryFactory(log)
	return NewCheckJobExistMessage(log, repository)
}
