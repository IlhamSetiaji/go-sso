package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IFindByOrganizationIDUseCaseRequest struct {
	OrganizationID string `json:"organization_id"`
}

type IFindByOrganizationIDUseCaseResponse struct {
	JobLevels *[]response.JobLevelResponse `json:"job_levels"`
}

type IFindByOrganizationIDUseCase interface {
	Execute(request *IFindByOrganizationIDUseCaseRequest) (*IFindByOrganizationIDUseCaseResponse, error)
}

type FindByOrganizationIDUseCase struct {
	Log                    *logrus.Logger
	JobLevelRepository     repository.IJobLevelRepository
	OrgStructureRepository repository.IOrganizationStructureRepository
	JobRepository          repository.IJobRepository
}

func NewFindByOrganizationIDUseCase(
	log *logrus.Logger,
	jobLevelRepository repository.IJobLevelRepository,
	orgStructureRepository repository.IOrganizationStructureRepository,
	jobRepository repository.IJobRepository,
) IFindByOrganizationIDUseCase {
	return &FindByOrganizationIDUseCase{
		Log:                    log,
		JobLevelRepository:     jobLevelRepository,
		OrgStructureRepository: orgStructureRepository,
		JobRepository:          jobRepository,
	}
}

func (uc *FindByOrganizationIDUseCase) Execute(req *IFindByOrganizationIDUseCaseRequest) (*IFindByOrganizationIDUseCaseResponse, error) {
	jobs, err := uc.JobRepository.FindAllByKeys(map[string]interface{}{
		"organization_id": req.OrganizationID,
	})
	if err != nil {
		return nil, err
	}

	// distinct job -> job levels
	jobLevelMap := make(map[string]*entity.JobLevel)
	for _, job := range *jobs {
		jobLevel := job.JobLevel
		jobLevelMap[jobLevel.ID.String()] = &jobLevel
	}

	jobLevels := make([]response.JobLevelResponse, 0)
	for _, jobLevel := range jobLevelMap {
		jobLevelResponse := dto.ConvertToSingleJobLevelResponse(jobLevel)
		jobLevels = append(jobLevels, *jobLevelResponse)
	}

	return &IFindByOrganizationIDUseCaseResponse{
		JobLevels: &jobLevels,
	}, nil
}

func FindByOrganizationIDUseCaseFactory(log *logrus.Logger) IFindByOrganizationIDUseCase {
	jobLevelRepository := repository.JobLevelRepositoryFactory(log)
	orgStructureRepository := repository.OrganizationStructureRepositoryFactory(log)
	jobRepository := repository.JobRepositoryFactory(log)
	return NewFindByOrganizationIDUseCase(log, jobLevelRepository, orgStructureRepository, jobRepository)
}
