package request

import (
	"app/go-sso/internal/entity"
)

type CreateUserRequest struct {
	EmployeeID  string            `form:"employee_id" validate:"omitempty"`
	Name        string            `form:"name" validate:"required"`
	Email       string            `form:"email" validate:"required,email"`
	Username    string            `form:"username" validate:"required"`
	Gender      entity.UserGender `form:"gender" validate:"required,userGender"`
	MobilePhone string            `form:"mobile_phone" validate:"omitempty,numeric,min=10,max=13,startswith=62"`
	RoleIDs     []string          `form:"role_ids[]" validate:"required,dive"`
	Status      entity.UserStatus `form:"status" validate:"required,userStatus"`
}
