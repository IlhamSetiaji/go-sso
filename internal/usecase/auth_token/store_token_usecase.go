package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IStoreTokenUseCaseRequest struct {
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}

type IStoreTokenUseCaseResponse struct {
	AuthToken entity.AuthToken `json:"auth_token"`
}

type IStoreTokenUseCase interface {
	Execute(request IStoreTokenUseCaseRequest) (*IStoreTokenUseCaseResponse, error)
}

type StoreTokenUseCase struct {
	Log                 *logrus.Logger
	AuthTokenRepository repository.IAuthTokenRepository
	UserRepository      repository.IUserRepository
}

func NewStoreTokenUseCase(log *logrus.Logger, authTokenRepository repository.IAuthTokenRepository, userRepository repository.IUserRepository) IStoreTokenUseCase {
	return &StoreTokenUseCase{
		Log:                 log,
		AuthTokenRepository: authTokenRepository,
		UserRepository:      userRepository,
	}
}

func (uc *StoreTokenUseCase) Execute(request IStoreTokenUseCaseRequest) (*IStoreTokenUseCaseResponse, error) {
	authToken := entity.AuthToken{
		UserID:    request.UserID,
		Token:     request.Token,
		ExpiredAt: request.ExpiredAt,
	}

	user, err := uc.UserRepository.FindById(request.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("[StoreTokenUseCase.Execute] User not found")
	}

	if err := uc.AuthTokenRepository.StoreAuthToken(user, &authToken); err != nil {
		return nil, err
	}

	return &IStoreTokenUseCaseResponse{
		AuthToken: authToken,
	}, nil
}

func StoreTokenUseCaseFactory(log *logrus.Logger) IStoreTokenUseCase {
	authTokenRepository := repository.AuthTokenRepositoryFactory(log)
	userRepository := repository.UserRepositoryFactory(log)
	return NewStoreTokenUseCase(log, authTokenRepository, userRepository)
}
