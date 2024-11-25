package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IFindTokenUseCaseRequest struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

type IFindTokenUseCaseResponse struct {
	AuthToken *entity.AuthToken `json:"auth_token"`
}

type IFindTokenUseCase interface {
	Execute(request IFindTokenUseCaseRequest) (*IFindTokenUseCaseResponse, error)
}

type FindTokenUseCase struct {
	Log                 *logrus.Logger
	AuthTokenRepository repository.IAuthTokenRepository
}

func NewFindTokenUseCase(log *logrus.Logger, authTokenRepository repository.IAuthTokenRepository) IFindTokenUseCase {
	return &FindTokenUseCase{
		Log:                 log,
		AuthTokenRepository: authTokenRepository,
	}
}

func (uc *FindTokenUseCase) Execute(request IFindTokenUseCaseRequest) (*IFindTokenUseCaseResponse, error) {
	authToken, err := uc.AuthTokenRepository.FindAuthToken(request.UserID, request.Token)
	if err != nil {
		return nil, err
	}

	return &IFindTokenUseCaseResponse{
		AuthToken: authToken,
	}, nil
}

func FindTokenUseCaseFactory(log *logrus.Logger) IFindTokenUseCase {
	authTokenRepository := repository.AuthTokenRepositoryFactory(log)
	return NewFindTokenUseCase(log, authTokenRepository)
}
