package request

import (
	"app/go-sso/internal/entity"

	"github.com/go-playground/validator/v10"
)

func UserStatusValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true
	}
	switch entity.UserStatus(status) {
	case entity.USER_ACTIVE, entity.USER_INACTIVE, entity.USER_PENDING:
		return true
	default:
		return false
	}
}

func UserGenderValidation(fl validator.FieldLevel) bool {
	gender := fl.Field().String()
	if gender == "" {
		return true
	}
	switch entity.UserGender(gender) {
	case entity.MALE, entity.FEMALE:
		return true
	default:
		return false
	}
}

func RoleStatusValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true
	}
	switch entity.RoleStatus(status) {
	case entity.ROLE_ACTIVE, entity.ROLE_INACTIVE:
		return true
	default:
		return false
	}
}

type RabbitMQRequest struct {
	ID          string                 `json:"id"`
	MessageType string                 `json:"message_type"`
	MessageData map[string]interface{} `json:"message_data"`
	ReplyTo     string                 `json:"reply_to"`
}
