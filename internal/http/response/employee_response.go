package response

import (
	"time"

	"github.com/google/uuid"
)

type EmployeeResponse struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
	EndDate        time.Time `json:"end_date"`
	RetirementDate time.Time `json:"retirement_date"`
	Email          string    `json:"email"`
	MobilePhone    string    `json:"mobile_phone"`
	NIK            string    `json:"nik"`
	MidsuitID      string    `json:"midsuit_id"`

	Organization   OrganizationResponse            `json:"organization"`
	EmployeeJob    *EmployeeJobResponse            `json:"employee_job"`
	KanbanProgress *EmployeeKanbanProgressResponse `json:"kanban_progress"`
}

type EmployeeKanbanProgressResponse struct {
	TotalTask  int `json:"total_task"`
	ToDo       int `json:"to_do"`
	InProgress int `json:"in_progress"`
	NeedReview int `json:"need_review"`
	Completed  int `json:"completed"`
}
