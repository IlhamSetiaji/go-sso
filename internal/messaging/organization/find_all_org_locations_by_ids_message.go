package messaging

import (
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IFindAllOrgLocationsByIDsMessageRequest struct {
	IncludedIDs []string `json:"included_ids"`
}

type IFindAllOrgLocationsByIDsMessageResponse struct {
	OrgLocations *[]response.OrganizationLocationResponse `json:"org_locations"`
}

type IFindAllOrgLocationsByIDsMessage interface {
	Execute(request *IFindAllOrgLocationsByIDsMessageRequest) (*IFindAllOrgLocationsByIDsMessageResponse, error)
}

type FindAllOrgLocationsByIDsMessage struct {
	Log        *logrus.Logger
	OrgLocRepo repository.IOrganizationLocationRepository
}

func NewFindAllOrgLocationsByIDsMessage(
	log *logrus.Logger,
	orgLocRepo repository.IOrganizationLocationRepository,
) IFindAllOrgLocationsByIDsMessage {
	return &FindAllOrgLocationsByIDsMessage{
		Log:        log,
		OrgLocRepo: orgLocRepo,
	}
}

func (uc *FindAllOrgLocationsByIDsMessage) Execute(req *IFindAllOrgLocationsByIDsMessageRequest) (*IFindAllOrgLocationsByIDsMessageResponse, error) {
	orgLocs, err := uc.OrgLocRepo.FindAllOrganizationLocations(req.IncludedIDs)
	if err != nil {
		return nil, err
	}

	return &IFindAllOrgLocationsByIDsMessageResponse{
		OrgLocations: dto.ConvertToOrganizationLocationResponse(orgLocs),
	}, nil
}

func FindAllOrgLocationsByIDsMessageFactory(log *logrus.Logger) IFindAllOrgLocationsByIDsMessage {
	orgLocRepo := repository.OrganizationLocationRepositoryFactory(log)
	return NewFindAllOrgLocationsByIDsMessage(log, orgLocRepo)
}
