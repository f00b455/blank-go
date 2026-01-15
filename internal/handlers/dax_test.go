package handlers

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/f00b455/blank-go/mocks"
	"github.com/f00b455/blank-go/pkg/dax"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewDAXHandler(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	service := dax.NewService(mockRepo)
	handler := NewDAXHandler(service)

	assert.NotNil(t, handler)
	assert.Equal(t, service, handler.service)
}

func TestDAXHandler_ImportCSV_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockRepository)
	service := dax.NewService(mockRepo)
	handler := NewDAXHandler(service)

	// Setup mock expectations
	mockRepo.On("BulkUpsert", mock.AnythingOfType("[]dax.DAXRecord")).Return(nil)

	// Create CSV content
	csvContent := `company,ticker,report_type,metric,year,value,currency
Siemens AG,SIE,income,EBITDA,2025,15859000000.0,EUR
SAP SE,SAP,income,Net Income,2025,8500000000.0,EUR`

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.csv")
	_, _ = part.Write([]byte(csvContent))
	_ = writer.Close()

	// Create test request
	req := httptest.NewRequest(http.MethodPost, "/api/v1/dax/import", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call handler
	handler.ImportCSV(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "records_imported")
	mockRepo.AssertExpectations(t)
}

func TestDAXHandler_ImportCSV_NoFile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockRepository)
	service := dax.NewService(mockRepo)
	handler := NewDAXHandler(service)

	// Create test request without file
	req := httptest.NewRequest(http.MethodPost, "/api/v1/dax/import", nil)
	w := httptest.NewRecorder()

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call handler
	handler.ImportCSV(c)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "file is required")
}

func TestDAXHandler_ImportCSV_InvalidCSV(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockRepository)
	service := dax.NewService(mockRepo)
	handler := NewDAXHandler(service)

	// Create invalid CSV content (missing required fields)
	csvContent := `company,ticker
Siemens AG,SIE`

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.csv")
	_, _ = part.Write([]byte(csvContent))
	_ = writer.Close()

	// Create test request
	req := httptest.NewRequest(http.MethodPost, "/api/v1/dax/import", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call handler
	handler.ImportCSV(c)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "missing required fields")
}

func TestDAXHandler_ImportCSV_RepositoryError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockRepository)
	service := dax.NewService(mockRepo)
	handler := NewDAXHandler(service)

	// Setup mock to return error
	mockRepo.On("BulkUpsert", mock.AnythingOfType("[]dax.DAXRecord")).
		Return(errors.New("database error"))

	// Create valid CSV content
	csvContent := `company,ticker,report_type,metric,year,value,currency
Siemens AG,SIE,income,EBITDA,2025,15859000000.0,EUR`

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.csv")
	_, _ = part.Write([]byte(csvContent))
	_ = writer.Close()

	// Create test request
	req := httptest.NewRequest(http.MethodPost, "/api/v1/dax/import", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call handler
	handler.ImportCSV(c)

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "database error")
	mockRepo.AssertExpectations(t)
}

func TestDAXHandler_GetAll_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockRepository)
	service := dax.NewService(mockRepo)
	handler := NewDAXHandler(service)

	// Setup mock expectations
	expectedRecords := []dax.DAXRecord{
		{Ticker: "SIE", Year: 2025},
		{Ticker: "SAP", Year: 2024},
	}
	mockRepo.On("FindAll", 1, 10).Return(expectedRecords, 2, nil)

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/dax?page=1&limit=10", nil)
	w := httptest.NewRecorder()

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call handler
	handler.GetAll(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "SIE")
	assert.Contains(t, w.Body.String(), "SAP")
	mockRepo.AssertExpectations(t)
}

func TestDAXHandler_GetAll_RepositoryError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockRepository)
	service := dax.NewService(mockRepo)
	handler := NewDAXHandler(service)

	// Setup mock to return error
	mockRepo.On("FindAll", 1, 10).Return([]dax.DAXRecord{}, 0, errors.New("database error"))

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/dax", nil)
	w := httptest.NewRecorder()

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call handler
	handler.GetAll(c)

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "database error")
	mockRepo.AssertExpectations(t)
}

func TestDAXHandler_GetByFilters_WithTickerAndYear(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockRepository)
	service := dax.NewService(mockRepo)
	handler := NewDAXHandler(service)

	// Setup mock expectations
	year := 2025
	expectedRecords := []dax.DAXRecord{
		{Ticker: "SIE", Year: 2025},
	}
	mockRepo.On("FindByFilters", "SIE", &year, 1, 10).Return(expectedRecords, 1, nil)

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/dax/filter?ticker=SIE&year=2025&page=1&limit=10", nil)
	w := httptest.NewRecorder()

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call handler
	handler.GetByFilters(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "SIE")
	mockRepo.AssertExpectations(t)
}

func TestDAXHandler_GetByFilters_InvalidYear(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockRepository)
	service := dax.NewService(mockRepo)
	handler := NewDAXHandler(service)

	// Create test request with invalid year
	req := httptest.NewRequest(http.MethodGet, "/api/v1/dax/filter?year=invalid", nil)
	w := httptest.NewRecorder()

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call handler
	handler.GetByFilters(c)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid year parameter")
}

func TestDAXHandler_GetMetrics_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockRepository)
	service := dax.NewService(mockRepo)
	handler := NewDAXHandler(service)

	// Setup mock expectations
	expectedMetrics := []string{"EBITDA", "Net Income"}
	mockRepo.On("GetMetrics", "SIE").Return(expectedMetrics, nil)

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/dax/metrics?ticker=SIE", nil)
	w := httptest.NewRecorder()

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call handler
	handler.GetMetrics(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "EBITDA")
	assert.Contains(t, w.Body.String(), "Net Income")
	mockRepo.AssertExpectations(t)
}

func TestDAXHandler_GetMetrics_NoTicker(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockRepository)
	service := dax.NewService(mockRepo)
	handler := NewDAXHandler(service)

	// Create test request without ticker
	req := httptest.NewRequest(http.MethodGet, "/api/v1/dax/metrics", nil)
	w := httptest.NewRecorder()

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call handler
	handler.GetMetrics(c)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "ticker parameter is required")
}

func TestParseIntQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		queryValue   string
		defaultValue int
		expected     int
	}{
		{
			name:         "Valid positive integer",
			queryValue:   "5",
			defaultValue: 10,
			expected:     5,
		},
		{
			name:         "Empty string returns default",
			queryValue:   "",
			defaultValue: 10,
			expected:     10,
		},
		{
			name:         "Invalid integer returns default",
			queryValue:   "abc",
			defaultValue: 10,
			expected:     10,
		},
		{
			name:         "Zero returns default",
			queryValue:   "0",
			defaultValue: 10,
			expected:     10,
		},
		{
			name:         "Negative returns default",
			queryValue:   "-5",
			defaultValue: 10,
			expected:     10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/?key="+tt.queryValue, nil)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			result := parseIntQuery(c, "key", tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}
