package messaging

import (
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IFindEmployeeByMidsuitIDMessageRequest struct {
	MidsuitID string `json:"midsuit_id"`
}

type IFindEmployeeByMidsuitIDMessageResponse struct {
	EmployeeID string                     `json:"employee_id"`
	Employee   *response.EmployeeResponse `json:"employee"`
}

type IFindEmployeeByMidsuitIDMessage interface {
	Execute(req IFindEmployeeByMidsuitIDMessageRequest) (*IFindEmployeeByMidsuitIDMessageResponse, error)
}

type FindEmployeeByMidsuitIDMessage struct {
	Log                *logrus.Logger
	EmployeeRepository repository.IEmployeeRepository
}

func NewFindEmployeeByMidsuitIDMessage(log *logrus.Logger, employeeRepository repository.IEmployeeRepository) IFindEmployeeByMidsuitIDMessage {
	return &FindEmployeeByMidsuitIDMessage{
		Log:                log,
		EmployeeRepository: employeeRepository,
	}
}

func (m *FindEmployeeByMidsuitIDMessage) Execute(req IFindEmployeeByMidsuitIDMessageRequest) (*IFindEmployeeByMidsuitIDMessageResponse, error) {
	employee, err := m.EmployeeRepository.FindByMidsuitID(req.MidsuitID)
	if err != nil {
		return nil, err
	}

	employeeResponse := dto.ConvertToSingleEmployeeResponse(employee)

	return &IFindEmployeeByMidsuitIDMessageResponse{
		EmployeeID: employee.ID.String(),
		Employee:   employeeResponse,
	}, nil
}

func FindEmployeeByMidsuitIDMessageFactory(log *logrus.Logger) IFindEmployeeByMidsuitIDMessage {
	repository := repository.EmployeeRepositoryFactory(log)
	return NewFindEmployeeByMidsuitIDMessage(log, repository)
}
