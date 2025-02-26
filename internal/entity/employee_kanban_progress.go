package entity

type EmployeeKanbanProgress struct {
	TotalTask  int `json:"total_task"`
	ToDo       int `json:"to_do"`
	InProgress int `json:"in_progress"`
	NeedReview int `json:"need_review"`
	Completed  int `json:"completed"`
}
