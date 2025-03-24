package messaging

import (
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IFindJobByIDMessageRequest struct {
	JobID uuid.UUID `json:"job_id"`
}

type IFindJobByIDMessageResponse struct {
	JobID     uuid.UUID `json:"job_id"`
	Name      string    `json:"name"`
	MidsuitID string    `json:"midsuit_id"`
}

type IFindJobByIDMessage interface {
	Execute(request IFindJobByIDMessageRequest) (*IFindJobByIDMessageResponse, error)
}

type FindJobByIDMessage struct {
	Log        *logrus.Logger
	Repository repository.IJobRepository
}

func NewFindJobByIDMessage(log *logrus.Logger, repository repository.IJobRepository) IFindJobByIDMessage {
	return &FindJobByIDMessage{
		Log:        log,
		Repository: repository,
	}
}

func (m *FindJobByIDMessage) Execute(request IFindJobByIDMessageRequest) (*IFindJobByIDMessageResponse, error) {
	job, err := m.Repository.FindById(request.JobID)
	if err != nil {
		return nil, err
	}

	return &IFindJobByIDMessageResponse{
		JobID:     job.ID,
		Name:      job.Name,
		MidsuitID: job.MidsuitID,
	}, nil
}

func FindJobByIDMessageFactory(log *logrus.Logger) IFindJobByIDMessage {
	repository := repository.JobRepositoryFactory(log)
	return NewFindJobByIDMessage(log, repository)
}
