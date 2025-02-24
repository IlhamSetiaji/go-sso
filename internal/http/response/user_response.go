package response

import (
	"app/go-sso/internal/entity"
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID                  uuid.UUID              `json:"id"`
	VerifiedUserProfile string                 `json:"verified_user_profile"`
	ChoosedRole         string                 `json:"choosed_role"`
	EmployeeID          uuid.UUID              `json:"employee_id"`
	OauthID             string                 `json:"oauth_id"`
	Username            string                 `json:"username"`
	Email               string                 `json:"email"`
	Name                string                 `json:"name"`
	MobilePhone         string                 `json:"mobile_phone"`
	Address             string                 `json:"address"`
	EmailVerifiedAt     time.Time              `json:"email_verified_at"`
	Gender              entity.UserGender      `json:"gender"`
	Photo               string                 `json:"photo"`
	Status              entity.UserStatus      `json:"status"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
	Roles               []RoleResponse         `json:"roles"`
	UserProfile         map[string]interface{} `json:"user_profile"`

	Employee EmployeeResponse `json:"employee"`
}
