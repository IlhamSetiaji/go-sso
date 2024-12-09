package response

import "github.com/google/uuid"

type JobLevelResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Level string    `json:"level"`
}
