package usecase

import (
	"app/go-sso/internal/repository"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IUpdateEmployeeMidsuitUsecaseRequest struct {
	ID        string `json:"id" validate:"required"`
	MidsuitID string `json:"midsuit_id" validate:"required"`
}

type IUpdateEmployeeMidsuitUsecaseResponse struct {
	EmployeeID string `json:"employee_id"`
}

type IUpdateEmployeeMidsuitUsecase interface {
	Execute(request *IUpdateEmployeeMidsuitUsecaseRequest) (*IUpdateEmployeeMidsuitUsecaseResponse, error)
}

type UpdateEmployeeMidsuitUsecase struct {
	Log          *logrus.Logger
	EmployeeRepo repository.IEmployeeRepository
}

func UpdateEmployeeMidsuitUsecaseFactory(log *logrus.Logger) IUpdateEmployeeMidsuitUsecase {
	employeeRepo := repository.EmployeeRepositoryFactory(log)
	return &UpdateEmployeeMidsuitUsecase{
		Log:          log,
		EmployeeRepo: employeeRepo,
	}
}

func (u *UpdateEmployeeMidsuitUsecase) Execute(request *IUpdateEmployeeMidsuitUsecaseRequest) (*IUpdateEmployeeMidsuitUsecaseResponse, error) {
	parsedEmployeeID, err := uuid.Parse(request.ID)
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}
	employee, err := u.EmployeeRepo.FindById(parsedEmployeeID)
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	if employee == nil {
		return nil, errors.New("employee not found")
	}

	employee.MidsuitID = request.MidsuitID

	_, err = u.EmployeeRepo.UpdateEmployeeMidsuitID(employee.ID, employee.MidsuitID)
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	return &IUpdateEmployeeMidsuitUsecaseResponse{
		EmployeeID: employee.ID.String(),
	}, nil
}
