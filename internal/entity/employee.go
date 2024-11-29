package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	ID             uuid.UUID    `json:"id" gorm:"type:char(36);primaryKey"`
	OrganizationID uuid.UUID    `json:"organization_id" gorm:"type:char(36)"`
	Name           string       `json:"name"`
	EndDate        *time.Time   `json:"end_date" gorm:"type:date"`
	RetirementDate *time.Time   `json:"retirement_date" gorm:"type:date"`
	Email          string       `json:"email" gorm:"unique"`
	MobilePhone    string       `json:"mobile_phone" gorm:"unique"`
	Organization   Organization `json:"organization" gorm:"foreignKey:OrganizationID;references:ID;constraint:OnDelete:CASCADE"`
}

func (employee *Employee) BeforeCreate(tx *gorm.DB) (err error) {
	employee.ID = uuid.New()
	employee.CreatedAt = time.Now().Add(time.Hour * 7)
	employee.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (employee *Employee) BeforeUpdate(tx *gorm.DB) (err error) {
	employee.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (Employee) TableName() string {
	return "employees"
}
