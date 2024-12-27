package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IFindAllEmployeesNotInUserUseCaseResponse struct {
	Employees []entity.Employee
}

type IFindAllEmployeesNotInUserUseCase interface {
	Execute() (*IFindAllEmployeesNotInUserUseCaseResponse, error)
}

type findAllEmployeesNotInUserUseCase struct {
	Log                *logrus.Logger
	EmployeeRepository repository.IEmployeeRepository
}

func NewFindAllEmployeesNotInUserUseCase(log *logrus.Logger, employeeRepository repository.IEmployeeRepository) IFindAllEmployeesNotInUserUseCase {
	return &findAllEmployeesNotInUserUseCase{
		Log:                log,
		EmployeeRepository: employeeRepository,
	}
}

func (u *findAllEmployeesNotInUserUseCase) Execute() (*IFindAllEmployeesNotInUserUseCaseResponse, error) {
	employees, err := u.EmployeeRepository.FindAllEmployeesNotInUsers()
	if err != nil {
		return nil, err
	}

	return &IFindAllEmployeesNotInUserUseCaseResponse{
		Employees: *employees,
	}, nil
}

func FindAllEmployeesNotInUserUseCaseFactory(log *logrus.Logger) IFindAllEmployeesNotInUserUseCase {
	employeeRepository := repository.EmployeeRepositoryFactory(log)
	return NewFindAllEmployeesNotInUserUseCase(log, employeeRepository)
}
