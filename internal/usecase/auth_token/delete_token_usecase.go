package usecase

import (
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IDeleteTokenUseCaseRequest struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

type IDeleteTokenUseCaseResponse struct {
	Message string `json:"message"`
}

type IDeleteTokenUseCase interface {
	Execute(request IDeleteTokenUseCaseRequest) (*IDeleteTokenUseCaseResponse, error)
}

type DeleteTokenUseCase struct {
	Log                 *logrus.Logger
	AuthTokenRepository repository.IAuthTokenRepository
}

func NewDeleteTokenUseCase(log *logrus.Logger, authTokenRepository repository.IAuthTokenRepository) IDeleteTokenUseCase {
	return &DeleteTokenUseCase{
		Log:                 log,
		AuthTokenRepository: authTokenRepository,
	}
}

func (uc *DeleteTokenUseCase) Execute(request IDeleteTokenUseCaseRequest) (*IDeleteTokenUseCaseResponse, error) {
	if err := uc.AuthTokenRepository.DeleteAuthToken(request.UserID, request.Token); err != nil {
		return nil, err
	}

	return &IDeleteTokenUseCaseResponse{
		Message: "Token deleted",
	}, nil
}

func DeleteTokenUseCaseFactory(log *logrus.Logger) IDeleteTokenUseCase {
	authTokenRepository := repository.AuthTokenRepositoryFactory(log)
	return NewDeleteTokenUseCase(log, authTokenRepository)
}
