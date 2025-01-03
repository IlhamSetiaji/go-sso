package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IUpdateUserUseCaseRequest struct {
	User    *entity.User `json:"user"`
	RoleIDs []string     `json:"role_ids,omitempty"`
}

type IUpdateUserUseCaseResponse struct {
	User *entity.User `json:"user"`
}

type IUpdateUserUseCase interface {
	Execute(request IUpdateUserUseCaseRequest) (IUpdateUserUseCaseResponse, error)
}

type UpdateUserUseCase struct {
	Log            *logrus.Logger
	userRepository repository.IUserRepository
}

func NewUpdateUserUseCase(log *logrus.Logger, userRepository repository.IUserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		Log:            log,
		userRepository: userRepository,
	}
}

func (uc *UpdateUserUseCase) Execute(request IUpdateUserUseCaseRequest) (IUpdateUserUseCaseResponse, error) {
	uc.Log.Info("Update user usecase")
	userExist, err := uc.userRepository.FindById(request.User.ID)
	if err != nil {
		uc.Log.Error("Update user usecase: " + err.Error())
		return IUpdateUserUseCaseResponse{}, errors.New("[UpdateUserUseCase] user not found: " + err.Error())
	}

	if userExist == nil {
		uc.Log.Error("Update user usecase: user not found")
		return IUpdateUserUseCaseResponse{}, errors.New("[UpdateUserUseCase] user not found")
	}

	var roleUUIDs []uuid.UUID
	for _, roleID := range request.RoleIDs {
		roleUUID, err := uuid.Parse(roleID)
		if err != nil {
			uc.Log.Error("Update user usecase: " + err.Error())
			return IUpdateUserUseCaseResponse{}, errors.New("[UpdateUserUseCase] invalid role ID: " + err.Error())
		}
		roleUUIDs = append(roleUUIDs, roleUUID)
	}

	user, err := uc.userRepository.UpdateUser(request.User, roleUUIDs)
	if err != nil {
		uc.Log.Error("Update user usecase: " + err.Error())
		return IUpdateUserUseCaseResponse{}, errors.New("[UpdateUserUseCase] error update user: " + err.Error())
	}

	return IUpdateUserUseCaseResponse{
		User: user,
	}, nil
}

func UpdateUserUseCaseFactory(log *logrus.Logger) *UpdateUserUseCase {
	userRepository := repository.UserRepositoryFactory(log)
	return NewUpdateUserUseCase(log, userRepository)
}
