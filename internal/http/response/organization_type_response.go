package response

import "github.com/google/uuid"

type OrganizationTypeResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
