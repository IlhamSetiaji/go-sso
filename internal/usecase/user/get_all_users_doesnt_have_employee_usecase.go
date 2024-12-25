package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IGetAllUsersDoesntHaveEmployeeUsecaseResponse struct {
	Users []entity.User `json:"users"`
}

type IGetAllUsersDoesntHaveEmployeeUsecase interface {
	Execute() (*IGetAllUsersDoesntHaveEmployeeUsecaseResponse, error)
}

type GetAllUsersDoesntHaveEmployeeUsecase struct {
	Log          *logrus.Logger
	UserRepo     repository.IUserRepository
	EmployeeRepo repository.IEmployeeRepository
}

func GetAllUsersDoesntHaveEmployeeUsecaseFactory(log *logrus.Logger) IGetAllUsersDoesntHaveEmployeeUsecase {
	userRepo := repository.UserRepositoryFactory(log)
	employeeRepo := repository.EmployeeRepositoryFactory(log)
	return &GetAllUsersDoesntHaveEmployeeUsecase{
		Log:          log,
		UserRepo:     userRepo,
		EmployeeRepo: employeeRepo,
	}
}

func (u *GetAllUsersDoesntHaveEmployeeUsecase) Execute() (*IGetAllUsersDoesntHaveEmployeeUsecaseResponse, error) {
	users, err := u.UserRepo.GetAllUsersDoesNotHaveEmployee()
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	return &IGetAllUsersDoesntHaveEmployeeUsecaseResponse{
		Users: *users,
	}, nil
}
