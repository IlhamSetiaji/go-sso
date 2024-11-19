package request

type LoginWebRequest struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}
