package messaging

import (
	"app/go-sso/internal/repository"
	"strconv"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IFindJobLevelByIDMessageRequest struct {
	JobLevelID uuid.UUID `json:"job_level_id"`
}

type IFindJobLevelByIDMessageResponse struct {
	JobLevelID string `json:"job_level_id"`
	Name       string `json:"name"`
	Level      int    `json:"level"`
}

type IFindJobLevelByIDMessage interface {
	Execute(request IFindJobLevelByIDMessageRequest) (*IFindJobLevelByIDMessageResponse, error)
}

type FindJobLevelByIDMessage struct {
	Log        *logrus.Logger
	Repository repository.IJobLevelRepository
}

func NewFindJobLevelByIDMessage(log *logrus.Logger, repository repository.IJobLevelRepository) IFindJobLevelByIDMessage {
	return &FindJobLevelByIDMessage{
		Log:        log,
		Repository: repository,
	}
}

func (m *FindJobLevelByIDMessage) Execute(request IFindJobLevelByIDMessageRequest) (*IFindJobLevelByIDMessageResponse, error) {
	jobLevel, err := m.Repository.FindById(request.JobLevelID)
	if err != nil {
		return nil, err
	}

	level, err := strconv.Atoi(jobLevel.Level)
	if err != nil {
		return nil, err
	}

	return &IFindJobLevelByIDMessageResponse{
		JobLevelID: jobLevel.ID.String(),
		Name:       jobLevel.Name,
		Level:      level,
	}, nil
}

func FindJobLevelByIDMessageFactory(log *logrus.Logger) IFindJobLevelByIDMessage {
	repository := repository.JobLevelRepositoryFactory(log)
	return NewFindJobLevelByIDMessage(log, repository)
}
