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

	Organization OrganizationResponse   `json:"organization"`
	EmployeeJob  map[string]interface{} `json:"employee_job"`
}
