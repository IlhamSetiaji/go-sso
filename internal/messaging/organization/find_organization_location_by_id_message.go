package messaging

import (
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IFindOrganizationLocationByIDMessageRequest struct {
	OrganizationLocationID uuid.UUID `json:"organization_location_id"`
}

type IFindOrganizationLocationByIDMessageResponse struct {
	OrganizationLocationID string `json:"organization_location_id"`
	Name                   string `json:"name"`
}

type IFindOrganizationLocationByIDMessage interface {
	Execute(request IFindOrganizationLocationByIDMessageRequest) (*IFindOrganizationLocationByIDMessageResponse, error)
}

type FindOrganizationLocationByIDMessage struct {
	Log        *logrus.Logger
	Repository repository.IOrganizationLocationRepository
}

func NewFindOrganizationLocationByIDMessage(log *logrus.Logger, repository repository.IOrganizationLocationRepository) IFindOrganizationLocationByIDMessage {
	return &FindOrganizationLocationByIDMessage{
		Log:        log,
		Repository: repository,
	}
}

func (m *FindOrganizationLocationByIDMessage) Execute(request IFindOrganizationLocationByIDMessageRequest) (*IFindOrganizationLocationByIDMessageResponse, error) {
	organizationLocation, err := m.Repository.FindById(request.OrganizationLocationID)
	if err != nil {
		return nil, err
	}

	return &IFindOrganizationLocationByIDMessageResponse{
		OrganizationLocationID: organizationLocation.ID.String(),
		Name:                   organizationLocation.Name,
	}, nil
}

func FindOrganizationLocationByIDMessageFactory(log *logrus.Logger) IFindOrganizationLocationByIDMessage {
	repository := repository.OrganizationLocationRepositoryFactory(log)
	return NewFindOrganizationLocationByIDMessage(log, repository)
}
