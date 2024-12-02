package request

type DeleteRoleRequest struct {
	ID string `form:"id" validate:"required"`
}
