package config

import (
	"app/go-sso/internal/http/request"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func NewValidator(viper *viper.Viper) *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("userStatus", request.UserStatusValidation)
	return validate
}
