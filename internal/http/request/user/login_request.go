package request

type LoginRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required"`
	ChoosedRole string `json:"choosed_role"`
}
