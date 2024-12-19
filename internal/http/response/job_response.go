package response

import (
	"github.com/google/uuid"
)

type JobResponse struct {
	ID                        uuid.UUID          `json:"id"`
	Name                      string             `json:"name"`
	OrganizationStructureID   uuid.UUID          `json:"organization_structure_id"`
	OrganizationStructureName string             `json:"organization_structure_name"`
	OrganizationID            uuid.UUID          `json:"organization_id"`
	OrganizationName          string             `json:"organization_name"`
	ParentID                  *uuid.UUID         `json:"parent_id"`
	Level                     int                `json:"level"` // Add level for hierarchy depth
	Path                      string             `json:"path"`  // Store full path for easy traversal
	Existing                  int                `json:"existing"`
	Promotion                 int                `json:"promotion"`
	JobLevel                  JobLevelResponse   `json:"job_level"`
	Parent                    *ParentJobResponse `json:"parent"`
	Children                  []JobResponse      `json:"children"`
}

type ParentJobResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
