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
