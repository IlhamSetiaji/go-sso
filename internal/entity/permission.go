package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null;unique"`
	Label     string    `json:"label" gorm:"type:varchar(255);not null"`
	GuardName string    `json:"guard_name" gorm:"type:varchar(255);not null;default:'web'"` // default guard name is web
	Roles     []Role    `json:"roles" gorm:"many2many:role_permissions;"`                   // many to many relationship
}

func (permission *Permission) BeforeCreate(tx *gorm.DB) (err error) {
	permission.ID = uuid.New()
	return
}

func (Permission) TableName() string {
	return "permissions"
}
