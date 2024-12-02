package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IGetAllUsersUseCaseResponse struct {
	Users *[]entity.User `json:"users"`
}

type IGetAllUsersUseCase interface {
	Execute() (*IGetAllUsersUseCaseResponse, error)
}

type GetAllUsersUseCase struct {
	Log            *logrus.Logger
	UserRepository repository.IUserRepository
}

func NewGetAllUsersUseCase(log *logrus.Logger, userRepository repository.IUserRepository) IGetAllUsersUseCase {
	return &GetAllUsersUseCase{
		Log:            log,
		UserRepository: userRepository,
	}
}

func (uc *GetAllUsersUseCase) Execute() (*IGetAllUsersUseCaseResponse, error) {
	users, err := uc.UserRepository.GetAllUsers()
	if err != nil {
		uc.Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to get all users")
		return nil, err
	}

	return &IGetAllUsersUseCaseResponse{
		Users: users,
	}, nil
}

func GetAllUsersUseCaseFactory(log *logrus.Logger) IGetAllUsersUseCase {
	userRepository := repository.UserRepositoryFactory(log)
	return NewGetAllUsersUseCase(log, userRepository)
}
