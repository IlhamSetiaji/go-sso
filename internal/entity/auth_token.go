package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthToken struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:char(36);not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Token     string    `json:"token" gorm:"type:text;not null"`
	ExpiredAt time.Time `json:"expired_at" gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (authToken *AuthToken) BeforeCreate(tx *gorm.DB) (err error) {
	authToken.ID = uuid.New()
	authToken.CreatedAt = time.Now().Add(time.Hour * 7)
	authToken.UpdatedAt = time.Now().Add(time.Hour * 7)
	if !authToken.ExpiredAt.IsZero() {
		authToken.ExpiredAt = authToken.ExpiredAt.Add(time.Hour * 7)
	}
	return nil
}

func (authToken *AuthToken) BeforeUpdate(tx *gorm.DB) (err error) {
	authToken.UpdatedAt = time.Now().Add(time.Hour * 7)
	if !authToken.ExpiredAt.IsZero() {
		authToken.ExpiredAt = authToken.ExpiredAt.Add(time.Hour * 7)
	}
	return nil
}

func (AuthToken) TableName() string {
	return "auth_tokens"
}
