package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model    `json:"-"`
	ID            uuid.UUID   `json:"id" gorm:"type:char(36);primary_key"`
	ApplicationID uuid.UUID   `json:"application_id" gorm:"type:char(36);not null"`
	Application   Application `json:"application" gorm:"foreignKey:ApplicationID;references:ID;constraint:OnDelete:CASCADE"`
	Name          string      `json:"name" gorm:"type:varchar(255);not null"`
	Label         string      `json:"label" gorm:"type:varchar(255);not null"`
	GuardName     string      `json:"guard_name" gorm:"type:varchar(255);not null;default:'web'"` // default guard name is web
	Description   string      `json:"description" gorm:"type:text;default:null"`
	Roles         []Role      `json:"roles" gorm:"many2many:role_permissions;"` // many to many relationship
}

func (permission *Permission) BeforeCreate(tx *gorm.DB) (err error) {
	permission.ID = uuid.New()
	permission.CreatedAt = time.Now().Add(time.Hour * 7)
	permission.UpdatedAt = time.Now().Add(time.Hour * 7)
	return
}

func (permission *Permission) BeforeUpdate(tx *gorm.DB) (err error) {
	permission.UpdatedAt = time.Now().Add(time.Hour * 7)
	return
}

func (Permission) TableName() string {
	return "permissions"
}
