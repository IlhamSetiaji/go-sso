package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationLocation struct {
	gorm.Model
	ID             uuid.UUID     `json:"id" gorm:"type:char(36);primaryKey"`
	OrganizationID uuid.UUID     `json:"organization_id" gorm:"type:char(36)"`
	Organization   Organization  `json:"organization" gorm:"foreignKey:OrganizationID;references:ID;constraint:OnDelete:CASCADE"`
	Name           string        `json:"name"`
	EmployeeJobs   []EmployeeJob `json:"employee_jobs" gorm:"foreignKey:OrganizationLocationID;references:ID;constraint:OnDelete:CASCADE"`
}

func (organizationLocation *OrganizationLocation) BeforeCreate(tx *gorm.DB) (err error) {
	organizationLocation.ID = uuid.New()
	organizationLocation.CreatedAt = time.Now().Add(time.Hour * 7)
	organizationLocation.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (organizationLocation *OrganizationLocation) BeforeUpdate(tx *gorm.DB) (err error) {
	organizationLocation.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (OrganizationLocation) TableName() string {
	return "organization_locations"
}
