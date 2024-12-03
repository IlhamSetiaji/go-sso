package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IFindApplicationByNameUsecaseRequest struct {
	Name string `json:"name"`
}

type IFindApplicationByNameUsecaseResponse struct {
	Application *entity.Application `json:"application"`
}

type IFindApplicationByNameUsecase interface {
	Execute(request *IFindApplicationByNameUsecaseRequest) (*IFindApplicationByNameUsecaseResponse, error)
}

type FindApplicationByNameUsecase struct {
	Log                   *logrus.Logger
	ApplicationRepository repository.IApplicationRepository
}

func NewFindApplicationByNameUsecase(log *logrus.Logger, applicationRepository repository.IApplicationRepository) *FindApplicationByNameUsecase {
	return &FindApplicationByNameUsecase{
		Log:                   log,
		ApplicationRepository: applicationRepository,
	}
}

func (u *FindApplicationByNameUsecase) Execute(request *IFindApplicationByNameUsecaseRequest) (*IFindApplicationByNameUsecaseResponse, error) {
	application, err := u.ApplicationRepository.FindApplicationByName(request.Name)
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}
	return &IFindApplicationByNameUsecaseResponse{
		Application: application,
	}, nil
}

func FindApplicationByNameUsecaseFactory(log *logrus.Logger) IFindApplicationByNameUsecase {
	applicationRepository := repository.ApplicationRepositoryFactory(log)
	return NewFindApplicationByNameUsecase(log, applicationRepository)
}
