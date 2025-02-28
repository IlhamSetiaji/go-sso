package response

import "github.com/google/uuid"

type OrganizationResponse struct {
	ID                             uuid.UUID `json:"id"`
	OrganizationTypeID             uuid.UUID `json:"organization_type_id"`
	Name                           string    `json:"name"`
	Region                         string    `json:"region"`
	Address                        string    `json:"address"`
	OrganizationType               OrganizationTypeResponse
	OrganizationStructureResponses []OrganizationStructureMinimalResponse `json:"organization_structures"`
}

type OrganizationMinimalResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type OrganizationForMessageResponse struct {
	ID                 uuid.UUID `json:"id"`
	OrganizationTypeID uuid.UUID `json:"organization_type_id"`
	Name               string    `json:"name"`
}
