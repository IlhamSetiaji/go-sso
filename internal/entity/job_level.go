package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobLevel struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	Name       string    `json:"name" gorm:"type:varchar(255);not null"`
	Level      string    `json:"level" gorm:"type:varchar(10);not null"`
	MidsuitID  string    `json:"midsuit_id" gorm:"type:varchar(255)"`

	OrganizationStructures []OrganizationStructure `json:"employee_jobs" gorm:"foreignKey:JobLevelID;references:ID;constraint:OnDelete:CASCADE"`
	Jobs                   []Job                   `json:"jobs" gorm:"foreignKey:JobLevelID;references:ID;constraint:OnDelete:CASCADE"`
	// Grades                 []Grade                 `json:"grades" gorm:"foreignKey:JobLevelID;references:ID constraint:OnDelete:CASCADE"`
}

func (jobLevel *JobLevel) BeforeCreate(tx *gorm.DB) (err error) {
	jobLevel.ID = uuid.New()
	jobLevel.CreatedAt = time.Now().Add(time.Hour * 7)
	jobLevel.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (jobLevel *JobLevel) BeforeUpdate(tx *gorm.DB) (err error) {
	jobLevel.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (JobLevel) TableName() string {
	return "job_levels"
}
