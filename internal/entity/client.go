package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Client struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	Name        string    `json:"name" gorm:"unique;not null"`
	Secret      string    `json:"secret" gorm:"unique;not null"`
	RedirectURI string    `json:"redirect_uri" gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt
	Roles       []Role `json:"roles" gorm:"foreignKey:ClientID;references:ID"`
}

func (client *Client) BeforeCreate(tx *gorm.DB) (err error) {
	client.ID = uuid.New()
	return
}

func (Client) TableName() string {
	return "clients"
}
