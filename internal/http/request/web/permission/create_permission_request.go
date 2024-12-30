package request

type CreatePermissionRequest struct {
	Name          string `form:"name" validate:"required"`
	GuardName     string `form:"guard_name" validate:"required"`
	Label         string `form:"label" validate:"required"`
	ApplicationID string `form:"application_id" validate:"required"`
	Description   string `form:"description" validate:"omitempty"`
}
