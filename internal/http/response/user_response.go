package response

import (
	"app/go-sso/internal/entity"
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID              uuid.UUID         `json:"id"`
	EmployeeID      uuid.UUID         `json:"employee_id"`
	OauthID         string            `json:"oauth_id"`
	Username        string            `json:"username"`
	Email           string            `json:"email"`
	Name            string            `json:"name"`
	MobilePhone     string            `json:"mobile_phone"`
	EmailVerifiedAt time.Time         `json:"email_verified_at"`
	Gender          entity.UserGender `json:"gender"`
	Photo           string            `json:"photo"`
	Status          entity.UserStatus `json:"status"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`

	Employee EmployeeResponse `json:"employee"`
}
