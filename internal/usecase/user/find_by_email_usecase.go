package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"
	"errors"

	"github.com/sirupsen/logrus"
)

type FindByEmailUseCase struct {
	Log            *logrus.Logger
	UserRepository repository.UserRepositoryInterface
}

type FindByEmailUseCaseResponse struct {
	User *entity.User `json:"user"`
}

func FindByEmailUseCaseFactory(
	log *logrus.Logger,
) *FindByEmailUseCase {
	return &FindByEmailUseCase{
		Log:            log,
		UserRepository: repository.UserRepositoryFactory(log),
	}
}

func (uc *FindByEmailUseCase) FindByEmail(email string) (*FindByEmailUseCaseResponse, error) {
	user, err := uc.UserRepository.FindByEmail(email)
	if err != nil {
		return nil, errors.New("[FindByEmailUseCase.FindByEmail] " + err.Error())
	}

	if user == nil {
		uc.Log.Panicf("User not found")
		return nil, errors.New("User not found")
	}

	return &FindByEmailUseCaseResponse{
		User: user,
	}, nil
}
