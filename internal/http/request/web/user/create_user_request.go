package request

import (
	"app/go-sso/internal/entity"
)

type CreateUserRequest struct {
	Name        string            `form:"name" validate:"required"`
	Email       string            `form:"email" validate:"required,email"`
	Username    string            `form:"username" validate:"required"`
	Gender      entity.UserGender `form:"gender" validate:"required,userGender"`
	MobilePhone string            `form:"mobile_phone" validate:"omitempty,numeric,min=10,max=13,startswith=62"`
	RoleID      string            `form:"role_id" validate:"required,uuid"`
	Status      entity.UserStatus `form:"status" validate:"required,userStatus"`
}
