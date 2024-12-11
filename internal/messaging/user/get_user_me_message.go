package messaging

import (
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IGetUserMeMessageRequest struct {
	UserID uuid.UUID `json:"user_id"`
}

type IGetUserMeMessageResponse struct {
	User response.UserResponse `json:"user"`
}

type IGetUserMeMessage interface {
	Execute(request IGetUserMeMessageRequest) (*IGetUserMeMessageResponse, error)
}

type GetUserMeMessage struct {
	Log            *logrus.Logger
	UserRepository repository.IUserRepository
}

func NewGetUserMeMessage(log *logrus.Logger, userRepository repository.IUserRepository) IGetUserMeMessage {
	return &GetUserMeMessage{
		Log:            log,
		UserRepository: userRepository,
	}
}

func (m *GetUserMeMessage) Execute(request IGetUserMeMessageRequest) (*IGetUserMeMessageResponse, error) {
	user, err := m.UserRepository.FindById(request.UserID)
	if err != nil {
		return nil, err
	}

	return &IGetUserMeMessageResponse{
		User: *dto.ConvertToSingleUserResponse(user),
	}, nil
}

func GetUserMeMessageFactory(log *logrus.Logger) IGetUserMeMessage {
	repository := repository.UserRepositoryFactory(log)
	return NewGetUserMeMessage(log, repository)
}
