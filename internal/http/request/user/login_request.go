package request

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// package request

// type LoginWebRequest struct {
// 	Email    string `form:"email" binding:"required"`
// 	Password string `form:"password" binding:"required"`
// }
