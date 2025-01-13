package request

type VerifyEmailRequest struct {
	Email string `form:"email" validate:"required,email"`
	Token int    `form:"token" validate:"required"`
}
