package messaging

import (
	"app/go-sso/internal/repository"
	"errors"

	"github.com/sirupsen/logrus"
)

type IGetUsersByPermissionNamesMessageRequest struct {
	PermissionNames []string `json:"permission_names"`
}

type IGetUsersByPermissionNamesMessageResponse struct {
	UserIDs []string `json:"user_ids"`
}

type IGetUsersByPermissionNamesMessage interface {
	Execute(request IGetUsersByPermissionNamesMessageRequest) (*IGetUsersByPermissionNamesMessageResponse, error)
}

type GetUsersByPermissionNamesMessage struct {
	Log        *logrus.Logger
	Repository repository.IUserRepository
}

func NewGetUsersByPermissionNamesMessage(log *logrus.Logger, repository repository.IUserRepository) IGetUsersByPermissionNamesMessage {
	return &GetUsersByPermissionNamesMessage{
		Log:        log,
		Repository: repository,
	}
}

func (m *GetUsersByPermissionNamesMessage) Execute(request IGetUsersByPermissionNamesMessageRequest) (*IGetUsersByPermissionNamesMessageResponse, error) {
	users, err := m.Repository.GetAllUsersByPermissionNames(request.PermissionNames)
	if err != nil {
		return nil, err
	}

	if len(*users) == 0 {
		return nil, errors.New("no users found")
	}

	userIDs := make([]string, len(*users))
	for i, user := range *users {
		userIDs[i] = user.ID.String()
	}

	return &IGetUsersByPermissionNamesMessageResponse{
		UserIDs: userIDs,
	}, nil
}
