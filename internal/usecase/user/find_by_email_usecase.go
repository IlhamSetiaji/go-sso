package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"
	"errors"

	"github.com/sirupsen/logrus"
)

type IFindByEmailUseCaseRequest struct {
	Email string `json:"email"`
}

type IFindByEmailUseCaseResponse struct {
	User *entity.User `json:"user"`
}

type IFindByEmailUseCase interface {
	Execute(request IFindByEmailUseCaseRequest) (*IFindByEmailUseCaseResponse, error)
}

type FindByEmailUseCase struct {
	Log            *logrus.Logger
	UserRepository repository.IUserRepository
}

func NewFindByEmailUseCase(log *logrus.Logger, userRepository repository.IUserRepository) IFindByEmailUseCase {
	return &FindByEmailUseCase{
		Log:            log,
		UserRepository: userRepository,
	}
}

func (uc *FindByEmailUseCase) Execute(request IFindByEmailUseCaseRequest) (*IFindByEmailUseCaseResponse, error) {
	user, err := uc.UserRepository.FindByEmail(request.Email)
	if err != nil {
		return nil, errors.New("[FindByEmailUseCase.FindByEmail] " + err.Error())
	}

	if user == nil {
		uc.Log.Panicf("User not found")
		return nil, errors.New("User not found")
	}

	return &IFindByEmailUseCaseResponse{
		User: user,
	}, nil
}

func FindByEmailUseCaseFactory(log *logrus.Logger) IFindByEmailUseCase {
	userRepository := repository.UserRepositoryFactory(log)
	return NewFindByEmailUseCase(log, userRepository)
}
