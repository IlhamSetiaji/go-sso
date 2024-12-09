package usecase

import (
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type ICountEmployeeRetiredEndByDateRangeUseCaseRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type ICountEmployeeRetiredEndByDateRangeUseCaseResponse struct {
	Total int64 `json:"total"`
}

type ICountEmployeeRetiredEndByDateRangeUseCase interface {
	Execute(request *ICountEmployeeRetiredEndByDateRangeUseCaseRequest) (*ICountEmployeeRetiredEndByDateRangeUseCaseResponse, error)
}

type CountEmployeeRetiredEndByDateRangeUseCase struct {
	Log                *logrus.Logger
	EmployeeRepository repository.IEmployeeRepository
}

func NewCountEmployeeRetiredEndByDateRangeUseCase(
	log *logrus.Logger,
	employeeRepository repository.IEmployeeRepository,
) ICountEmployeeRetiredEndByDateRangeUseCase {
	return &CountEmployeeRetiredEndByDateRangeUseCase{
		Log:                log,
		EmployeeRepository: employeeRepository,
	}
}

func (uc *CountEmployeeRetiredEndByDateRangeUseCase) Execute(req *ICountEmployeeRetiredEndByDateRangeUseCaseRequest) (*ICountEmployeeRetiredEndByDateRangeUseCaseResponse, error) {
	total, err := uc.EmployeeRepository.CountEmployeeRetiredEndByDateRange(req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	return &ICountEmployeeRetiredEndByDateRangeUseCaseResponse{
		Total: total,
	}, nil
}

func CountEmployeeRetiredEndByDateRangeUseCaseFactory(log *logrus.Logger) ICountEmployeeRetiredEndByDateRangeUseCase {
	employeeRepository := repository.EmployeeRepositoryFactory(log)
	return NewCountEmployeeRetiredEndByDateRangeUseCase(log, employeeRepository)
}
