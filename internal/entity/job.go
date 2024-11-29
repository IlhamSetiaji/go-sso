package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	ID                      uuid.UUID             `json:"id" gorm:"type:char(36);primaryKey"`
	Name                    string                `json:"name"`
	OrganizationStructureID uuid.UUID             `json:"organization_structure_id" gorm:"type:char(36)"`
	ParentID                *uuid.UUID            `json:"parent_id" gorm:"type:char(36)"`
	OrganizationStructure   OrganizationStructure `json:"organization_structure" gorm:"foreignKey:OrganizationStructureID;references:ID;constraint:OnDelete:CASCADE"`
	Parent                  *Job                  `json:"parent" gorm:"foreignKey:ParentID;references:ID;constraint:OnDelete:CASCADE"`
	Children                []Job                 `json:"children" gorm:"foreignKey:ParentID;references:ID"`
	EmployeeJobs            []EmployeeJob         `json:"employee_jobs" gorm:"foreignKey:JobID;references:ID"`
}

func (job *Job) BeforeCreate(tx *gorm.DB) (err error) {
	job.ID = uuid.New()
	job.CreatedAt = time.Now().Add(time.Hour * 7)
	job.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (job *Job) BeforeUpdate(tx *gorm.DB) (err error) {
	job.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (Job) TableName() string {
	return "jobs"
}
