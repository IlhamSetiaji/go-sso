package request

import "app/go-sso/internal/entity"

type CreateRoleRequest struct {
	Name          string            `form:"name" validate:"required"`
	GuardName     string            `form:"guard_name" validate:"required"`
	ApplicationID string            `form:"application_id" validate:"required"`
	Status        entity.RoleStatus `form:"status" validate:"required,roleStatus"`
}
