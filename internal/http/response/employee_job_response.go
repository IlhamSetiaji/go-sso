package response

import (
	"time"

	"github.com/google/uuid"
)

type EmployeeJobResponse struct {
	ID                      uuid.UUID `json:"id"`
	EmpOrganizationID       uuid.UUID `json:"emp_organization_id"`
	JobID                   uuid.UUID `json:"job_id"`
	JobName                 string    `json:"job_name"`
	JobNameChinese          string    `json:"job_name_chinese"`
	EmployeeID              uuid.UUID `json:"employee_id"`
	OrganizationLocationID  uuid.UUID `json:"organization_location_id"`
	OrganizationStructureID uuid.UUID `json:"organization_structure_id"`
	Name                    string    `json:"name"`
	JobLevelID              string    `json:"job_level_id"`
	MidsuitID               string    `json:"midsuit_id"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`

	EmpOrganization       *OrganizationMinimalResponse   `json:"emp_organization"`
	Job                   *JobResponse                   `json:"job"`
	OrganizationLocation  *OrganizationLocationResponse  `json:"organization_location"`
	OrganizationStructure *OrganizationStructureResponse `json:"organization_structure"`
	Grade                 *GradeResponse                 `json:"grade"`
}
