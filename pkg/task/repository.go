package task

import (
	"errors"
	"sort"
	"strings"
	"sync"
)

var (
	// ErrTaskNotFound is returned when a task is not found
	ErrTaskNotFound = errors.New("task not found")
)

// Repository defines the interface for task storage operations
type Repository interface {
	Create(task *Task) error
	GetByID(id string) (*Task, error)
	GetAll(filter FilterOptions) ([]*Task, error)
	Update(task *Task) error
	Delete(id string) error
}

// InMemoryRepository implements Repository using in-memory storage
type InMemoryRepository struct {
	mu    sync.RWMutex
	tasks map[string]*Task
}

// NewInMemoryRepository creates a new in-memory repository
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		tasks: make(map[string]*Task),
	}
}

// Create adds a new task to the repository
func (r *InMemoryRepository) Create(task *Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tasks[task.ID] = task
	return nil
}

// GetByID retrieves a task by its ID
func (r *InMemoryRepository) GetByID(id string) (*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, ErrTaskNotFound
	}

	return task, nil
}

// GetAll retrieves all tasks with optional filtering
func (r *InMemoryRepository) GetAll(filter FilterOptions) ([]*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*Task
	for _, task := range r.tasks {
		if matchesFilter(task, filter) {
			result = append(result, task)
		}
	}

	// Sort tasks
	sortTasks(result, filter)

	// Apply pagination
	start := filter.Offset
	if start > len(result) {
		return []*Task{}, nil
	}

	end := start + filter.Limit
	if filter.Limit == 0 || end > len(result) {
		end = len(result)
	}

	return result[start:end], nil
}

// Update modifies an existing task
func (r *InMemoryRepository) Update(task *Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.ID]; !exists {
		return ErrTaskNotFound
	}

	r.tasks[task.ID] = task
	return nil
}

// Delete removes a task from the repository
func (r *InMemoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return ErrTaskNotFound
	}

	delete(r.tasks, id)
	return nil
}

// matchesFilter checks if a task matches the given filter criteria
func matchesFilter(task *Task, filter FilterOptions) bool {
	if filter.Status != nil && task.Status != *filter.Status {
		return false
	}

	if filter.Priority != nil && task.Priority != *filter.Priority {
		return false
	}

	if filter.Tag != nil {
		found := false
		for _, tag := range task.Tags {
			if tag == *filter.Tag {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// sortTasks sorts tasks based on the filter options
func sortTasks(tasks []*Task, filter FilterOptions) {
	if filter.SortBy == "" {
		filter.SortBy = "created_at"
	}

	sort.Slice(tasks, func(i, j int) bool {
		var less bool
		switch filter.SortBy {
		case "created_at":
			less = tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
		case "updated_at":
			less = tasks[i].UpdatedAt.Before(tasks[j].UpdatedAt)
		case "due_date":
			if tasks[i].DueDate == nil && tasks[j].DueDate == nil {
				less = false
			} else if tasks[i].DueDate == nil {
				less = false
			} else if tasks[j].DueDate == nil {
				less = true
			} else {
				less = tasks[i].DueDate.Before(*tasks[j].DueDate)
			}
		case "priority":
			less = priorityValue(tasks[i].Priority) < priorityValue(tasks[j].Priority)
		case "title":
			less = strings.ToLower(tasks[i].Title) < strings.ToLower(tasks[j].Title)
		default:
			less = tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
		}

		if filter.SortDesc {
			return !less
		}
		return less
	})
}

// priorityValue returns numeric value for priority comparison
func priorityValue(p Priority) int {
	switch p {
	case PriorityLow:
		return 1
	case PriorityMedium:
		return 2
	case PriorityHigh:
		return 3
	default:
		return 0
	}
}
