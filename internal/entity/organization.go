package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Organization struct {
	gorm.Model
	ID                     uuid.UUID               `json:"id" gorm:"type:char(36);primaryKey"`
	Name                   string                  `json:"name"`
	OrganizationLocations  []OrganizationLocation  `json:"organization_locations" gorm:"foreignKey:OrganizationID;references:ID"`
	OrganizationStructures []OrganizationStructure `json:"organization_structures" gorm:"foreignKey:OrganizationID;references:ID"`
	Employees              []Employee              `json:"employees" gorm:"foreignKey:OrganizationID;references:ID"`
	AlternateEmployeeJobs  []EmployeeJob           `json:"alternate_employee_jobs" gorm:"foreignKey:EmpOrganizationID;references:ID"`
}

func (organization *Organization) BeforeCreate(tx *gorm.DB) (err error) {
	organization.ID = uuid.New()
	organization.CreatedAt = time.Now().Add(time.Hour * 7)
	organization.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (organization *Organization) BeforeUpdate(tx *gorm.DB) (err error) {
	organization.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (Organization) TableName() string {
	return "organizations"
}
