package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobLevel struct {
	gorm.Model
	ID                     uuid.UUID               `json:"id" gorm:"type:char(36);primaryKey"`
	Name                   string                  `json:"name" gorm:"type:varchar(255);not null"`
	Level                  string                  `json:"level" gorm:"type:varchar(10);not null"`
	OrganizationStructures []OrganizationStructure `json:"employee_jobs" gorm:"foreignKey:JobLevelID;references:ID"`
}
