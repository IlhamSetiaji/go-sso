package response

import "github.com/google/uuid"

type OrganizationResponse struct {
	ID                             uuid.UUID `json:"id"`
	OrganizationTypeID             uuid.UUID `json:"organization_type_id"`
	Name                           string    `json:"name"`
	OrganizationType               OrganizationTypeResponse
	OrganizationStructureResponses []OrganizationStructureMinimalResponse `json:"organization_structures"`
}
