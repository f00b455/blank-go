package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "returns healthy status",
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "healthy",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req, _ := http.NewRequest("GET", "/health", nil)
			c.Request = req

			// Execute
			HealthCheck(c)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedBody["status"], response["status"])
		})
	}
}

func TestPing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "returns pong message",
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "pong",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req, _ := http.NewRequest("GET", "/ping", nil)
			c.Request = req

			// Execute
			Ping(c)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedBody["message"], response["message"])
		})
	}
}

func TestDetailedHealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		startTime      time.Time
		expectedStatus int
	}{
		{
			name:           "returns detailed health information",
			startTime:      time.Now().Add(-5 * time.Second),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "returns detailed health with different start time",
			startTime:      time.Now().Add(-1 * time.Minute),
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req, _ := http.NewRequest("GET", "/health/detailed", nil)
			c.Request = req

			// Create handler with start time
			handler := DetailedHealthCheck(tt.startTime)

			// Execute
			handler(c)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Check required fields
			assert.Equal(t, "healthy", response["status"])
			assert.NotEmpty(t, response["timestamp"])
			assert.NotEmpty(t, response["version"])
			assert.NotZero(t, response["uptime_seconds"])

			// Check system info
			system, ok := response["system"].(map[string]interface{})
			require.True(t, ok, "system should be present")
			assert.NotEmpty(t, system["go_version"])
			assert.NotZero(t, system["goroutines"])

			// Check checks info
			checks, ok := response["checks"].(map[string]interface{})
			require.True(t, ok, "checks should be present")
			assert.Equal(t, "ok", checks["api"])
		})
	}
}
