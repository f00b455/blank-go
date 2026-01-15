package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call handler
	HealthCheck(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "healthy")
	assert.Contains(t, w.Body.String(), `"status":"healthy"`)
}

func TestPing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call handler
	Ping(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
	assert.Contains(t, w.Body.String(), `"message":"pong"`)
}

func TestDetailedHealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create start time
	startTime := time.Now().Add(-5 * time.Minute)

	// Get handler function
	handler := DetailedHealthCheck(startTime)

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/health/detailed", nil)
	w := httptest.NewRecorder()

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call handler
	handler(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "healthy")
	assert.Contains(t, w.Body.String(), "timestamp")
	assert.Contains(t, w.Body.String(), "version")
	assert.Contains(t, w.Body.String(), "uptime_seconds")
	assert.Contains(t, w.Body.String(), "system")
	assert.Contains(t, w.Body.String(), "go_version")
	assert.Contains(t, w.Body.String(), "goroutines")
	assert.Contains(t, w.Body.String(), "memory_alloc_mb")
	assert.Contains(t, w.Body.String(), "memory_sys_mb")
	assert.Contains(t, w.Body.String(), "gc_runs")
	assert.Contains(t, w.Body.String(), "checks")
	assert.Contains(t, w.Body.String(), `"api":"ok"`)
}

func TestDetailedHealthCheck_UptimeCalculation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create start time 10 seconds ago
	startTime := time.Now().Add(-10 * time.Second)

	// Get handler function
	handler := DetailedHealthCheck(startTime)

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/health/detailed", nil)
	w := httptest.NewRecorder()

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call handler
	handler(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	// Uptime should be at least 10 seconds (allowing for test execution time)
	assert.Contains(t, w.Body.String(), "uptime_seconds")
}

func TestDetailedHealthCheck_MemoryMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create start time
	startTime := time.Now()

	// Get handler function
	handler := DetailedHealthCheck(startTime)

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/health/detailed", nil)
	w := httptest.NewRecorder()

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call handler
	handler(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	// Verify memory metrics are present
	body := w.Body.String()
	assert.Contains(t, body, "memory_alloc_mb")
	assert.Contains(t, body, "memory_sys_mb")
	assert.Contains(t, body, "gc_runs")
	// Memory values should be numeric (greater than or equal to 0)
	assert.NotContains(t, body, `"memory_alloc_mb":null`)
	assert.NotContains(t, body, `"memory_sys_mb":null`)
}

func TestDetailedHealthCheck_SystemInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create start time
	startTime := time.Now()

	// Get handler function
	handler := DetailedHealthCheck(startTime)

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/health/detailed", nil)
	w := httptest.NewRecorder()

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call handler
	handler(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	// Verify Go version is present
	assert.Contains(t, body, "go_version")
	assert.Contains(t, body, "go")
	// Verify goroutines count is present
	assert.Contains(t, body, "goroutines")
}
