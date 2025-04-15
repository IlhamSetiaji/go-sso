package response

import (
	"time"

	"github.com/google/uuid"
)

type GradeResponse struct {
	ID         uuid.UUID         `json:"id"`
	JobLevelID uuid.UUID         `json:"job_level_id"`
	Name       string            `json:"name"`
	MidsuitID  string            `json:"midsuit_id"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	JobLevel   *JobLevelResponse `json:"job_level"`
}
