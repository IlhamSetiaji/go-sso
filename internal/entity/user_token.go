package entity

import (
	"time"

	"gorm.io/gorm"
)

type UserTokenType string

const (
	UserTokenVerification  UserTokenType = "VERIFICATION"
	UserTokenResetPassword UserTokenType = "RESET_PASSWORD"
)

type UserToken struct {
	Email     string        `json:"email" gorm:"type:varchar(255);not null"`
	Token     int           `json:"token" gorm:"type:int;not null"`
	TokenType UserTokenType `json:"token_type" gorm:"not null"`
	ExpiredAt time.Time     `json:"expired_at"`
	CreatedAt time.Time     `gorm:"autoCreateTime"`
	UpdatedAt time.Time     `gorm:"autoUpdateTime"`
}

func (userToken *UserToken) BeforeCreate(tx *gorm.DB) (err error) {
	userToken.ExpiredAt = time.Now().Add(time.Hour * 3)
	userToken.CreatedAt = time.Now()
	userToken.UpdatedAt = time.Now()
	return nil
}

func (UserToken) TableName() string {
	return "user_tokens"
}
