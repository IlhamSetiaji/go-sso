package entity

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Username        string    `json:"username"`
	IsEmployee      bool      `json:"is_employee"`
	EmailVerifiedAt time.Time `json:"email_verified_at"`
}
