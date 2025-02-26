package usecase

import (
	"app/go-sso/internal/entity"
	messaging "app/go-sso/internal/messaging/employee"
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

	countKanbanMessageFactory := messaging.CountKanbanProgressByEmployeeIDMessageFactory(uc.Log)

	for i, employee := range *employees {
		var employeeKanbanProgress entity.EmployeeKanbanProgress

		resp, err := countKanbanMessageFactory.Execute(&messaging.ICountKanbanProgressByEmployeeIDMessageRequest{
			EmployeeID: employee.ID.String(),
		})
		if err != nil {
			return nil, err
		}

		employeeKanbanProgress.TotalTask = resp.TotalTask
		employeeKanbanProgress.ToDo = resp.ToDo
		employeeKanbanProgress.InProgress = resp.InProgress
		employeeKanbanProgress.NeedReview = resp.NeedReview
		employeeKanbanProgress.Completed = resp.Completed

		employee.EmployeeKanbanProgress = &employeeKanbanProgress
		(*employees)[i] = employee
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
