package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthToken struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:char(36);not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Token     string    `json:"token" gorm:"type:text;not null"`
	ExpiredAt time.Time `json:"expired_at" gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (authToken *AuthToken) BeforeCreate(tx *gorm.DB) (err error) {
	authToken.ID = uuid.New()
	return
}

func (AuthToken) TableName() string {
	return "auth_tokens"
}
