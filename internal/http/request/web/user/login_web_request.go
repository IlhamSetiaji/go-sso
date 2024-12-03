package request

type LoginWebRequest struct {
	Email    string `form:"email" validate:"required"`
	Password string `form:"password" validate:"required"`
	State    string `form:"state" validate:"omitempty"`
}
