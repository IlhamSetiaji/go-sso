package messaging

import (
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IFindAllJobsByOrganizationIDMessageRequest struct {
	OrganizationID string `json:"organization_id"`
}

type IFindAllJobsByOrganizationIDMessageResponse struct {
	Jobs *[]response.JobResponse `json:"jobs"`
}

type IFindAllJobsByOrganizationIDMessage interface {
	Execute(request *IFindAllJobsByOrganizationIDMessageRequest) (*IFindAllJobsByOrganizationIDMessageResponse, error)
}

type FindAllJobsByOrganizationIDMessage struct {
	Log              *logrus.Logger
	JobRepo          repository.IJobRepository
	OrgStructureRepo repository.IOrganizationStructureRepository
}

func NewFindAllJobsByOrganizationIDMessage(
	log *logrus.Logger,
	jobRepo repository.IJobRepository,
	orgStructureRepo repository.IOrganizationStructureRepository,
) IFindAllJobsByOrganizationIDMessage {
	return &FindAllJobsByOrganizationIDMessage{
		Log:              log,
		JobRepo:          jobRepo,
		OrgStructureRepo: orgStructureRepo,
	}
}

func (uc *FindAllJobsByOrganizationIDMessage) Execute(req *IFindAllJobsByOrganizationIDMessageRequest) (*IFindAllJobsByOrganizationIDMessageResponse, error) {
	orgStructures, err := uc.OrgStructureRepo.FindAllOrgStructuresByOrganizationID(uuid.MustParse(req.OrganizationID))
	if err != nil {
		return nil, err
	}

	var orgStructureIDs []uuid.UUID
	for _, orgStructure := range *orgStructures {
		orgStructureIDs = append(orgStructureIDs, orgStructure.ID)
	}

	jobs, err := uc.JobRepo.GetJobsByOrganizationStructureIDs(orgStructureIDs)
	if err != nil {
		return nil, err
	}

	return &IFindAllJobsByOrganizationIDMessageResponse{
		Jobs: dto.ConvertToJobResponse(jobs),
	}, nil
}

func FindAllJobsByOrganizationIDMessageFactory(log *logrus.Logger) IFindAllJobsByOrganizationIDMessage {
	jobRepo := repository.JobRepositoryFactory(log)
	orgStructureRepo := repository.OrganizationStructureRepositoryFactory(log)
	return NewFindAllJobsByOrganizationIDMessage(log, jobRepo, orgStructureRepo)
}
