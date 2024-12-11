package response

import "github.com/google/uuid"

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

type OrganizationStructureMinimalResponse struct {
	ID               uuid.UUID              `json:"id"`
	OrganizationID   uuid.UUID              `json:"organization_id"`
	Name             string                 `json:"name"`
	Level            int                    `json:"level"`
	Path             string                 `json:"path"`
	ParentID         *uuid.UUID             `json:"parent_id,omitempty"`
	JobLevelResponse map[string]interface{} `json:"job_level"`
}
