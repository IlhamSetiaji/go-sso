package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleStatus string

const (
	ROLE_ACTIVE   RoleStatus = "ACTIVE"
	ROLE_INACTIVE RoleStatus = "INACTIVE"
)

type Role struct {
	// gorm.Model
	ID            uuid.UUID    `json:"id" gorm:"type:char(36);primaryKey"`
	ApplicationID uuid.UUID    `json:"application_id" gorm:"type:char(36);not null"`
	Application   Application  `json:"application" gorm:"foreignKey:ApplicationID;references:ID;constraint:OnDelete:CASCADE"`
	Name          string       `json:"name" gorm:"not null"`
	GuardName     string       `json:"guard_name" gorm:"default:web"`
	Status        RoleStatus   `json:"status" gorm:"default:ACTIVE"`
	Users         []User       `json:"users" gorm:"many2many:user_roles;"`             // many to many relationship
	Permissions   []Permission `json:"permissions" gorm:"many2many:role_permissions;"` // many to many relationship
	CreatedAt     time.Time    `gorm:"autoCreateTime"`
	UpdatedAt     time.Time    `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt
}

func (role *Role) BeforeCreate(tx *gorm.DB) (err error) {
	role.ID = uuid.New()
	role.CreatedAt = time.Now().Add(time.Hour * 7)
	role.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (role *Role) BeforeUpdate(tx *gorm.DB) (err error) {
	role.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (Role) TableName() string {
	return "roles"
}
