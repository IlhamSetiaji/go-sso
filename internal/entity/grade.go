package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Grade struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	JobLevelID uuid.UUID `json:"job_level_id" gorm:"type:char(36)"`
	Name       string    `json:"name" gorm:"type:varchar(255)"`
	MidsuitID  string    `json:"midsuit_id" gorm:"type:varchar(255)"`

	JobLevel     *JobLevel     `json:"job_level" gorm:"foreignKey:JobLevelID;references:ID;constraint:OnDelete:CASCADE"`
	EmployeeJobs []EmployeeJob `json:"employee_jobs" gorm:"foreignKey:GradeID;references:ID;constraint:OnDelete:CASCADE"`
}

func (g *Grade) BeforeCreate(tx *gorm.DB) (err error) {
	g.ID = uuid.New()
	g.CreatedAt = time.Now().Add(time.Hour * 7)
	g.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (g *Grade) BeforeUpdate(tx *gorm.DB) (err error) {
	g.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (Grade) TableName() string {
	return "grades"
}
