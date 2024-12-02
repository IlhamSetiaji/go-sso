package request

type DeletePermissionRequest struct {
	ID string `form:"id" validate:"required"`
}
