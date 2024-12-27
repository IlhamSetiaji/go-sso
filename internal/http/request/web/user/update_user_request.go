package request

import "app/go-sso/internal/entity"

type UpdateUserRequest struct {
	ID          string            `form:"id" validate:"required"`
	Name        string            `form:"name" validate:"required"`
	Email       string            `form:"email" validate:"required,email"`
	Username    string            `form:"username" validate:"required"`
	Gender      entity.UserGender `form:"gender" validate:"required,userGender"`
	MobilePhone string            `form:"mobile_phone" validate:"omitempty,numeric,min=10,max=13,startswith=62"`
	Status      entity.UserStatus `form:"status" validate:"required,userStatus"`
	RoleID      string            `form:"role_id" validate:"omitempty,uuid"`
}
