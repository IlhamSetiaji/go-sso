package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserRole struct {
	UserID    uuid.UUID `json:"user_id" gorm:"type:char(36);primaryKey"`
	RoleID    uuid.UUID `json:"role_id" gorm:"type:char(36);primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User User `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Role Role `json:"role" gorm:"foreignKey:RoleID;references:ID"`
}

func (userRole *UserRole) BeforeCreate() (err error) {
	userRole.CreatedAt = time.Now().Add(time.Hour * 7)
	userRole.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (userRole *UserRole) BeforeUpdate() (err error) {
	userRole.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (UserRole) TableName() string {
	return "user_roles"
}
