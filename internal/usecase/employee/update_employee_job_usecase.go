package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IUpdateEmployeeJobUsecaseRequest struct {
	ID                     string `form:"id" validate:"required"`
	EmployeeID             string `form:"employee_id" validate:"required"`
	Name                   string `form:"name" validate:"omitempty"`
	JobID                  string `form:"job_id" validate:"required"`
	EmpOrganizationID      string `form:"emp_organization_id" validate:"required"`
	OrganizationLocationID string `form:"organization_location_id" validate:"required"`
}

type IUpdateEmployeeJobUsecaseResponse struct {
	EmployeeJobID string `json:"employee_job_id"`
}

type IUpdateEmployeeJobUsecase interface {
	Execute(request *IUpdateEmployeeJobUsecaseRequest) (*IUpdateEmployeeJobUsecaseResponse, error)
}

type UpdateEmployeeJobUsecase struct {
	Log             *logrus.Logger
	EmployeeRepo    repository.IEmployeeRepository
	JobRepo         repository.IJobRepository
	OrgRepo         repository.IOrganizationRepository
	OrgLocationRepo repository.IOrganizationLocationRepository
}

func UpdateEmployeeJobUsecaseFactory(log *logrus.Logger) IUpdateEmployeeJobUsecase {
	employeeRepo := repository.EmployeeRepositoryFactory(log)
	jobRepo := repository.JobRepositoryFactory(log)
	orgRepo := repository.OrganizationRepositoryFactory(log)
	orgLocationRepo := repository.OrganizationLocationRepositoryFactory(log)
	return &UpdateEmployeeJobUsecase{
		Log:             log,
		EmployeeRepo:    employeeRepo,
		JobRepo:         jobRepo,
		OrgRepo:         orgRepo,
		OrgLocationRepo: orgLocationRepo,
	}
}

func (u *UpdateEmployeeJobUsecase) Execute(request *IUpdateEmployeeJobUsecaseRequest) (*IUpdateEmployeeJobUsecaseResponse, error) {
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
		ID:                     uuid.MustParse(request.ID),
		EmployeeID:             &employee.ID,
		Name:                   request.Name,
		JobID:                  uuid.MustParse(request.JobID),
		EmpOrganizationID:      uuid.MustParse(request.EmpOrganizationID),
		OrganizationLocationID: uuid.MustParse(request.OrganizationLocationID),
	}

	_, err = u.EmployeeRepo.UpdateEmployeeJob(employeeJob)
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	return &IUpdateEmployeeJobUsecaseResponse{
		EmployeeJobID: employeeJob.ID.String(),
	}, nil
}
