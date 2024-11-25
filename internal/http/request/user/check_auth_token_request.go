package request

type CheckAuthTokenRequest struct {
	UserID string `json:"user_id" validate:"required"`
	Token  string `json:"token" validate:"required"`
}
