package task

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrInvalidTitle is returned when title is empty
	ErrInvalidTitle = errors.New("title is required")
	// ErrInvalidStatus is returned when status is invalid
	ErrInvalidStatus = errors.New("invalid status")
	// ErrInvalidPriority is returned when priority is invalid
	ErrInvalidPriority = errors.New("invalid priority")
)

// Service handles business logic for tasks
type Service struct {
	repo Repository
}

// NewService creates a new task service
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateTaskRequest represents the data needed to create a task
type CreateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Priority    string     `json:"priority,omitempty"`
	Status      string     `json:"status,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
}

// UpdateTaskRequest represents the data for updating a task
type UpdateTaskRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Priority    *string    `json:"priority,omitempty"`
	Status      *string    `json:"status,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
}

// ValidateCreateRequest validates a create task request (pure function)
func ValidateCreateRequest(req CreateTaskRequest) error {
	if strings.TrimSpace(req.Title) == "" {
		return ErrInvalidTitle
	}

	if req.Status != "" && !IsValidStatus(req.Status) {
		return ErrInvalidStatus
	}

	if req.Priority != "" && !IsValidPriority(req.Priority) {
		return ErrInvalidPriority
	}

	return nil
}

// BuildTaskFromRequest creates a Task from CreateTaskRequest (pure function)
func BuildTaskFromRequest(req CreateTaskRequest, now time.Time) *Task {
	status := StatusPending
	if req.Status != "" {
		status = Status(req.Status)
	}

	priority := PriorityMedium
	if req.Priority != "" {
		priority = Priority(req.Priority)
	}

	return &Task{
		ID:          uuid.New().String(),
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
		Status:      status,
		Priority:    priority,
		DueDate:     req.DueDate,
		Tags:        req.Tags,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// ApplyUpdate applies an UpdateTaskRequest to an existing Task (pure function)
func ApplyUpdate(task *Task, req UpdateTaskRequest, now time.Time) (*Task, error) {
	updated := *task

	if req.Title != nil {
		trimmed := strings.TrimSpace(*req.Title)
		if trimmed == "" {
			return nil, ErrInvalidTitle
		}
		updated.Title = trimmed
	}

	if req.Description != nil {
		updated.Description = strings.TrimSpace(*req.Description)
	}

	if req.Status != nil {
		if !IsValidStatus(*req.Status) {
			return nil, ErrInvalidStatus
		}
		updated.Status = Status(*req.Status)
	}

	if req.Priority != nil {
		if !IsValidPriority(*req.Priority) {
			return nil, ErrInvalidPriority
		}
		updated.Priority = Priority(*req.Priority)
	}

	if req.DueDate != nil {
		updated.DueDate = req.DueDate
	}

	if req.Tags != nil {
		updated.Tags = req.Tags
	}

	updated.UpdatedAt = now

	return &updated, nil
}

// Create creates a new task
func (s *Service) Create(req CreateTaskRequest) (*Task, error) {
	if err := ValidateCreateRequest(req); err != nil {
		return nil, err
	}

	task := BuildTaskFromRequest(req, time.Now())

	if err := s.repo.Create(task); err != nil {
		return nil, err
	}

	return task, nil
}

// GetByID retrieves a task by ID
func (s *Service) GetByID(id string) (*Task, error) {
	return s.repo.GetByID(id)
}

// GetAll retrieves all tasks with optional filters
func (s *Service) GetAll(filter FilterOptions) ([]*Task, error) {
	return s.repo.GetAll(filter)
}

// Update updates an existing task
func (s *Service) Update(id string, req UpdateTaskRequest) (*Task, error) {
	existingTask, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	updatedTask, err := ApplyUpdate(existingTask, req, time.Now())
	if err != nil {
		return nil, err
	}

	if err := s.repo.Update(updatedTask); err != nil {
		return nil, err
	}

	return updatedTask, nil
}

// Delete removes a task by ID
func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
