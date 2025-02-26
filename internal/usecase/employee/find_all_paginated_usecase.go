package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IFindAllPaginatedUseCaseRequest struct {
	Page         int    `json:"page"`
	PageSize     int    `json:"page_size"`
	Search       string `json:"search"`
	IsOnboarding string `json:"is_onboarding"`
}

type IFindAllPaginatedUseCaseResponse struct {
	Employees *[]entity.Employee `json:"employees"`
	Total     int64              `json:"total"`
}

type IFindAllPaginatedUseCase interface {
	Execute(request *IFindAllPaginatedUseCaseRequest) (*IFindAllPaginatedUseCaseResponse, error)
}

type FindAllPaginatedUseCase struct {
	Log                *logrus.Logger
	EmployeeRepository repository.IEmployeeRepository
}

func NewFindAllPaginatedUseCase(
	log *logrus.Logger,
	employeeRepository repository.IEmployeeRepository,
) IFindAllPaginatedUseCase {
	return &FindAllPaginatedUseCase{
		Log:                log,
		EmployeeRepository: employeeRepository,
	}
}

func (uc *FindAllPaginatedUseCase) Execute(req *IFindAllPaginatedUseCaseRequest) (*IFindAllPaginatedUseCaseResponse, error) {
	var isOnboarding string
	if req.IsOnboarding != "" {
		if req.IsOnboarding == "YES" || req.IsOnboarding == "NO" {
			isOnboarding = req.IsOnboarding
		} else {
			isOnboarding = ""
		}
	}
	employees, total, err := uc.EmployeeRepository.FindAllPaginated(req.Page, req.PageSize, req.Search, isOnboarding)
	if err != nil {
		return nil, err
	}

	return &IFindAllPaginatedUseCaseResponse{
		Employees: employees,
		Total:     total,
	}, nil
}

func FindAllPaginatedUseCaseFactory(log *logrus.Logger) IFindAllPaginatedUseCase {
	employeeRepository := repository.EmployeeRepositoryFactory(log)
	return NewFindAllPaginatedUseCase(log, employeeRepository)
}
