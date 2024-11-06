package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserRole struct {
	UserID    uuid.UUID `json:"user_id" gorm:"type:char(36);primaryKey"`
	RoleID    uuid.UUID `json:"role_id" gorm:"type:char(36);primaryKey"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	User User `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Role Role `json:"role" gorm:"foreignKey:RoleID;references:ID"`
}

func (UserRole) TableName() string {
	return "user_roles"
}
