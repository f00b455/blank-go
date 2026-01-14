package task

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryRepository_Create(t *testing.T) {
	repo := NewInMemoryRepository()
	task := &Task{
		ID:        "test-id",
		Title:     "Test Task",
		Status:    StatusPending,
		Priority:  PriorityMedium,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(task)
	require.NoError(t, err)

	// Verify task was created
	retrieved, err := repo.GetByID(task.ID)
	require.NoError(t, err)
	assert.Equal(t, task.ID, retrieved.ID)
	assert.Equal(t, task.Title, retrieved.Title)
}

func TestInMemoryRepository_GetByID(t *testing.T) {
	repo := NewInMemoryRepository()
	task := &Task{
		ID:        "test-id",
		Title:     "Test Task",
		Status:    StatusPending,
		Priority:  PriorityMedium,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create task
	err := repo.Create(task)
	require.NoError(t, err)

	t.Run("existing task", func(t *testing.T) {
		retrieved, err := repo.GetByID("test-id")
		require.NoError(t, err)
		assert.Equal(t, task.ID, retrieved.ID)
	})

	t.Run("non-existent task", func(t *testing.T) {
		retrieved, err := repo.GetByID("non-existent")
		assert.ErrorIs(t, err, ErrTaskNotFound)
		assert.Nil(t, retrieved)
	})
}

func TestInMemoryRepository_Update(t *testing.T) {
	repo := NewInMemoryRepository()
	task := &Task{
		ID:        "test-id",
		Title:     "Original",
		Status:    StatusPending,
		Priority:  PriorityMedium,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create task
	err := repo.Create(task)
	require.NoError(t, err)

	t.Run("update existing task", func(t *testing.T) {
		task.Title = "Updated"
		err := repo.Update(task)
		require.NoError(t, err)

		retrieved, err := repo.GetByID(task.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated", retrieved.Title)
	})

	t.Run("update non-existent task", func(t *testing.T) {
		nonExistent := &Task{ID: "non-existent", Title: "Test"}
		err := repo.Update(nonExistent)
		assert.ErrorIs(t, err, ErrTaskNotFound)
	})
}

func TestInMemoryRepository_Delete(t *testing.T) {
	repo := NewInMemoryRepository()
	task := &Task{
		ID:        "test-id",
		Title:     "Test Task",
		Status:    StatusPending,
		Priority:  PriorityMedium,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create task
	err := repo.Create(task)
	require.NoError(t, err)

	t.Run("delete existing task", func(t *testing.T) {
		err := repo.Delete(task.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = repo.GetByID(task.ID)
		assert.ErrorIs(t, err, ErrTaskNotFound)
	})

	t.Run("delete non-existent task", func(t *testing.T) {
		err := repo.Delete("non-existent")
		assert.ErrorIs(t, err, ErrTaskNotFound)
	})
}

func TestInMemoryRepository_GetAll(t *testing.T) {
	repo := NewInMemoryRepository()

	// Create test tasks
	tasks := []*Task{
		{
			ID:        "1",
			Title:     "Task 1",
			Status:    StatusPending,
			Priority:  PriorityHigh,
			Tags:      []string{"work", "urgent"},
			CreatedAt: time.Date(2026, 1, 10, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2026, 1, 10, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        "2",
			Title:     "Task 2",
			Status:    StatusInProgress,
			Priority:  PriorityMedium,
			Tags:      []string{"personal"},
			CreatedAt: time.Date(2026, 1, 11, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2026, 1, 11, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        "3",
			Title:     "Task 3",
			Status:    StatusCompleted,
			Priority:  PriorityLow,
			Tags:      []string{"work"},
			CreatedAt: time.Date(2026, 1, 12, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2026, 1, 12, 10, 0, 0, 0, time.UTC),
		},
	}

	for _, task := range tasks {
		err := repo.Create(task)
		require.NoError(t, err)
	}

	t.Run("get all tasks", func(t *testing.T) {
		result, err := repo.GetAll(FilterOptions{})
		require.NoError(t, err)
		assert.Len(t, result, 3)
	})

	t.Run("filter by status", func(t *testing.T) {
		status := StatusPending
		result, err := repo.GetAll(FilterOptions{Status: &status})
		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, StatusPending, result[0].Status)
	})

	t.Run("filter by priority", func(t *testing.T) {
		priority := PriorityHigh
		result, err := repo.GetAll(FilterOptions{Priority: &priority})
		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, PriorityHigh, result[0].Priority)
	})

	t.Run("filter by tag", func(t *testing.T) {
		tag := "work"
		result, err := repo.GetAll(FilterOptions{Tag: &tag})
		require.NoError(t, err)
		assert.Len(t, result, 2)
	})

	t.Run("pagination", func(t *testing.T) {
		result, err := repo.GetAll(FilterOptions{Limit: 2, Offset: 0})
		require.NoError(t, err)
		assert.Len(t, result, 2)

		result, err = repo.GetAll(FilterOptions{Limit: 2, Offset: 2})
		require.NoError(t, err)
		assert.Len(t, result, 1)
	})

	t.Run("sort by created_at asc", func(t *testing.T) {
		result, err := repo.GetAll(FilterOptions{SortBy: "created_at", SortDesc: false})
		require.NoError(t, err)
		assert.Equal(t, "1", result[0].ID)
		assert.Equal(t, "2", result[1].ID)
		assert.Equal(t, "3", result[2].ID)
	})

	t.Run("sort by created_at desc", func(t *testing.T) {
		result, err := repo.GetAll(FilterOptions{SortBy: "created_at", SortDesc: true})
		require.NoError(t, err)
		assert.Equal(t, "3", result[0].ID)
		assert.Equal(t, "2", result[1].ID)
		assert.Equal(t, "1", result[2].ID)
	})

	t.Run("sort by priority", func(t *testing.T) {
		result, err := repo.GetAll(FilterOptions{SortBy: "priority", SortDesc: true})
		require.NoError(t, err)
		assert.Equal(t, PriorityHigh, result[0].Priority)
	})
}

func TestMatchesFilter(t *testing.T) {
	task := &Task{
		Status:   StatusPending,
		Priority: PriorityHigh,
		Tags:     []string{"work", "urgent"},
	}

	tests := []struct {
		name   string
		filter FilterOptions
		want   bool
	}{
		{
			name:   "no filter",
			filter: FilterOptions{},
			want:   true,
		},
		{
			name: "matching status",
			filter: FilterOptions{
				Status: statusPtr(StatusPending),
			},
			want: true,
		},
		{
			name: "non-matching status",
			filter: FilterOptions{
				Status: statusPtr(StatusCompleted),
			},
			want: false,
		},
		{
			name: "matching priority",
			filter: FilterOptions{
				Priority: priorityPtr(PriorityHigh),
			},
			want: true,
		},
		{
			name: "matching tag",
			filter: FilterOptions{
				Tag: stringPtr("work"),
			},
			want: true,
		},
		{
			name: "non-matching tag",
			filter: FilterOptions{
				Tag: stringPtr("personal"),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchesFilter(task, tt.filter)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestPriorityValue(t *testing.T) {
	tests := []struct {
		priority Priority
		want     int
	}{
		{PriorityLow, 1},
		{PriorityMedium, 2},
		{PriorityHigh, 3},
		{"invalid", 0},
	}

	for _, tt := range tests {
		t.Run(string(tt.priority), func(t *testing.T) {
			result := priorityValue(tt.priority)
			assert.Equal(t, tt.want, result)
		})
	}
}
