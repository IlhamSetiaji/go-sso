package usecase

import (
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IGetJobsByOrganizationIDUseCaseRequest struct {
	OrganizationID uuid.UUID `json:"organization_id"`
}

type IGetJobsByOrganizationIDUseCaseResponse struct {
	Jobs *[]response.JobResponse `json:"jobs"`
}

type IGetJobsByOrganizationIDUseCase interface {
	Execute(request *IGetJobsByOrganizationIDUseCaseRequest) (*IGetJobsByOrganizationIDUseCaseResponse, error)
}

type GetJobsByOrganizationIDUseCase struct {
	Log                    *logrus.Logger
	OrgStructureRepository repository.IOrganizationStructureRepository
	JobRepository          repository.IJobRepository
}

func NewGetJobsByOrganizationIDUseCase(log *logrus.Logger, orgStructureRepository repository.IOrganizationStructureRepository, jobRepository repository.IJobRepository) IGetJobsByOrganizationIDUseCase {
	return &GetJobsByOrganizationIDUseCase{
		Log:                    log,
		OrgStructureRepository: orgStructureRepository,
		JobRepository:          jobRepository,
	}
}

func (uc *GetJobsByOrganizationIDUseCase) Execute(request *IGetJobsByOrganizationIDUseCaseRequest) (*IGetJobsByOrganizationIDUseCaseResponse, error) {
	organizationID := request.OrganizationID

	organizationStructures, err := uc.OrgStructureRepository.GetOrganizationSructuresByOrganizationID(organizationID)
	if err != nil {
		return nil, err
	}

	var organizationStructureIDs []uuid.UUID
	for _, organizationStructure := range *organizationStructures {
		organizationStructureIDs = append(organizationStructureIDs, organizationStructure.ID)
	}

	jobs, err := uc.JobRepository.GetJobsByOrganizationStructureIDs(organizationStructureIDs)
	if err != nil {
		return nil, err
	}

	return &IGetJobsByOrganizationIDUseCaseResponse{
		Jobs: dto.ConvertToJobResponse(jobs),
	}, nil
}

func GetJobsByOrganizationIDUseCaseFactory(log *logrus.Logger) IGetJobsByOrganizationIDUseCase {
	orgStructureRepository := repository.OrganizationStructureRepositoryFactory(log)
	jobRepository := repository.JobRepositoryFactory(log)
	return NewGetJobsByOrganizationIDUseCase(log, orgStructureRepository, jobRepository)
}
