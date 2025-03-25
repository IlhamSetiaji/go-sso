package usecase

import (
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	"app/go-sso/internal/repository"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IGetAllByJobLevelIDUseCaseRequest struct {
	JobLevelID string `json:"job_level_id"`
}

type IGetAllByJobLevelIDUseCaseResponse struct {
	Grade []response.GradeResponse `json:"grade"`
}

type IGetAllByJobLevelIDUseCase interface {
	Execute(req *IGetAllByJobLevelIDUseCaseRequest) (*IGetAllByJobLevelIDUseCaseResponse, error)
}

type GetAllByJobLevelIDUseCase struct {
	Log                *logrus.Logger
	JobLevelRepository repository.IJobLevelRepository
	GradeRepository    repository.IGradeRepository
}

func NewGetAllByJobLevelIDUseCase(
	log *logrus.Logger,
	jobLevelRepository repository.IJobLevelRepository,
	gradeRepository repository.IGradeRepository,
) IGetAllByJobLevelIDUseCase {
	return &GetAllByJobLevelIDUseCase{
		Log:                log,
		JobLevelRepository: jobLevelRepository,
		GradeRepository:    gradeRepository,
	}
}

func (uc *GetAllByJobLevelIDUseCase) Execute(req *IGetAllByJobLevelIDUseCaseRequest) (*IGetAllByJobLevelIDUseCaseResponse, error) {
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

	return &IGetAllByJobLevelIDUseCaseResponse{
		Grade: gradeResps,
	}, nil
}

func GetAllByJobLevelIDUseCaseFactory(
	log *logrus.Logger,
) IGetAllByJobLevelIDUseCase {
	jobLevelRepository := repository.JobLevelRepositoryFactory(log)
	gradeRepository := repository.GradeRepositoryFactory(log)
	return NewGetAllByJobLevelIDUseCase(log, jobLevelRepository, gradeRepository)
}
