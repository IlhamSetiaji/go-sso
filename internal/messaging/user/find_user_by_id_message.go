package messaging

import (
	"app/go-sso/internal/repository"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IFindUserByIDMessageRequest struct {
	UserID uuid.UUID `json:"user_id"`
}

type IFindUserByIDMessageResponse struct {
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
}

type IFindUserByIDMessage interface {
	Execute(request IFindUserByIDMessageRequest) (*IFindUserByIDMessageResponse, error)
}

type FindUserByIDMessage struct {
	Log        *logrus.Logger
	Repository repository.IUserRepository
}

func NewFindUserByIDMessage(log *logrus.Logger, repository repository.IUserRepository) IFindUserByIDMessage {
	return &FindUserByIDMessage{
		Log:        log,
		Repository: repository,
	}
}

func (m *FindUserByIDMessage) Execute(request IFindUserByIDMessageRequest) (*IFindUserByIDMessageResponse, error) {
	user, err := m.Repository.FindById(request.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return &IFindUserByIDMessageResponse{
		UserID: user.ID,
		Name:   user.Name,
	}, nil
}

func FindUserByIDMessageFactory(log *logrus.Logger) IFindUserByIDMessage {
	repository := repository.UserRepositoryFactory(log)
	return NewFindUserByIDMessage(log, repository)
}
