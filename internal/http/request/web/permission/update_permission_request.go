package request

type UpdatePermissionRequest struct {
	ID            string `form:"id" validate:"required"`
	Name          string `form:"name" validate:"required"`
	GuardName     string `form:"guard_name" validate:"required"`
	Label         string `form:"label" validate:"required"`
	ApplicationID string `form:"application_id" validate:"required"`
}
