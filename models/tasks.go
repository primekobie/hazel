package models

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// TaskStatus defines allowed statuses for a task.
type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "started"
	StatusDone       TaskStatus = "complete"
)

// TaskPriority defines allowed priorities for a task.
type TaskPriority string

const (
	PriorityLow    TaskPriority = "low"
	PriorityMedium TaskPriority = "medium"
	PriorityHigh   TaskPriority = "high"
)

// Task represents a single work item within a project.
type Task struct {
	Id          uuid.UUID    `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Project     *Project     `json:"project,omitemtpy"`
	AssignedTo  []User       `json:"assignedTo,omitempty"`
	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
	Due         time.Time    `json:"due,omitzero"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

type TaskStore interface {
	Create(ctx context.Context, task *Task) error
	Update(ctx context.Context, task *Task) error
	Get(ctx context.Context, id uuid.UUID) (Task, error)
	GetAllForProject(ctx context.Context, projectId uuid.UUID)
	Delete(ctx context.Context, id uuid.UUID) error
}
