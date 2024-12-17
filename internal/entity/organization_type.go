package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationType struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	Name       string    `json:"name"`
	MidsuitID  string    `json:"midsuit_id"`
	Category   string    `json:"category"`

	Organizations []Organization `json:"organizations" gorm:"foreignKey:OrganizationTypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (organizationType *OrganizationType) BeforeCreate(tx *gorm.DB) (err error) {
	organizationType.ID = uuid.New()
	organizationType.CreatedAt = time.Now().Add(time.Hour * 7)
	organizationType.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (organizationType *OrganizationType) BeforeUpdate(tx *gorm.DB) (err error) {
	organizationType.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (OrganizationType) TableName() string {
	return "organization_types"
}
