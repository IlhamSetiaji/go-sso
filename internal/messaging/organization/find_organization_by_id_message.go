package messaging

import (
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IFindOrganizationByIDMessageRequest struct {
	OrganizationID uuid.UUID `json:"organization_id"`
}

type IFindOrganizationByIDMessageResponse struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
}

type IFindOrganizationByIDMessage interface {
	Execute(request IFindOrganizationByIDMessageRequest) (*IFindOrganizationByIDMessageResponse, error)
}

type FindOrganizationByIDMessage struct {
	Log        *logrus.Logger
	Repository repository.IOrganizationRepository
}

func NewFindOrganizationByIDMessage(log *logrus.Logger, repository repository.IOrganizationRepository) IFindOrganizationByIDMessage {
	return &FindOrganizationByIDMessage{
		Log:        log,
		Repository: repository,
	}
}

func (m *FindOrganizationByIDMessage) Execute(request IFindOrganizationByIDMessageRequest) (*IFindOrganizationByIDMessageResponse, error) {
	organization, err := m.Repository.FindById(request.OrganizationID)
	if err != nil {
		return nil, err
	}

	return &IFindOrganizationByIDMessageResponse{
		OrganizationID: organization.ID,
		Name:           organization.Name,
	}, nil
}

func FindOrganizationByIDMessageFactory(log *logrus.Logger) IFindOrganizationByIDMessage {
	repository := repository.OrganizationRepositoryFactory(log)
	return NewFindOrganizationByIDMessage(log, repository)
}