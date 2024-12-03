package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IFindByIdUseCaseRequest struct {
	ID uuid.UUID `json:"id"`
}

type IFindByIdUseCaseResponse struct {
	User *entity.User `json:"user"`
}

type IFindByIdUseCase interface {
	Execute(request *IFindByIdUseCaseRequest) (*IFindByIdUseCaseResponse, error)
}

type FindByIdUseCase struct {
	Log            *logrus.Logger
	UserRepository repository.IUserRepository
}

func NewFindByIdUseCase(log *logrus.Logger, userRepository repository.IUserRepository) IFindByIdUseCase {
	return &FindByIdUseCase{
		Log:            log,
		UserRepository: userRepository,
	}
}

func (u *FindByIdUseCase) Execute(request *IFindByIdUseCaseRequest) (*IFindByIdUseCaseResponse, error) {
	user, err := u.UserRepository.FindById(request.ID)
	if err != nil {
		u.Log.Error("[FindByIdUseCase.Execute] " + err.Error())
		return nil, err
	}

	return &IFindByIdUseCaseResponse{
		User: user,
	}, nil
}

func FindByIdUseCaseFactory(log *logrus.Logger) IFindByIdUseCase {
	userRepository := repository.UserRepositoryFactory(log)
	return NewFindByIdUseCase(log, userRepository)
}
