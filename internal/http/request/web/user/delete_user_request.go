package request

type DeleteUserRequest struct {
	ID string `form:"id" validate:"required"`
}
