package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IStoreEmployeeJobUsecaseRequest struct {
	EmployeeID              string `form:"employee_id" validate:"required"`
	Name                    string `form:"name" validate:"omitempty"`
	JobID                   string `form:"job_id" validate:"required"`
	EmpOrganizationID       string `form:"emp_organization_id" validate:"required"`
	OrganizationLocationID  string `form:"organization_location_id" validate:"required"`
	OrganizationStructureID string `form:"organization_structure_id" validate:"omitempty"`
}

type IStoreEmployeeJobUsecaseResponse struct {
	EmployeeJobID string `json:"employee_job_id"`
}

type IStoreEmployeeJobUsecase interface {
	Execute(request *IStoreEmployeeJobUsecaseRequest) (*IStoreEmployeeJobUsecaseResponse, error)
}

type StoreEmployeeJobUsecase struct {
	Log             *logrus.Logger
	EmployeeRepo    repository.IEmployeeRepository
	JobRepo         repository.IJobRepository
	OrgRepo         repository.IOrganizationRepository
	OrgLocationRepo repository.IOrganizationLocationRepository
}

func StoreEmployeeJobUsecaseFactory(log *logrus.Logger) IStoreEmployeeJobUsecase {
	employeeRepo := repository.EmployeeRepositoryFactory(log)
	jobRepo := repository.JobRepositoryFactory(log)
	orgRepo := repository.OrganizationRepositoryFactory(log)
	orgLocationRepo := repository.OrganizationLocationRepositoryFactory(log)
	return &StoreEmployeeJobUsecase{
		Log:             log,
		EmployeeRepo:    employeeRepo,
		JobRepo:         jobRepo,
		OrgRepo:         orgRepo,
		OrgLocationRepo: orgLocationRepo,
	}
}

func (u *StoreEmployeeJobUsecase) Execute(request *IStoreEmployeeJobUsecaseRequest) (*IStoreEmployeeJobUsecaseResponse, error) {
	employee, err := u.EmployeeRepo.FindById(uuid.MustParse(request.EmployeeID))
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	if employee == nil {
		return nil, errors.New("Employee not found")
	}

	job, err := u.JobRepo.FindById(uuid.MustParse(request.JobID))
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	if job == nil {
		return nil, errors.New("Job not found")
	}

	org, err := u.OrgRepo.FindById(uuid.MustParse(request.EmpOrganizationID))
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	if org == nil {
		return nil, errors.New("Organization not found")
	}

	orgLocation, err := u.OrgLocationRepo.FindById(uuid.MustParse(request.OrganizationLocationID))
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	if orgLocation == nil {
		return nil, errors.New("Organization Location not found")
	}

	employeeJob := &entity.EmployeeJob{
		Name:                    request.Name,
		EmployeeID:              &employee.ID,
		JobID:                   uuid.MustParse(request.JobID),
		EmpOrganizationID:       uuid.MustParse(request.EmpOrganizationID),
		OrganizationLocationID:  uuid.MustParse(request.OrganizationLocationID),
		OrganizationStructureID: uuid.MustParse(request.OrganizationStructureID),
	}

	_, err = u.EmployeeRepo.StoreEmployeeJob(employeeJob)
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	return &IStoreEmployeeJobUsecaseResponse{
		EmployeeJobID: employeeJob.ID.String(),
	}, nil
}
