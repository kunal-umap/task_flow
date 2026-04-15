package models

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string
type TaskPriority string

const (
	StatusTodo       TaskStatus   = "todo"
	StatusInProgress TaskStatus   = "in_progress"
	StatusDone       TaskStatus   = "done"
	PriorityLow      TaskPriority = "low"
	PriorityMedium   TaskPriority = "medium"
	PriorityHigh     TaskPriority = "high"
)

type Task struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	Title       string       `json:"title" db:"title"`
	Description string       `json:"description,omitempty" db:"description"`
	Status      TaskStatus   `json:"status" db:"status"`
	Priority    TaskPriority `json:"priority" db:"priority"`
	ProjectID   uuid.UUID    `json:"project_id" db:"project_id"`
	AssigneeID  *uuid.UUID   `json:"assignee_id,omitempty" db:"assignee_id"`
	DueDate     *time.Time   `json:"due_date,omitempty" db:"due_date"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
}
