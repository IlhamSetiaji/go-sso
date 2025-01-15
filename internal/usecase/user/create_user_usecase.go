package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ICreateUserUseCaseRequest struct {
	User    *entity.User `json:"user"`
	RoleIDs []string     `json:"role_ids[]"`
}

type ICreateUserUseCaseResponse struct {
	User *entity.User `json:"user"`
}

type ICreateUserUseCase interface {
	Execute(request ICreateUserUseCaseRequest) (ICreateUserUseCaseResponse, error)
}

type CreateUserUseCase struct {
	Log            *logrus.Logger
	UserRepository repository.IUserRepository
}

func NewCreateUserUseCase(log *logrus.Logger, userRepository repository.IUserRepository) ICreateUserUseCase {
	return &CreateUserUseCase{
		Log:            log,
		UserRepository: userRepository,
	}
}

func (uc *CreateUserUseCase) Execute(request ICreateUserUseCaseRequest) (ICreateUserUseCaseResponse, error) {
	uc.Log.Info("CreateUserUseCase.Execute")

	// Convert RoleIDs from []string to []uuid.UUID
	var roleUUIDs []uuid.UUID
	for _, roleID := range request.RoleIDs {
		roleUUID, err := uuid.Parse(roleID)
		if err != nil {
			return ICreateUserUseCaseResponse{}, err
		}
		roleUUIDs = append(roleUUIDs, roleUUID)
	}

	request.User.EmailVerifiedAt = time.Now()

	user, err := uc.UserRepository.CreateUser(request.User, roleUUIDs)
	if err != nil {
		return ICreateUserUseCaseResponse{}, err
	}

	return ICreateUserUseCaseResponse{
		User: user,
	}, nil
}

func CreateUserUseCaseFactory(log *logrus.Logger) ICreateUserUseCase {
	userRepository := repository.UserRepositoryFactory(log)
	return &CreateUserUseCase{
		Log:            log,
		UserRepository: userRepository,
	}
}
