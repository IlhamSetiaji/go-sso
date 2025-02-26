package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IStoreEmployeeUsecaseRequest struct {
	OrganizationID string `form:"organization_id" validate:"required"` // Add this line
	UserID         string `form:"user_id" validate:"omitempty"`
	Name           string `form:"name" validate:"required"`
	Email          string `form:"email" validate:"required"`
	MobilePhone    string `form:"mobile_phone" validate:"omitempty"`
	EndDate        string `form:"end_date" validate:"required,datetime=2006-01-02"`
	RetirementDate string `form:"retirement_date" validate:"omitempty,datetime=2006-01-02"`
}

type IStoreEmployeeUsecaseResponse struct {
	EmployeeID string `json:"employee_id"`
}

type IStoreEmployeeUsecase interface {
	Execute(request *IStoreEmployeeUsecaseRequest) (*IStoreEmployeeUsecaseResponse, error)
}

type StoreEmployeeUsecase struct {
	Log          *logrus.Logger
	UserRepo     repository.IUserRepository
	EmployeeRepo repository.IEmployeeRepository
}

func StoreEmployeeUsecaseFactory(log *logrus.Logger) IStoreEmployeeUsecase {
	userRepo := repository.UserRepositoryFactory(log)
	employeeRepo := repository.EmployeeRepositoryFactory(log)
	return &StoreEmployeeUsecase{
		Log:          log,
		UserRepo:     userRepo,
		EmployeeRepo: employeeRepo,
	}
}

func (u *StoreEmployeeUsecase) Execute(request *IStoreEmployeeUsecaseRequest) (*IStoreEmployeeUsecaseResponse, error) {
	var endDate time.Time
	var err error
	if request.EndDate != "" {
		endDate, err = time.Parse("2006-01-02", request.EndDate)
		if err != nil {
			u.Log.Error(err)
			return nil, err
		}
	}

	var retirementDate time.Time
	if request.RetirementDate != "" {
		retirementDate, err = time.Parse("2006-01-02", request.RetirementDate)
		if err != nil {
			u.Log.Error(err)
			return nil, err
		}
	}

	employee := &entity.Employee{
		Name:           request.Name,
		Email:          request.Email,
		MobilePhone:    request.MobilePhone,
		EndDate:        endDate,
		RetirementDate: retirementDate,
		OrganizationID: uuid.MustParse(request.OrganizationID), // Add this line
		IsOnboarding:   "NO",
	}

	employee, err = u.EmployeeRepo.Store(employee)
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	if request.UserID != "" {
		user, err := u.UserRepo.FindById(uuid.MustParse(request.UserID))
		if err != nil {
			u.Log.Error(err)
			return nil, err
		}

		if user == nil {
			return nil, errors.New("User not found")
		}

		_, err = u.UserRepo.UpdateUserOnly(&entity.User{
			ID:         user.ID,
			EmployeeID: &employee.ID,
		})

		if err != nil {
			u.Log.Error(err)
			return nil, err
		}
	}

	return &IStoreEmployeeUsecaseResponse{
		EmployeeID: employee.ID.String(),
	}, nil
}
