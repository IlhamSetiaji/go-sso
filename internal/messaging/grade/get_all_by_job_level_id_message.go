package messaging

import (
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	"app/go-sso/internal/repository"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IGetAllByJobLevelIDMessageRequest struct {
	JobLevelID string `json:"job_level_id"`
}

type IGetAllByJobLevelIDMessageResponse struct {
	Grade []response.GradeResponse `json:"grade"`
}

type IGetAllByJobLevelIDMessage interface {
	Execute(req *IGetAllByJobLevelIDMessageRequest) (*IGetAllByJobLevelIDMessageResponse, error)
}

type GetAllByJobLevelIDMessage struct {
	Log                *logrus.Logger
	JobLevelRepository repository.IJobLevelRepository
	GradeRepository    repository.IGradeRepository
}

func NewGetAllByJobLevelIDMessage(
	log *logrus.Logger,
	jobLevelRepository repository.IJobLevelRepository,
	gradeRepository repository.IGradeRepository,
) IGetAllByJobLevelIDMessage {
	return &GetAllByJobLevelIDMessage{
		Log:                log,
		JobLevelRepository: jobLevelRepository,
		GradeRepository:    gradeRepository,
	}
}

func (uc *GetAllByJobLevelIDMessage) Execute(req *IGetAllByJobLevelIDMessageRequest) (*IGetAllByJobLevelIDMessageResponse, error) {
	jobLevelID, err := uuid.Parse(req.JobLevelID)
	if err != nil {
		return nil, err
	}

	jobLevel, err := uc.JobLevelRepository.FindById(jobLevelID)
	if err != nil {
		return nil, err
	}

	if jobLevel == nil {
		return nil, errors.New("job level not found")
	}

	grades, err := uc.GradeRepository.FindAllByJobLevelID(jobLevelID)
	if err != nil {
		return nil, err
	}

	var gradeResps []response.GradeResponse
	for _, grade := range *grades {
		gradeResps = append(gradeResps, *dto.ConvertToSingleGradeResponse(&grade))

	}

	return &IGetAllByJobLevelIDMessageResponse{
		Grade: gradeResps,
	}, nil
}
