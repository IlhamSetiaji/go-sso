package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
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
}

func NewFindByOrganizationIDUseCase(
	log *logrus.Logger,
	jobLevelRepository repository.IJobLevelRepository,
	orgStructureRepository repository.IOrganizationStructureRepository,
) IFindByOrganizationIDUseCase {
	return &FindByOrganizationIDUseCase{
		Log:                    log,
		JobLevelRepository:     jobLevelRepository,
		OrgStructureRepository: orgStructureRepository,
	}
}

func (uc *FindByOrganizationIDUseCase) Execute(req *IFindByOrganizationIDUseCaseRequest) (*IFindByOrganizationIDUseCaseResponse, error) {
	orgStructures, err := uc.OrgStructureRepository.FindByOrganizationId(uuid.MustParse(req.OrganizationID))
	if err != nil {
		return nil, err
	}

	var jobLevels []entity.JobLevel
	var jobLevelIDs []uuid.UUID

	for _, orgStructure := range *orgStructures {
		jobLevelIDs = append(jobLevelIDs, orgStructure.JobLevelID)
	}

	jobLevelPtrs, err := uc.JobLevelRepository.FindByIds(jobLevelIDs)
	if err != nil {
		return nil, err
	}

	jobLevels = append(jobLevels, *jobLevelPtrs...)

	return &IFindByOrganizationIDUseCaseResponse{
		JobLevels: dto.ConvertToJobLevelResponse(&jobLevels),
	}, nil
}

func FindByOrganizationIDUseCaseFactory(log *logrus.Logger) IFindByOrganizationIDUseCase {
	jobLevelRepository := repository.JobLevelRepositoryFactory(log)
	orgStructureRepository := repository.OrganizationStructureRepositoryFactory(log)
	return NewFindByOrganizationIDUseCase(log, jobLevelRepository, orgStructureRepository)
}
