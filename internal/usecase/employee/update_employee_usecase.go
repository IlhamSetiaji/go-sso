package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IUpdateEmployeeUsecaseRequest struct {
	ID             string `form:"id" validate:"required"`
	OrganizationID string `form:"organization_id" validate:"required"`
	Name           string `form:"name" validate:"required"`
	Email          string `form:"email" validate:"required"`
	MobilePhone    string `form:"mobile_phone" validate:"omitempty"`
	EndDate        string `form:"end_date" validate:"required,datetime=2006-01-02"`
	RetirementDate string `form:"retirement_date" validate:"omitempty,datetime=2006-01-02"`
}

type IUpdateEmployeeUsecaseResponse struct {
	EmployeeID string `json:"employee_id"`
}

type IUpdateEmployeeUsecase interface {
	Execute(request *IUpdateEmployeeUsecaseRequest) (*IUpdateEmployeeUsecaseResponse, error)
}

type UpdateEmployeeUsecase struct {
	Log          *logrus.Logger
	UserRepo     repository.IUserRepository
	EmployeeRepo repository.IEmployeeRepository
}

func UpdateEmployeeUsecaseFactory(log *logrus.Logger) IUpdateEmployeeUsecase {
	userRepo := repository.UserRepositoryFactory(log)
	employeeRepo := repository.EmployeeRepositoryFactory(log)
	return &UpdateEmployeeUsecase{
		Log:          log,
		UserRepo:     userRepo,
		EmployeeRepo: employeeRepo,
	}
}

func (u *UpdateEmployeeUsecase) Execute(request *IUpdateEmployeeUsecaseRequest) (*IUpdateEmployeeUsecaseResponse, error) {
	employee, err := u.EmployeeRepo.FindById(uuid.MustParse(request.ID))
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	if employee == nil {
		return nil, errors.New("Employee not found")
	}

	employee.Name = request.Name
	employee.Email = request.Email
	employee.MobilePhone = request.MobilePhone
	employee.EndDate, err = time.Parse("2006-01-02", request.EndDate)
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	if request.RetirementDate != "" {
		employee.RetirementDate, err = time.Parse("2006-01-02", request.RetirementDate)
		if err != nil {
			u.Log.Error(err)
			return nil, err
		}
	}

	employee, err = u.EmployeeRepo.Update(&entity.Employee{
		ID:             employee.ID,
		OrganizationID: uuid.MustParse(request.OrganizationID),
		Name:           employee.Name,
		Email:          employee.Email,
		MobilePhone:    employee.MobilePhone,
		EndDate:        employee.EndDate,
		RetirementDate: employee.RetirementDate,
	})
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	return &IUpdateEmployeeUsecaseResponse{
		EmployeeID: employee.ID.String(),
	}, nil
}
