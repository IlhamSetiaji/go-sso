package usecase

import (
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IDeleteUserUseCaseRequest struct {
	ID uuid.UUID `json:"id"`
}

type IDeleteUserUseCase interface {
	Execute(request IDeleteUserUseCaseRequest) error
}

type DeleteUserUseCase struct {
	Log            *logrus.Logger
	userRepository repository.IUserRepository
}

func NewDeleteUserUseCase(log *logrus.Logger, userRepository repository.IUserRepository) IDeleteUserUseCase {
	return &DeleteUserUseCase{
		Log:            log,
		userRepository: userRepository,
	}
}

func (uc *DeleteUserUseCase) Execute(request IDeleteUserUseCaseRequest) error {
	uc.Log.Info("Delete user usecase")

	err := uc.userRepository.DeleteUser(request.ID)

	if err != nil {
		uc.Log.Error("Delete user usecase error: " + err.Error())
		return err
	}

	return nil
}

func DeleteUserUseCaseFactory(log *logrus.Logger) IDeleteUserUseCase {
	userRepository := repository.UserRepositoryFactory(log)
	return &DeleteUserUseCase{
		Log:            log,
		userRepository: userRepository,
	}
}
