package response

import (
	"time"

	"github.com/google/uuid"
)

type RoleResponse struct {
	ID              uuid.UUID            `json:"id"`
	ApplicationID   uuid.UUID            `json:"application_id"`
	ApplicationName string               `json:"application_name"`
	Name            string               `json:"name"`
	GuardName       string               `json:"guard_name"`
	Status          string               `json:"status"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
	Permissions     []PermissionResponse `json:"permissions"`
}
