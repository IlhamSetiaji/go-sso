package response

import "github.com/google/uuid"

type OrganizationResponse struct {
	ID                 uuid.UUID `json:"id"`
	OrganizationTypeID uuid.UUID `json:"organization_type_id"`
	Name               string    `json:"name"`
	OrganizationType   OrganizationTypeResponse
}

type OrganizationTypeResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type JobLevelResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Level string    `json:"level"`
}

type OrganizationStructureResponse struct {
	ID             uuid.UUID                       `json:"id"`
	OrganizationID uuid.UUID                       `json:"organization_id"`
	Name           string                          `json:"name"`
	JobLevelID     uuid.UUID                       `json:"job_level_id"`
	ParentID       *uuid.UUID                      `json:"parent_id,omitempty"`
	Level          int                             `json:"level"`
	Path           string                          `json:"path"`
	Organization   OrganizationResponse            `json:"organization"`
	JobLevel       JobLevelResponse                `json:"job_level"`
	Parent         *OrganizationStructureResponse  `json:"parent,omitempty"`
	Children       []OrganizationStructureResponse `json:"children,omitempty"`
}
