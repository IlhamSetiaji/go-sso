package request

import "app/go-sso/internal/entity"

type UpdateRoleRequest struct {
	ID            string            `form:"id" validate:"required"`
	Name          string            `form:"name" validate:"required"`
	GuardName     string            `form:"guard_name" validate:"required"`
	Status        entity.RoleStatus `form:"status" validate:"required,roleStatus"`
	ApplicationID string            `form:"application_id" validate:"required"`
}
