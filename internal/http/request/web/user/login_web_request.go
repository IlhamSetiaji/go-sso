package request

type LoginWebRequest struct {
	Email    string `form:"email" validate:"required"`
	Password string `form:"password" validate:"required"`
	State    string `form:"state" validate:"omitempty"`
}

type ChooseRolesWebRequest struct {
	ID     string `form:"id" validate:"required"`
	RoleID string `form:"role_id" validate:"required"`
	State  string `form:"state" validate:"omitempty"`
}
