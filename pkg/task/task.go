package task

import (
	"time"
)

// Status represents the current state of a task
type Status string

const (
	// StatusPending represents a task that hasn't been started
	StatusPending Status = "pending"
	// StatusInProgress represents a task that is currently being worked on
	StatusInProgress Status = "in_progress"
	// StatusCompleted represents a task that has been finished
	StatusCompleted Status = "completed"
)

// Priority represents the importance level of a task
type Priority string

const (
	// PriorityLow represents low priority tasks
	PriorityLow Priority = "low"
	// PriorityMedium represents medium priority tasks
	PriorityMedium Priority = "medium"
	// PriorityHigh represents high priority tasks
	PriorityHigh Priority = "high"
)

// Task represents a todo item
type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Status      Status     `json:"status"`
	Priority    Priority   `json:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// IsValidStatus checks if a status string is valid
func IsValidStatus(s string) bool {
	return s == string(StatusPending) ||
		s == string(StatusInProgress) ||
		s == string(StatusCompleted)
}

// IsValidPriority checks if a priority string is valid
func IsValidPriority(p string) bool {
	return p == string(PriorityLow) ||
		p == string(PriorityMedium) ||
		p == string(PriorityHigh)
}

// FilterOptions contains parameters for filtering tasks
type FilterOptions struct {
	Status   *Status
	Priority *Priority
	Tag      *string
	Limit    int
	Offset   int
	SortBy   string
	SortDesc bool
}
