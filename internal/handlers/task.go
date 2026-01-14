package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/f00b455/blank-go/pkg/task"
	"github.com/gin-gonic/gin"
)

// TaskHandler handles HTTP requests for task operations
type TaskHandler struct {
	service *task.Service
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(service *task.Service) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains error information
type ErrorDetail struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

// CreateTask handles POST /api/v1/tasks
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req task.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request body",
				Details: []string{err.Error()},
			},
		})
		return
	}

	createdTask, err := h.service.Create(req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, createdTask)
}

// GetTask handles GET /api/v1/tasks/:id
func (h *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")

	foundTask, err := h.service.GetByID(id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, foundTask)
}

// ListTasks handles GET /api/v1/tasks
func (h *TaskHandler) ListTasks(c *gin.Context) {
	filter := parseFilterOptions(c)

	tasks, err := h.service.GetAll(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to retrieve tasks",
			},
		})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// UpdateTask handles PUT /api/v1/tasks/:id
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")

	var req task.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request body",
				Details: []string{err.Error()},
			},
		})
		return
	}

	updatedTask, err := h.service.Update(id, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}

// DeleteTask handles DELETE /api/v1/tasks/:id
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")

	err := h.service.Delete(id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task deleted successfully",
	})
}

// handleServiceError converts service errors to HTTP responses
func handleServiceError(c *gin.Context, err error) {
	if errors.Is(err, task.ErrTaskNotFound) {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: ErrorDetail{
				Code:    "NOT_FOUND",
				Message: "task not found",
			},
		})
		return
	}

	if errors.Is(err, task.ErrInvalidTitle) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "VALIDATION_ERROR",
				Message: "title is required",
			},
		})
		return
	}

	if errors.Is(err, task.ErrInvalidStatus) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "VALIDATION_ERROR",
				Message: "invalid status",
			},
		})
		return
	}

	if errors.Is(err, task.ErrInvalidPriority) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "VALIDATION_ERROR",
				Message: "invalid priority",
			},
		})
		return
	}

	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: ErrorDetail{
			Code:    "INTERNAL_ERROR",
			Message: "An unexpected error occurred",
		},
	})
}

// parseFilterOptions parses query parameters into FilterOptions
func parseFilterOptions(c *gin.Context) task.FilterOptions {
	filter := task.FilterOptions{
		Limit:    0,
		Offset:   0,
		SortBy:   "created_at",
		SortDesc: false,
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	if status := c.Query("status"); status != "" && task.IsValidStatus(status) {
		s := task.Status(status)
		filter.Status = &s
	}

	if priority := c.Query("priority"); priority != "" && task.IsValidPriority(priority) {
		p := task.Priority(priority)
		filter.Priority = &p
	}

	if tag := c.Query("tag"); tag != "" {
		filter.Tag = &tag
	}

	if sortBy := c.Query("sort_by"); sortBy != "" {
		filter.SortBy = sortBy
	}

	if order := c.Query("order"); order == "desc" {
		filter.SortDesc = true
	}

	return filter
}
