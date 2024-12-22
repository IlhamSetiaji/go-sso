package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IFindByIdOnlyUseCaseRequest struct {
	ID uuid.UUID `json:"id"`
}

type IFindByIdOnlyUseCaseResponse struct {
	User *entity.User `json:"user"`
}

type IFindByIdOnlyUseCase interface {
	Execute(request *IFindByIdOnlyUseCaseRequest) (*IFindByIdOnlyUseCaseResponse, error)
}

type FindByIdOnlyUseCase struct {
	Log            *logrus.Logger
	UserRepository repository.IUserRepository
}

func NewFindByIdOnlyUseCase(log *logrus.Logger, userRepository repository.IUserRepository) IFindByIdOnlyUseCase {
	return &FindByIdOnlyUseCase{
		Log:            log,
		UserRepository: userRepository,
	}
}

func (u *FindByIdOnlyUseCase) Execute(request *IFindByIdOnlyUseCaseRequest) (*IFindByIdOnlyUseCaseResponse, error) {
	user, err := u.UserRepository.FindByIdOnly(request.ID)
	if err != nil {
		u.Log.Error("[FindByIdOnlyUseCase.Execute] " + err.Error())
		return nil, err
	}

	return &IFindByIdOnlyUseCaseResponse{
		User: user,
	}, nil
}

func FindByIdOnlyUseCaseFactory(log *logrus.Logger) IFindByIdOnlyUseCase {
	userRepository := repository.UserRepositoryFactory(log)
	return NewFindByIdOnlyUseCase(log, userRepository)
}
