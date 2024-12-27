package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IFindAllOrganizationLocationsUsecaseResponse struct {
	OrganizationLocations []entity.OrganizationLocation `json:"organization_locations"`
}

type IFindAllOrganizationLocationsUsecase interface {
	Execute() (*IFindAllOrganizationLocationsUsecaseResponse, error)
}

type FindAllOrganizationLocationsUsecase struct {
	Log             *logrus.Logger
	OrgLocationRepo repository.IOrganizationLocationRepository
}

func FindAllOrganizationLocationsUsecaseFactory(log *logrus.Logger) IFindAllOrganizationLocationsUsecase {
	orgLocationRepo := repository.OrganizationLocationRepositoryFactory(log)
	return &FindAllOrganizationLocationsUsecase{
		Log:             log,
		OrgLocationRepo: orgLocationRepo,
	}
}

func (u *FindAllOrganizationLocationsUsecase) Execute() (*IFindAllOrganizationLocationsUsecaseResponse, error) {
	organizationLocations, err := u.OrgLocationRepo.FindAllOrganizationLocations([]string{})
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	return &IFindAllOrganizationLocationsUsecaseResponse{
		OrganizationLocations: *organizationLocations,
	}, nil
}
