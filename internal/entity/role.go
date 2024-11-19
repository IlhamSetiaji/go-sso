package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	// gorm.Model
	ID          uuid.UUID    `json:"id" gorm:"type:char(36);primaryKey"`
	ClientID    uuid.UUID    `json:"client_id" gorm:"type:char(36);not null"`
	Client      Client       `json:"client" gorm:"foreignKey:ClientID;references:ID"`
	Name        string       `json:"name" gorm:"not null"`
	GuardName   string       `json:"guard_name" gorm:"default:web"`
	Users       []User       `json:"users" gorm:"many2many:user_roles;"`             // many to many relationship
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"` // many to many relationship
	CreatedAt   time.Time    `gorm:"autoCreateTime"`
	UpdatedAt   time.Time    `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt
}

func (role *Role) BeforeCreate(tx *gorm.DB) (err error) {
	role.ID = uuid.New()
	return
}

func (Role) TableName() string {
	return "roles"
}
