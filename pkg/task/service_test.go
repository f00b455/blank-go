package task

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateCreateRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateTaskRequest
		wantErr error
	}{
		{
			name: "valid request with all fields",
			req: CreateTaskRequest{
				Title:       "Test task",
				Description: "Description",
				Priority:    "high",
				Status:      "pending",
			},
			wantErr: nil,
		},
		{
			name: "valid request with required fields only",
			req: CreateTaskRequest{
				Title: "Test task",
			},
			wantErr: nil,
		},
		{
			name: "empty title",
			req: CreateTaskRequest{
				Title: "",
			},
			wantErr: ErrInvalidTitle,
		},
		{
			name: "whitespace-only title",
			req: CreateTaskRequest{
				Title: "   ",
			},
			wantErr: ErrInvalidTitle,
		},
		{
			name: "invalid status",
			req: CreateTaskRequest{
				Title:  "Test task",
				Status: "invalid",
			},
			wantErr: ErrInvalidStatus,
		},
		{
			name: "invalid priority",
			req: CreateTaskRequest{
				Title:    "Test task",
				Priority: "critical",
			},
			wantErr: ErrInvalidPriority,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateRequest(tt.req)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBuildTaskFromRequest(t *testing.T) {
	now := time.Date(2026, 1, 14, 12, 0, 0, 0, time.UTC)
	dueDate := time.Date(2026, 1, 20, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		req  CreateTaskRequest
		want func(*Task)
	}{
		{
			name: "with all fields",
			req: CreateTaskRequest{
				Title:       "Test task",
				Description: "Description",
				Priority:    "high",
				Status:      "in_progress",
				DueDate:     &dueDate,
				Tags:        []string{"work", "urgent"},
			},
			want: func(task *Task) {
				assert.NotEmpty(t, task.ID)
				assert.Equal(t, "Test task", task.Title)
				assert.Equal(t, "Description", task.Description)
				assert.Equal(t, PriorityHigh, task.Priority)
				assert.Equal(t, StatusInProgress, task.Status)
				assert.Equal(t, &dueDate, task.DueDate)
				assert.Equal(t, []string{"work", "urgent"}, task.Tags)
				assert.Equal(t, now, task.CreatedAt)
				assert.Equal(t, now, task.UpdatedAt)
			},
		},
		{
			name: "with defaults",
			req: CreateTaskRequest{
				Title: "  Test task  ",
			},
			want: func(task *Task) {
				assert.Equal(t, "Test task", task.Title)
				assert.Equal(t, "", task.Description)
				assert.Equal(t, PriorityMedium, task.Priority)
				assert.Equal(t, StatusPending, task.Status)
				assert.Nil(t, task.DueDate)
				assert.Nil(t, task.Tags)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := BuildTaskFromRequest(tt.req, now)
			tt.want(task)
		})
	}
}

func TestApplyUpdate(t *testing.T) {
	now := time.Date(2026, 1, 14, 12, 0, 0, 0, time.UTC)
	updateTime := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
	originalTask := &Task{
		ID:          "test-id",
		Title:       "Original Title",
		Description: "Original Description",
		Status:      StatusPending,
		Priority:    PriorityLow,
		Tags:        []string{"old"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tests := []struct {
		name    string
		task    *Task
		req     UpdateTaskRequest
		want    func(*Task)
		wantErr error
	}{
		{
			name: "update all fields",
			task: originalTask,
			req: UpdateTaskRequest{
				Title:       stringPtr("New Title"),
				Description: stringPtr("New Description"),
				Status:      stringPtr("completed"),
				Priority:    stringPtr("high"),
				Tags:        []string{"new"},
			},
			want: func(updated *Task) {
				assert.Equal(t, "test-id", updated.ID)
				assert.Equal(t, "New Title", updated.Title)
				assert.Equal(t, "New Description", updated.Description)
				assert.Equal(t, StatusCompleted, updated.Status)
				assert.Equal(t, PriorityHigh, updated.Priority)
				assert.Equal(t, []string{"new"}, updated.Tags)
				assert.Equal(t, now, updated.CreatedAt)
				assert.Equal(t, updateTime, updated.UpdatedAt)
			},
			wantErr: nil,
		},
		{
			name: "update single field",
			task: originalTask,
			req: UpdateTaskRequest{
				Status: stringPtr("in_progress"),
			},
			want: func(updated *Task) {
				assert.Equal(t, "Original Title", updated.Title)
				assert.Equal(t, "Original Description", updated.Description)
				assert.Equal(t, StatusInProgress, updated.Status)
				assert.Equal(t, PriorityLow, updated.Priority)
				assert.Equal(t, []string{"old"}, updated.Tags)
				assert.Equal(t, updateTime, updated.UpdatedAt)
			},
			wantErr: nil,
		},
		{
			name: "empty title",
			task: originalTask,
			req: UpdateTaskRequest{
				Title: stringPtr(""),
			},
			want:    nil,
			wantErr: ErrInvalidTitle,
		},
		{
			name: "invalid status",
			task: originalTask,
			req: UpdateTaskRequest{
				Status: stringPtr("invalid"),
			},
			want:    nil,
			wantErr: ErrInvalidStatus,
		},
		{
			name: "invalid priority",
			task: originalTask,
			req: UpdateTaskRequest{
				Priority: stringPtr("critical"),
			},
			want:    nil,
			wantErr: ErrInvalidPriority,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated, err := ApplyUpdate(tt.task, tt.req, updateTime)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, updated)
			} else {
				require.NoError(t, err)
				require.NotNil(t, updated)
				tt.want(updated)
			}
		})
	}
}

func TestServiceCreate(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateTaskRequest
		wantErr error
	}{
		{
			name: "successful creation",
			req: CreateTaskRequest{
				Title:    "Test task",
				Priority: "high",
			},
			wantErr: nil,
		},
		{
			name: "empty title",
			req: CreateTaskRequest{
				Title: "",
			},
			wantErr: ErrInvalidTitle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewInMemoryRepository()
			service := NewService(repo)

			task, err := service.Create(tt.req)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, task)
			} else {
				require.NoError(t, err)
				require.NotNil(t, task)
				assert.NotEmpty(t, task.ID)
				assert.Equal(t, tt.req.Title, task.Title)
			}
		})
	}
}

func TestServiceGetByID(t *testing.T) {
	repo := NewInMemoryRepository()
	service := NewService(repo)

	// Create a task
	created, err := service.Create(CreateTaskRequest{Title: "Test task"})
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		wantErr error
	}{
		{
			name:    "existing task",
			id:      created.ID,
			wantErr: nil,
		},
		{
			name:    "non-existent task",
			id:      "non-existent",
			wantErr: ErrTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := service.GetByID(tt.id)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, task)
			} else {
				require.NoError(t, err)
				require.NotNil(t, task)
				assert.Equal(t, tt.id, task.ID)
			}
		})
	}
}

func TestServiceGetAll(t *testing.T) {
	repo := NewInMemoryRepository()
	service := NewService(repo)

	// Create test tasks
	_, err := service.Create(CreateTaskRequest{Title: "Task 1", Priority: "high", Status: "pending"})
	require.NoError(t, err)
	_, err = service.Create(CreateTaskRequest{Title: "Task 2", Priority: "low", Status: "completed"})
	require.NoError(t, err)

	tests := []struct {
		name      string
		filter    FilterOptions
		wantCount int
	}{
		{
			name:      "all tasks",
			filter:    FilterOptions{},
			wantCount: 2,
		},
		{
			name: "filter by status",
			filter: FilterOptions{
				Status: statusPtr(StatusPending),
			},
			wantCount: 1,
		},
		{
			name: "filter by priority",
			filter: FilterOptions{
				Priority: priorityPtr(PriorityHigh),
			},
			wantCount: 1,
		},
		{
			name: "pagination",
			filter: FilterOptions{
				Limit:  1,
				Offset: 0,
			},
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, err := service.GetAll(tt.filter)
			require.NoError(t, err)
			assert.Len(t, tasks, tt.wantCount)
		})
	}
}

func TestServiceUpdate(t *testing.T) {
	repo := NewInMemoryRepository()
	service := NewService(repo)

	// Create a task
	created, err := service.Create(CreateTaskRequest{Title: "Original"})
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		req     UpdateTaskRequest
		wantErr error
	}{
		{
			name: "successful update",
			id:   created.ID,
			req: UpdateTaskRequest{
				Title: stringPtr("Updated"),
			},
			wantErr: nil,
		},
		{
			name: "non-existent task",
			id:   "non-existent",
			req: UpdateTaskRequest{
				Title: stringPtr("Updated"),
			},
			wantErr: ErrTaskNotFound,
		},
		{
			name: "invalid update",
			id:   created.ID,
			req: UpdateTaskRequest{
				Title: stringPtr(""),
			},
			wantErr: ErrInvalidTitle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := service.Update(tt.id, tt.req)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, task)
			} else {
				require.NoError(t, err)
				require.NotNil(t, task)
				if tt.req.Title != nil {
					assert.Equal(t, *tt.req.Title, task.Title)
				}
			}
		})
	}
}

func TestServiceDelete(t *testing.T) {
	repo := NewInMemoryRepository()
	service := NewService(repo)

	// Create a task
	created, err := service.Create(CreateTaskRequest{Title: "To delete"})
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		wantErr error
	}{
		{
			name:    "successful deletion",
			id:      created.ID,
			wantErr: nil,
		},
		{
			name:    "non-existent task",
			id:      "non-existent",
			wantErr: ErrTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Delete(tt.id)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)

				// Verify task is deleted
				_, err := service.GetByID(tt.id)
				assert.ErrorIs(t, err, ErrTaskNotFound)
			}
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func statusPtr(s Status) *Status {
	return &s
}

func priorityPtr(p Priority) *Priority {
	return &p
}
