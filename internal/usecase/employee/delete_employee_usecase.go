package usecase

import (
	"app/go-sso/internal/repository"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IDeleteEmployeeUsecaseRequest struct {
	ID string `form:"id" validate:"required"`
}

type IDeleteEmployeeUsecaseResponse struct {
	EmployeeID string `json:"employee_id"`
}

type IDeleteEmployeeUsecase interface {
	Execute(request *IDeleteEmployeeUsecaseRequest) (*IDeleteEmployeeUsecaseResponse, error)
}

type DeleteEmployeeUsecase struct {
	Log          *logrus.Logger
	UserRepo     repository.IUserRepository
	EmployeeRepo repository.IEmployeeRepository
}

func DeleteEmployeeUsecaseFactory(log *logrus.Logger) IDeleteEmployeeUsecase {
	userRepo := repository.UserRepositoryFactory(log)
	employeeRepo := repository.EmployeeRepositoryFactory(log)
	return &DeleteEmployeeUsecase{
		Log:          log,
		UserRepo:     userRepo,
		EmployeeRepo: employeeRepo,
	}
}

func (u *DeleteEmployeeUsecase) Execute(request *IDeleteEmployeeUsecaseRequest) (*IDeleteEmployeeUsecaseResponse, error) {
	employee, err := u.EmployeeRepo.FindById(uuid.MustParse(request.ID))
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	if employee == nil {
		return nil, errors.New("Employee not found")
	}

	user, err := u.UserRepo.FindByIdOnly(uuid.MustParse(employee.User.ID.String()))
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	if user == nil {
		return nil, errors.New("User not found")
	}

	user.EmployeeID = nil

	if _, err := u.UserRepo.UpdateEmployeeIdToNull(user); err != nil {
		u.Log.Error(err)
		return nil, err
	}

	if err := u.EmployeeRepo.Delete(uuid.MustParse(request.ID)); err != nil {
		u.Log.Error(err)
		return nil, err
	}

	return &IDeleteEmployeeUsecaseResponse{EmployeeID: request.ID}, nil
}
