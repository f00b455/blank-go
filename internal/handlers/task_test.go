package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/f00b455/blank-go/pkg/task"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockTaskRepository implements task.Repository for testing
type MockTaskRepository struct {
	CreateFunc     func(t *task.Task) error
	GetByIDFunc    func(id string) (*task.Task, error)
	GetAllFunc     func(filter task.FilterOptions) ([]*task.Task, error)
	UpdateFunc     func(t *task.Task) error
	DeleteFunc     func(id string) error
}

func (m *MockTaskRepository) Create(t *task.Task) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(t)
	}
	return nil
}

func (m *MockTaskRepository) GetByID(id string) (*task.Task, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, errors.New("not implemented")
}

func (m *MockTaskRepository) GetAll(filter task.FilterOptions) ([]*task.Task, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(filter)
	}
	return nil, nil
}

func (m *MockTaskRepository) Update(t *task.Task) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(t)
	}
	return nil
}

func (m *MockTaskRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func TestNewTaskHandler(t *testing.T) {
	mockRepo := &MockTaskRepository{}
	service := task.NewService(mockRepo)
	handler := NewTaskHandler(service)

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.service)
}

func TestCreateTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockCreateFunc func(tk *task.Task) error
		expectedStatus int
		expectError    bool
	}{
		{
			name: "successful task creation",
			requestBody: map[string]interface{}{
				"title":    "Test Task",
				"priority": "high",
			},
			mockCreateFunc: func(tk *task.Task) error {
				return nil
			},
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name: "invalid JSON body",
			requestBody: "invalid json",
			mockCreateFunc: func(tk *task.Task) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "service error",
			requestBody: map[string]interface{}{
				"title": "Test Task",
			},
			mockCreateFunc: func(tk *task.Task) error {
				return task.ErrInvalidTitle
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockTaskRepository{
				CreateFunc: tt.mockCreateFunc,
			}
			service := task.NewService(mockRepo)
			handler := NewTaskHandler(service)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			handler.CreateTask(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestGetTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name            string
		taskID          string
		mockGetByIDFunc func(id string) (*task.Task, error)
		expectedStatus  int
		expectError     bool
	}{
		{
			name:   "successful task retrieval",
			taskID: "test-id",
			mockGetByIDFunc: func(id string) (*task.Task, error) {
				return &task.Task{
					ID:    "test-id",
					Title: "Test Task",
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:   "task not found",
			taskID: "nonexistent",
			mockGetByIDFunc: func(id string) (*task.Task, error) {
				return nil, task.ErrTaskNotFound
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockTaskRepository{
				GetByIDFunc: tt.mockGetByIDFunc,
			}
			service := task.NewService(mockRepo)
			handler := NewTaskHandler(service)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: tt.taskID}}
			req, _ := http.NewRequest("GET", "/api/v1/tasks/"+tt.taskID, nil)
			c.Request = req

			handler.GetTask(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestListTasks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name            string
		query           string
		mockGetAllFunc  func(filter task.FilterOptions) ([]*task.Task, error)
		expectedStatus  int
		expectedCount   int
	}{
		{
			name:  "list all tasks",
			query: "",
			mockGetAllFunc: func(filter task.FilterOptions) ([]*task.Task, error) {
				return []*task.Task{
					{ID: "1", Title: "Task 1"},
					{ID: "2", Title: "Task 2"},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:  "list tasks with filters",
			query: "?status=pending&priority=high",
			mockGetAllFunc: func(filter task.FilterOptions) ([]*task.Task, error) {
				return []*task.Task{
					{ID: "1", Title: "Task 1", Status: "pending", Priority: "high"},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockTaskRepository{
				GetAllFunc: tt.mockGetAllFunc,
			}
			service := task.NewService(mockRepo)
			handler := NewTaskHandler(service)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req, _ := http.NewRequest("GET", "/api/v1/tasks"+tt.query, nil)
			c.Request = req

			handler.ListTasks(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response []task.Task
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Len(t, response, tt.expectedCount)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		taskID         string
		requestBody    interface{}
		mockUpdateFunc func(tk *task.Task) error
		expectedStatus int
	}{
		{
			name:   "successful task update",
			taskID: "test-id",
			requestBody: map[string]interface{}{
				"title":  "Updated Task",
				"status": "completed",
			},
			mockUpdateFunc: func(tk *task.Task) error {
				return nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "task not found",
			taskID: "nonexistent",
			requestBody: map[string]interface{}{
				"title": "Updated Task",
			},
			mockUpdateFunc: func(tk *task.Task) error {
				return task.ErrTaskNotFound
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockTaskRepository{
				UpdateFunc: tt.mockUpdateFunc,
				GetByIDFunc: func(id string) (*task.Task, error) {
					if id == "test-id" {
						return &task.Task{ID: id, Title: "Old Title"}, nil
					}
					return nil, task.ErrTaskNotFound
				},
			}
			service := task.NewService(mockRepo)
			handler := NewTaskHandler(service)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: tt.taskID}}

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("PUT", "/api/v1/tasks/"+tt.taskID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			handler.UpdateTask(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestDeleteTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		taskID         string
		mockDeleteFunc func(id string) error
		expectedStatus int
	}{
		{
			name:   "successful task deletion",
			taskID: "test-id",
			mockDeleteFunc: func(id string) error {
				return nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "task not found",
			taskID: "nonexistent",
			mockDeleteFunc: func(id string) error {
				return task.ErrTaskNotFound
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockTaskRepository{
				DeleteFunc: tt.mockDeleteFunc,
			}
			service := task.NewService(mockRepo)
			handler := NewTaskHandler(service)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: tt.taskID}}
			req, _ := http.NewRequest("DELETE", "/api/v1/tasks/"+tt.taskID, nil)
			c.Request = req

			handler.DeleteTask(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestHandleServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{
			name:           "invalid title error",
			err:            task.ErrInvalidTitle,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid status error",
			err:            task.ErrInvalidStatus,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid priority error",
			err:            task.ErrInvalidPriority,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "not found error",
			err:            task.ErrTaskNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "generic error",
			err:            errors.New("generic error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			handleServiceError(c, tt.err)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestParseFilterOptions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name              string
		queryString       string
		expectStatusSet   bool
		expectedStatusVal string
		expectedLimit     int
		expectedOffset    int
		expectPrioritySet bool
		expectedPriority  string
		expectTagSet      bool
		expectedTag       string
		expectedSortBy    string
		expectedSortDesc  bool
	}{
		{
			name:              "parse status filter",
			queryString:       "?status=pending",
			expectStatusSet:   true,
			expectedStatusVal: "pending",
			expectedLimit:     0,
			expectedSortBy:    "created_at",
		},
		{
			name:            "parse limit",
			queryString:     "?limit=50",
			expectStatusSet: false,
			expectedLimit:   50,
			expectedSortBy:  "created_at",
		},
		{
			name:            "no filters",
			queryString:     "",
			expectStatusSet: false,
			expectedLimit:   0,
			expectedSortBy:  "created_at",
		},
		{
			name:           "parse offset",
			queryString:    "?offset=10",
			expectedOffset: 10,
			expectedSortBy: "created_at",
		},
		{
			name:              "parse priority filter",
			queryString:       "?priority=high",
			expectPrioritySet: true,
			expectedPriority:  "high",
			expectedSortBy:    "created_at",
		},
		{
			name:           "parse tag filter",
			queryString:    "?tag=important",
			expectTagSet:   true,
			expectedTag:    "important",
			expectedSortBy: "created_at",
		},
		{
			name:           "parse sort_by",
			queryString:    "?sort_by=title",
			expectedSortBy: "title",
		},
		{
			name:             "parse order desc",
			queryString:      "?order=desc",
			expectedSortBy:   "created_at",
			expectedSortDesc: true,
		},
		{
			name:            "invalid limit ignored",
			queryString:     "?limit=invalid",
			expectedLimit:   0,
			expectedSortBy:  "created_at",
		},
		{
			name:           "negative limit ignored",
			queryString:    "?limit=-5",
			expectedLimit:  0,
			expectedSortBy: "created_at",
		},
		{
			name:           "invalid offset ignored",
			queryString:    "?offset=invalid",
			expectedOffset: 0,
			expectedSortBy: "created_at",
		},
		{
			name:           "negative offset ignored",
			queryString:    "?offset=-5",
			expectedOffset: 0,
			expectedSortBy: "created_at",
		},
		{
			name:            "invalid status ignored",
			queryString:     "?status=invalid_status",
			expectStatusSet: false,
			expectedSortBy:  "created_at",
		},
		{
			name:              "invalid priority ignored",
			queryString:       "?priority=invalid_priority",
			expectPrioritySet: false,
			expectedSortBy:    "created_at",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req, _ := http.NewRequest("GET", "/api/v1/tasks"+tt.queryString, nil)
			c.Request = req

			filter := parseFilterOptions(c)

			if tt.expectStatusSet {
				require.NotNil(t, filter.Status)
				assert.Equal(t, tt.expectedStatusVal, string(*filter.Status))
			} else {
				assert.Nil(t, filter.Status)
			}
			assert.Equal(t, tt.expectedLimit, filter.Limit)
			assert.Equal(t, tt.expectedOffset, filter.Offset)
			assert.Equal(t, tt.expectedSortBy, filter.SortBy)
			assert.Equal(t, tt.expectedSortDesc, filter.SortDesc)

			if tt.expectPrioritySet {
				require.NotNil(t, filter.Priority)
				assert.Equal(t, tt.expectedPriority, string(*filter.Priority))
			} else {
				assert.Nil(t, filter.Priority)
			}

			if tt.expectTagSet {
				require.NotNil(t, filter.Tag)
				assert.Equal(t, tt.expectedTag, *filter.Tag)
			} else {
				assert.Nil(t, filter.Tag)
			}
		})
	}
}

func TestListTasks_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := &MockTaskRepository{
		GetAllFunc: func(filter task.FilterOptions) ([]*task.Task, error) {
			return nil, errors.New("database error")
		},
	}
	service := task.NewService(mockRepo)
	handler := NewTaskHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/api/v1/tasks", nil)
	c.Request = req

	handler.ListTasks(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "INTERNAL_ERROR")
}

func TestUpdateTask_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := &MockTaskRepository{}
	service := task.NewService(mockRepo)
	handler := NewTaskHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "test-id"}}

	req, _ := http.NewRequest("PUT", "/api/v1/tasks/test-id", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	handler.UpdateTask(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "INVALID_REQUEST")
}
