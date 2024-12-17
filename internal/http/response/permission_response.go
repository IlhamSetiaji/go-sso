package response

import (
	"time"

	"github.com/google/uuid"
)

type PermissionResponse struct {
	ID              uuid.UUID `json:"id"`
	ApplicationID   uuid.UUID `json:"application_id"`
	ApplicationName string    `json:"application_name"`
	Name            string    `json:"name"`
	Label           string    `json:"label"`
	GuardName       string    `json:"guard_name"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
