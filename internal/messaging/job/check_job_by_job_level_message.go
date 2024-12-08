package messaging

import (
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ICheckJobByJobLevelMessageRequest struct {
	JobID      uuid.UUID `json:"job_id"`
	JobLevelID uuid.UUID `json:"job_level_id"`
}

type ICheckJobByJobLevelMessageResponse struct {
	JobID      uuid.UUID `json:"job_id"`
	JobLevelID uuid.UUID `json:"job_level_id"`
	Exists     bool      `json:"exists"`
}

type ICheckJobByJobLevelMessage interface {
	Execute(request ICheckJobByJobLevelMessageRequest) (*ICheckJobByJobLevelMessageResponse, error)
}

type CheckJobByJobLevelMessage struct {
	Log                    *logrus.Logger
	Repository             repository.IJobRepository
	OrgStructureRepository repository.IOrganizationStructureRepository
	JobLevelRepostiory     repository.IJobLevelRepository
}

func NewCheckJobByJobLevelMessage(log *logrus.Logger, repository repository.IJobRepository, orgStructureRepository repository.IOrganizationStructureRepository, jobLevelRepository repository.IJobLevelRepository) ICheckJobByJobLevelMessage {
	return &CheckJobByJobLevelMessage{
		Log:                    log,
		Repository:             repository,
		OrgStructureRepository: orgStructureRepository,
		JobLevelRepostiory:     jobLevelRepository,
	}
}

func (m *CheckJobByJobLevelMessage) Execute(request ICheckJobByJobLevelMessageRequest) (*ICheckJobByJobLevelMessageResponse, error) {
	job, err := m.Repository.FindById(request.JobID)
	if err != nil {
		return nil, err
	}

	_, err = m.JobLevelRepostiory.FindById(request.JobLevelID)
	if err != nil {
		return nil, err
	}

	orgStructure, err := m.OrgStructureRepository.GetOrganizationSructuresByJobLevelID(request.JobLevelID)
	if err != nil {
		return nil, err
	}

	exists := false

	for _, os := range *orgStructure {
		for _, j := range os.Jobs {
			if j.ID == job.ID {
				exists = true
				break
			}
		}
	}

	return &ICheckJobByJobLevelMessageResponse{
		JobID:      request.JobID,
		JobLevelID: request.JobLevelID,
		Exists:     exists,
	}, nil
}

func CheckJobByJobLevelMessageFactory(log *logrus.Logger) ICheckJobByJobLevelMessage {
	jobRepository := repository.JobRepositoryFactory(log)
	orgStructureRepository := repository.OrganizationStructureRepositoryFactory(log)
	jobLevelRepository := repository.JobLevelRepositoryFactory(log)
	return NewCheckJobByJobLevelMessage(log, jobRepository, orgStructureRepository, jobLevelRepository)
}
