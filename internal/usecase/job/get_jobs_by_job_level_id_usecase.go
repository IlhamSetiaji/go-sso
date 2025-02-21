package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	jobMsg "app/go-sso/internal/messaging/job"
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IGetJobsByJobLevelIDUseCaseRequest struct {
	JobLevelID     string `json:"job_level_id"`
	OrganizationID string `json:"organization_id"`
}

type IGetJobsByJobLevelIDUseCaseResponse struct {
	Jobs *[]response.JobResponse `json:"jobs"`
}

type IGetJobsByJobLevelIDUseCase interface {
	Execute(request *IGetJobsByJobLevelIDUseCaseRequest) (*IGetJobsByJobLevelIDUseCaseResponse, error)
}

type GetJobsByJobLevelIDUseCase struct {
	Log                    *logrus.Logger
	OrgStructureRepository repository.IOrganizationStructureRepository
	JobRepository          repository.IJobRepository
}

func NewGetJobsByJobLevelIDUseCase(log *logrus.Logger, orgStructureRepository repository.IOrganizationStructureRepository, jobRepository repository.IJobRepository) IGetJobsByJobLevelIDUseCase {
	return &GetJobsByJobLevelIDUseCase{
		Log:                    log,
		OrgStructureRepository: orgStructureRepository,
		JobRepository:          jobRepository,
	}
}

func (uc *GetJobsByJobLevelIDUseCase) Execute(request *IGetJobsByJobLevelIDUseCaseRequest) (*IGetJobsByJobLevelIDUseCaseResponse, error) {
	jobLevelID := request.JobLevelID
	uuidJobLevelID, err := uuid.Parse(jobLevelID)
	if err != nil {
		return nil, err
	}

	var organizationStructures *[]entity.OrganizationStructure

	if request.OrganizationID == "" {
		organizationStructures, err = uc.OrgStructureRepository.GetOrganizationSructuresByJobLevelID(uuidJobLevelID)
		if err != nil {
			return nil, err
		}
	} else {
		organizationID, err := uuid.Parse(request.OrganizationID)
		if err != nil {
			return nil, err
		}

		organizationStructures, err = uc.OrgStructureRepository.GetOrganizationSructuresByJobLevelIDAndOrganizationID(uuidJobLevelID, organizationID)
		if err != nil {
			return nil, err
		}
	}

	var organizationStructureIDs []uuid.UUID
	for _, organizationStructure := range *organizationStructures {
		organizationStructureIDs = append(organizationStructureIDs, organizationStructure.ID)
	}

	jobs, err := uc.JobRepository.GetJobsByOrganizationStructureIDs(organizationStructureIDs)
	if err != nil {
		return nil, err
	}

	for i, job := range *jobs {
		children, err := uc.JobRepository.FindAllChildren(job.ID)
		if err != nil {
			return nil, err
		}
		(*jobs)[i].Children = children

		uc.Log.Info("job.ID: ", job.ID.String())

		msgFactory := jobMsg.FindJobPlafonByJobIDMessageFactory(uc.Log)
		msgReq := &jobMsg.IFindJobPlafonByJobIDMessageRequest{
			JobID: job.ID.String(),
		}
		msgResp, err := msgFactory.Execute(msgReq)
		if err != nil {
			return nil, err
		}

		(*jobs)[i].Plafon = msgResp.Plafon

	}

	return &IGetJobsByJobLevelIDUseCaseResponse{
		Jobs: dto.ConvertToJobResponse(jobs),
	}, nil
}

func GetJobsByJobLevelIDUseCaseFactory(log *logrus.Logger) IGetJobsByJobLevelIDUseCase {
	orgStructureRepository := repository.OrganizationStructureRepositoryFactory(log)
	jobRepository := repository.JobRepositoryFactory(log)
	return NewGetJobsByJobLevelIDUseCase(log, orgStructureRepository, jobRepository)
}
