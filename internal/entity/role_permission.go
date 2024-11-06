package entity

import (
	"time"

	"github.com/google/uuid"
)

type RolePermission struct {
	RoleID       uuid.UUID `json:"role_id" gorm:"type:char(36);primaryKey"`
	PermissionID uuid.UUID `json:"permission_id" gorm:"type:char(36);primaryKey"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Role       Role       `json:"role" gorm:"foreignKey:RoleID;references:ID"`
	Permission Permission `json:"permission" gorm:"foreignKey:PermissionID;references:ID"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
