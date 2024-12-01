package entity

import (
	"time"

	"github.com/google/uuid"
)

type RolePermission struct {
	RoleID       uuid.UUID `json:"role_id" gorm:"type:char(36);primaryKey"`
	PermissionID uuid.UUID `json:"permission_id" gorm:"type:char(36);primaryKey"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Role       Role       `json:"role" gorm:"foreignKey:RoleID;references:ID"`
	Permission Permission `json:"permission" gorm:"foreignKey:PermissionID;references:ID"`
}

func (rp *RolePermission) BeforeCreate() (err error) {
	rp.CreatedAt = time.Now().Add(time.Hour * 7)
	rp.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (rp *RolePermission) BeforeUpdate() (err error) {
	rp.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
