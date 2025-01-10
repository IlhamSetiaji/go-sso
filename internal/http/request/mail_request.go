package request

type MailRequest struct {
	Email   string `json:"email" validate:"required,email"`
	From    string `json:"from" validate:"required,email"`
	To      string `json:"to" validate:"required,email"`
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body" validate:"required"`
}
