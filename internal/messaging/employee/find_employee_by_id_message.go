package messaging

import (
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IFindEmployeeByIDMessageRequest struct {
	EmployeeID string `json:"employee_id"`
}

type IFindEmployeeByIDMessageResponse struct {
	EmployeeID string                     `json:"employee_id"`
	Employee   *response.EmployeeResponse `json:"employee"`
}

type IFindEmployeeByIDMessage interface {
	Execute(req IFindEmployeeByIDMessageRequest) (*IFindEmployeeByIDMessageResponse, error)
}

type FindEmployeeByIDMessage struct {
	Log                *logrus.Logger
	EmployeeRepository repository.IEmployeeRepository
}

func NewFindEmployeeByIDMessage(log *logrus.Logger, employeeRepository repository.IEmployeeRepository) IFindEmployeeByIDMessage {
	return &FindEmployeeByIDMessage{
		Log:                log,
		EmployeeRepository: employeeRepository,
	}
}

func (m *FindEmployeeByIDMessage) Execute(req IFindEmployeeByIDMessageRequest) (*IFindEmployeeByIDMessageResponse, error) {
	employee, err := m.EmployeeRepository.FindById(uuid.MustParse(req.EmployeeID))
	if err != nil {
		return nil, err
	}

	employeeResponse := dto.ConvertToSingleEmployeeResponse(employee)

	return &IFindEmployeeByIDMessageResponse{
		EmployeeID: employee.ID.String(),
		Employee:   employeeResponse,
	}, nil
}

func FindEmployeeByIDMessageFactory(log *logrus.Logger) IFindEmployeeByIDMessage {
	repository := repository.EmployeeRepositoryFactory(log)
	return NewFindEmployeeByIDMessage(log, repository)
}
