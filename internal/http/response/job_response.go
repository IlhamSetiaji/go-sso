package response

import (
	"github.com/google/uuid"
)

type JobResponse struct {
	ID                      uuid.UUID     `json:"id"`
	Name                    string        `json:"name"`
	OrganizationStructureID uuid.UUID     `json:"organization_structure_id"`
	ParentID                *uuid.UUID    `json:"parent_id"`
	Level                   int           `json:"level"` // Add level for hierarchy depth
	Path                    string        `json:"path"`  // Store full path for easy traversal
	Children                []JobResponse `json:"children"`
}
