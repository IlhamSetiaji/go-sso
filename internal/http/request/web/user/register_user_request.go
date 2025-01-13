package request

import (
	"app/go-sso/internal/entity"
)

type UserRegisterRequest struct {
	Username             string            `form:"username" validate:"required"`
	Email                string            `form:"email" validate:"required,email"`
	Name                 string            `form:"name" validate:"required"`
	Password             string            `form:"password" validate:"required"`
	PasswordConfirmation string            `form:"password_confirmation" validate:"required,eqfield=Password"`
	Gender               entity.UserGender `form:"gender" validate:"required,userGender"`
	MobilePhone          string            `form:"mobile_phone" validate:"required"`
	BirthDate            string            `form:"birth_date" validate:"required"`
	BirthPlace           string            `form:"birth_place" validate:"required"`
	Address              string            `form:"address" validate:"omitempty"`
	NoKTP                string            `form:"no_ktp" validate:"required"`
	KTP                  string            `form:"ktp" validate:"omitempty"`
}
