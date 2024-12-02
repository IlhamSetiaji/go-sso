package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IGetAllApplicationsUseCaseResponse struct {
	Applications *[]entity.Application `json:"applications"`
}

type IGetAllApplicationsUseCase interface {
	Execute() (*IGetAllApplicationsUseCaseResponse, error)
}

type GetAllApplicationsUseCase struct {
	Log                   *logrus.Logger
	ApplicationRepository repository.IApplicationRepository
}

func NewGetAllApplicationsUseCase(log *logrus.Logger, applicationRepository repository.IApplicationRepository) IGetAllApplicationsUseCase {
	return &GetAllApplicationsUseCase{
		Log:                   log,
		ApplicationRepository: applicationRepository,
	}
}

func (uc *GetAllApplicationsUseCase) Execute() (*IGetAllApplicationsUseCaseResponse, error) {
	applications, err := uc.ApplicationRepository.GetAllApplications()
	if err != nil {
		return nil, err
	}

	return &IGetAllApplicationsUseCaseResponse{
		Applications: applications,
	}, nil
}

func GetAllApplicationsUseCaseFactory(log *logrus.Logger) IGetAllApplicationsUseCase {
	applicationRepository := repository.ApplicationRepositoryFactory(log)
	return &GetAllApplicationsUseCase{
		Log:                   log,
		ApplicationRepository: applicationRepository,
	}
}
