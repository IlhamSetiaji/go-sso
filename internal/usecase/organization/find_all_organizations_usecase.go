package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IFindAllOrganizationsUsecaseResponse struct {
	Organizations []entity.Organization `json:"organizations"`
}

type IFindAllOrganizationsUsecase interface {
	Execute() (*IFindAllOrganizationsUsecaseResponse, error)
}

type FindAllOrganizationsUsecase struct {
	Log              *logrus.Logger
	OrganizationRepo repository.IOrganizationRepository
}

func FindAllOrganizationsUsecaseFactory(log *logrus.Logger) IFindAllOrganizationsUsecase {
	organizationRepo := repository.OrganizationRepositoryFactory(log)
	return &FindAllOrganizationsUsecase{
		Log:              log,
		OrganizationRepo: organizationRepo,
	}
}

func (u *FindAllOrganizationsUsecase) Execute() (*IFindAllOrganizationsUsecaseResponse, error) {
	organizations, err := u.OrganizationRepo.FindAllOrganizations()
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	return &IFindAllOrganizationsUsecaseResponse{
		Organizations: *organizations,
	}, nil
}
