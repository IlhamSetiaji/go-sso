package messaging

import (
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IFindAllOrgMessageRequest struct {
	IncludedIDs []string `json:"included_ids"`
}

type IFindAllOrgMessageResponse struct {
	Organizations *[]response.OrganizationForMessageResponse `json:"organizations"`
}

type IFindAllOrgMessage interface {
	Execute(request *IFindAllOrgMessageRequest) (*IFindAllOrgMessageResponse, error)
}

type FindAllOrgMessage struct {
	Log     *logrus.Logger
	OrgRepo repository.IOrganizationRepository
}

func NewFindAllOrgMessage(
	log *logrus.Logger,
	orgRepo repository.IOrganizationRepository,
) IFindAllOrgMessage {
	return &FindAllOrgMessage{
		Log:     log,
		OrgRepo: orgRepo,
	}
}

func (uc *FindAllOrgMessage) Execute(req *IFindAllOrgMessageRequest) (*IFindAllOrgMessageResponse, error) {
	var uuids []uuid.UUID
	for _, id := range req.IncludedIDs {
		uuid, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		uuids = append(uuids, uuid)
	}

	organizations, err := uc.OrgRepo.FindByIDs(uuids)
	if err != nil {
		return nil, err
	}

	return &IFindAllOrgMessageResponse{
		Organizations: dto.ConvertToOrganizationForMessageResponse(organizations),
	}, nil
}

func FindAllOrgMessageFactory(log *logrus.Logger) IFindAllOrgMessage {
	orgRepo := repository.OrganizationRepositoryFactory(log)
	return NewFindAllOrgMessage(log, orgRepo)
}
