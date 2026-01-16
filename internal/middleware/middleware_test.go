package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		checkHeaders   bool
	}{
		{
			name:           "sets CORS headers for GET request",
			method:         "GET",
			expectedStatus: http.StatusOK,
			checkHeaders:   true,
		},
		{
			name:           "sets CORS headers for POST request",
			method:         "POST",
			expectedStatus: http.StatusOK,
			checkHeaders:   true,
		},
		{
			name:           "handles OPTIONS preflight request",
			method:         "OPTIONS",
			expectedStatus: http.StatusNoContent,
			checkHeaders:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router := gin.New()
			router.Use(CORS())
			router.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})
			router.POST("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			// Create request
			req, _ := http.NewRequest(tt.method, "/test", nil)
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkHeaders {
				assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
				assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
				assert.Equal(t, "Origin, Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
			}
		})
	}
}

func TestRequestTimer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string
	}{
		{
			name: "adds response time header",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router := gin.New()
			router.Use(RequestTimer())
			router.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			// Create request
			req, _ := http.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, http.StatusOK, w.Code)
			assert.NotEmpty(t, w.Header().Get("X-Response-Time"))
		})
	}
}
