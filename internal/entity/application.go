package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Application struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	Name        string    `json:"name" gorm:"unique;not null"`
	Label       string    `json:"label" gorm:"not null"`
	Secret      string    `json:"secret" gorm:"unique;not null"`
	RedirectURI string    `json:"redirect_uri" gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt
	Roles       []Role       `json:"roles" gorm:"foreignKey:ApplicationID;references:ID"`
	Permissions []Permission `json:"permissions" gorm:"foreignKey:ApplicationID;references:ID"`
}

func (application *Application) BeforeCreate(tx *gorm.DB) (err error) {
	application.ID = uuid.New()
	application.CreatedAt = time.Now().Add(time.Hour * 7)
	application.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (application *Application) BeforeUpdate(tx *gorm.DB) (err error) {
	application.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (Application) TableName() string {
	return "applications"
}
