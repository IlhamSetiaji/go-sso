package messaging

import (
	"app/go-sso/internal/repository"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IFindAllOrgStructureChildrenIDsMessageRequest struct {
	ParentID string `json:"parent_id"`
}

type IFindAllOrgStructureChildrenIDsMessageResponse struct {
	ChildrenIDs []string `json:"children_ids"`
}

type IFindAllOrgStructureChildrenIDsMessage interface {
	Execute(request *IFindAllOrgStructureChildrenIDsMessageRequest) (*IFindAllOrgStructureChildrenIDsMessageResponse, error)
}

type FindAllOrgStructureChildrenIDsMessage struct {
	Log        *logrus.Logger
	Repository repository.IOrganizationStructureRepository
}

func NewFindAllOrgStructureChildrenIDsMessage(
	log *logrus.Logger,
	repository repository.IOrganizationStructureRepository,
) IFindAllOrgStructureChildrenIDsMessage {
	return &FindAllOrgStructureChildrenIDsMessage{
		Log:        log,
		Repository: repository,
	}
}

func (uc *FindAllOrgStructureChildrenIDsMessage) Execute(req *IFindAllOrgStructureChildrenIDsMessageRequest) (*IFindAllOrgStructureChildrenIDsMessageResponse, error) {
	parentUUID, err := uuid.Parse(req.ParentID)
	if err != nil {
		return nil, err
	}

	exist, err := uc.Repository.FindById(parentUUID)
	if err != nil {
		return nil, err
	}
	if exist == nil {
		return nil, errors.New("parent not found")
	}

	childrenIDs, err := uc.Repository.FindAllChildrenIDs(parentUUID)
	if err != nil {
		return nil, err
	}

	childrenIDStrings := make([]string, len(childrenIDs))
	for i, id := range childrenIDs {
		childrenIDStrings[i] = id.String()
	}

	return &IFindAllOrgStructureChildrenIDsMessageResponse{
		ChildrenIDs: childrenIDStrings,
	}, nil
}

func FindAllOrgStructureChildrenIDsMessageFactory(log *logrus.Logger) IFindAllOrgStructureChildrenIDsMessage {
	organizationStructureRepository := repository.OrganizationStructureRepositoryFactory(log)
	return NewFindAllOrgStructureChildrenIDsMessage(log, organizationStructureRepository)
}
