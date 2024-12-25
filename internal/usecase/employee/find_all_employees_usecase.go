package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IFindAllEmployeesUsecaseResponse struct {
	Employees []entity.Employee `json:"employees"`
}

type IFindAllEmployeesUsecase interface {
	Execute() (*IFindAllEmployeesUsecaseResponse, error)
}

type FindAllEmployeesUsecase struct {
	Log          *logrus.Logger
	EmployeeRepo repository.IEmployeeRepository
}

func FindAllEmployeesUsecaseFactory(log *logrus.Logger) IFindAllEmployeesUsecase {
	employeeRepo := repository.EmployeeRepositoryFactory(log)
	return &FindAllEmployeesUsecase{
		Log:          log,
		EmployeeRepo: employeeRepo,
	}
}

func (u *FindAllEmployeesUsecase) Execute() (*IFindAllEmployeesUsecaseResponse, error) {
	employees, err := u.EmployeeRepo.FindAllEmployees()
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	return &IFindAllEmployeesUsecaseResponse{
		Employees: *employees,
	}, nil
}
