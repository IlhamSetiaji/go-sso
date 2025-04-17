package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IFindEmployeeRecruitmentManagerUsecaseResponse struct {
	Employee *entity.Employee `json:"employee"`
}

type IFindEmployeeRecruitmentManagerUsecase interface {
	Execute() (*IFindEmployeeRecruitmentManagerUsecaseResponse, error)
}

type FindEmployeeRecruitmentManagerUsecase struct {
	Log                *logrus.Logger
	EmployeeRepository repository.IEmployeeRepository
}

func NewFindEmployeeRecruitmentManagerUsecase(
	log *logrus.Logger,
	employeeRepository repository.IEmployeeRepository,
) IFindEmployeeRecruitmentManagerUsecase {
	return &FindEmployeeRecruitmentManagerUsecase{
		Log:                log,
		EmployeeRepository: employeeRepository,
	}
}

func (uc *FindEmployeeRecruitmentManagerUsecase) Execute() (*IFindEmployeeRecruitmentManagerUsecaseResponse, error) {
	employee, err := uc.EmployeeRepository.FindEmployeeRecruitmentManager()
	if err != nil {
		return nil, err
	}

	return &IFindEmployeeRecruitmentManagerUsecaseResponse{
		Employee: &(*employee)[0],
	}, nil
}

func FindEmployeeRecruitmentManagerUsecaseFactory(log *logrus.Logger) IFindEmployeeRecruitmentManagerUsecase {
	employeeRepository := repository.EmployeeRepositoryFactory(log)
	return NewFindEmployeeRecruitmentManagerUsecase(log, employeeRepository)
}
