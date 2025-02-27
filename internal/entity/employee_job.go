package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmployeeJob struct {
	gorm.Model              `json:"-"`
	ID                      uuid.UUID  `json:"id" gorm:"type:char(36);primaryKey"`
	EmpOrganizationID       uuid.UUID  `json:"emp_organization_id" gorm:"type:char(36)"`
	JobID                   uuid.UUID  `json:"job_id" gorm:"type:char(36)"`
	EmployeeID              *uuid.UUID `json:"employee_id" gorm:"type:char(36); unique"`
	OrganizationLocationID  uuid.UUID  `json:"organization_location_id" gorm:"type:char(36)"`
	OrganizationStructureID uuid.UUID  `json:"organization_structure_id" gorm:"type:char(36)"`
	Name                    string     `json:"name" gorm:"type:varchar(255)"`

	EmpOrganization       *Organization          `json:"emp_organization" gorm:"foreignKey:EmpOrganizationID;references:ID;constraint:OnDelete:CASCADE"`
	Job                   *Job                   `json:"job" gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE"`
	Employee              *Employee              `json:"employee" gorm:"foreignKey:EmployeeID;references:ID;constraint:OnDelete:CASCADE"`
	OrganizationLocation  *OrganizationLocation  `json:"organization_location" gorm:"foreignKey:OrganizationLocationID;references:ID;constraint:OnDelete:CASCADE"`
	OrganizationStructure *OrganizationStructure `json:"organization_structure" gorm:"foreignKey:OrganizationStructureID;references:ID;constraint:OnDelete:CASCADE"`
}

func (e *EmployeeJob) BeforeCreate(tx *gorm.DB) (err error) {
	e.ID = uuid.New()
	e.CreatedAt = time.Now().Add(time.Hour * 7)
	e.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (e *EmployeeJob) BeforeUpdate(tx *gorm.DB) (err error) {
	e.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (EmployeeJob) TableName() string {
	return "employee_jobs"
}
