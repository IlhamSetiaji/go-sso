package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model     `json:"-"`
	ID             uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	OrganizationID uuid.UUID `json:"organization_id" gorm:"type:char(36)"`
	Name           string    `json:"name"`
	EndDate        time.Time `json:"end_date" gorm:"type:date"`
	RetirementDate time.Time `json:"retirement_date" gorm:"type:date"`
	Email          string    `json:"email" gorm:"default:null"`
	MobilePhone    string    `json:"mobile_phone" gorm:"default:null"`
	MidsuitID      string    `json:"midsuit_id"`
	SignaturePath  string    `json:"signature_path" gorm:"type:text"`
	IsCeoPic       bool      `json:"is_ceo_pic" gorm:"type:boolean;default:null"`
	IsOnboarding   string    `json:"is_onboarding" gorm:"type:boolean;default:YES"`

	Organization Organization `json:"organization" gorm:"foreignKey:OrganizationID;references:ID;constraint:OnDelete:CASCADE"`
	User         *User        `json:"user" gorm:"foreignKey:EmployeeID;references:ID"`
	EmployeeJob  *EmployeeJob `json:"employee_job" gorm:"foreignKey:EmployeeID;references:ID"`

	EmployeeKanbanProgress *EmployeeKanbanProgress `json:"employee_kanban_progress" gorm:"-"`
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

func (employee *Employee) BeforeDelete(tx *gorm.DB) (err error) {
	// Modify unique columns before soft delete
	if employee.DeletedAt.Valid {
		return nil
	}

	randomString := uuid.New().String()

	employee.Email = employee.Email + "_deleted" + randomString
	employee.MobilePhone = employee.MobilePhone + "_deleted" + randomString

	// Update the record with the modified unique columns
	if err := tx.Model(&employee).Where("id = ?", employee.ID).Updates(map[string]interface{}{
		"email":        employee.Email,
		"mobile_phone": employee.MobilePhone,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (Employee) TableName() string {
	return "employees"
}
