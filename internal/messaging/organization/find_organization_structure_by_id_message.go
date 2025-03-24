package messaging

import (
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IFindOrganizationStructureByIDMessageRequest struct {
	OrganizationStructureID uuid.UUID `json:"organization_structure_id"`
}

type IFindOrganizationStructureByIDMessageResponse struct {
	OrganizationStructureID string `json:"organization_structure_id"`
	Name                    string `json:"name"`
	MidsuitID               string `json:"midsuit_id"`
}

type IFindOrganizationStructureByIDMessage interface {
	Execute(request IFindOrganizationStructureByIDMessageRequest) (*IFindOrganizationStructureByIDMessageResponse, error)
}

type FindOrganizationStructureByIDMessage struct {
	Log        *logrus.Logger
	Repository repository.IOrganizationStructureRepository
}

func NewFindOrganizationStructureByIDMessage(log *logrus.Logger, repository repository.IOrganizationStructureRepository) IFindOrganizationStructureByIDMessage {
	return &FindOrganizationStructureByIDMessage{
		Log:        log,
		Repository: repository,
	}
}

func (m *FindOrganizationStructureByIDMessage) Execute(request IFindOrganizationStructureByIDMessageRequest) (*IFindOrganizationStructureByIDMessageResponse, error) {
	organizationStructure, err := m.Repository.FindById(request.OrganizationStructureID)
	if err != nil {
		return nil, err
	}

	return &IFindOrganizationStructureByIDMessageResponse{
		OrganizationStructureID: organizationStructure.ID.String(),
		Name:                    organizationStructure.Name,
		MidsuitID:               organizationStructure.MidsuitID,
	}, nil
}

func FindOrganizationStructureByIDMessageFactory(log *logrus.Logger) IFindOrganizationStructureByIDMessage {
	repository := repository.OrganizationStructureRepositoryFactory(log)
	return NewFindOrganizationStructureByIDMessage(log, repository)
}
