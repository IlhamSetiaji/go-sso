package messaging

import (
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IFindAllJobsMessageRequest struct {
	IncludedIDs []string `json:"included_ids"`
}

type IFindAllJobsMessageResponse struct {
	Jobs *[]response.JobResponse `json:"jobs"`
}

type IFindAllJobsMessage interface {
	Execute(request *IFindAllJobsMessageRequest) (*IFindAllJobsMessageResponse, error)
}

type FindAllJobsMessage struct {
	Log     *logrus.Logger
	JobRepo repository.IJobRepository
}

func NewFindAllJobsMessage(
	log *logrus.Logger,
	jobRepo repository.IJobRepository,
) IFindAllJobsMessage {
	return &FindAllJobsMessage{
		Log:     log,
		JobRepo: jobRepo,
	}
}

func (uc *FindAllJobsMessage) Execute(req *IFindAllJobsMessageRequest) (*IFindAllJobsMessageResponse, error) {
	jobs, err := uc.JobRepo.FindAllJobs(req.IncludedIDs)
	if err != nil {
		return nil, err
	}

	return &IFindAllJobsMessageResponse{
		Jobs: dto.ConvertToJobResponse(jobs),
	}, nil
}

func FindAllJobsMessageFactory(log *logrus.Logger) IFindAllJobsMessage {
	jobRepo := repository.JobRepositoryFactory(log)
	return NewFindAllJobsMessage(log, jobRepo)
}
