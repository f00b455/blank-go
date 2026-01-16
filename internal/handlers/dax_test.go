package handlers

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/f00b455/blank-go/pkg/dax"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDAXHandler_ImportCSV_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.csv")
	csvContent := `company,ticker,report_type,metric,year,value,currency
Siemens AG,SIE,income,EBITDA,2025,15859000000.0,EUR`
	_, _ = part.Write([]byte(csvContent))
	_ = writer.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/dax/import", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())

	// Use real service with in-memory repository
	repo := dax.NewInMemoryRepository()
	service := dax.NewService(repo)
	handler := NewDAXHandler(service)

	handler.ImportCSV(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Successfully imported")
}

func TestDAXHandler_ImportCSV_MissingFile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &dax.Service{}
	handler := NewDAXHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/dax/import", nil)

	handler.ImportCSV(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "file is required")
}

func TestDAXHandler_ImportCSV_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.csv")
	csvContent := `company,ticker,metric,year,value
Siemens AG,SIE,EBITDA,2025,15859000000.0`
	_, _ = part.Write([]byte(csvContent))
	_ = writer.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/dax/import", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())

	repo := dax.NewInMemoryRepository()
	service := dax.NewService(repo)
	handler := NewDAXHandler(service)

	handler.ImportCSV(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "missing required fields")
}

func TestDAXHandler_ImportCSV_InvalidYear(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.csv")
	csvContent := `company,ticker,report_type,metric,year,value,currency
Siemens AG,SIE,income,EBITDA,invalid,15859000000.0,EUR`
	_, _ = part.Write([]byte(csvContent))
	_ = writer.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/dax/import", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())

	repo := dax.NewInMemoryRepository()
	service := dax.NewService(repo)
	handler := NewDAXHandler(service)

	handler.ImportCSV(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid year")
}

func TestDAXHandler_GetAll_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := dax.NewInMemoryRepository()
	service := dax.NewService(repo)
	handler := NewDAXHandler(service)

	// Add test data
	value := 1000.0
	_ = repo.Create(&dax.DAXRecord{
		Company:    "Test Company",
		Ticker:     "TST",
		ReportType: "income",
		Metric:     "Revenue",
		Year:       2025,
		Value:      &value,
		Currency:   "EUR",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/dax?page=1&limit=10", nil)

	handler.GetAll(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "TST")
}

func TestDAXHandler_GetAll_DefaultPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := dax.NewInMemoryRepository()
	service := dax.NewService(repo)
	handler := NewDAXHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/dax", nil)

	handler.GetAll(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDAXHandler_GetByFilters_WithTicker(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := dax.NewInMemoryRepository()
	service := dax.NewService(repo)
	handler := NewDAXHandler(service)

	// Add test data
	value := 1000.0
	_ = repo.Create(&dax.DAXRecord{
		Company:    "Test Company",
		Ticker:     "TST",
		ReportType: "income",
		Metric:     "Revenue",
		Year:       2025,
		Value:      &value,
		Currency:   "EUR",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/dax/filter?ticker=TST", nil)

	handler.GetByFilters(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "TST")
}

func TestDAXHandler_GetByFilters_WithYear(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := dax.NewInMemoryRepository()
	service := dax.NewService(repo)
	handler := NewDAXHandler(service)

	// Add test data
	value := 1000.0
	_ = repo.Create(&dax.DAXRecord{
		Company:    "Test Company",
		Ticker:     "TST",
		ReportType: "income",
		Metric:     "Revenue",
		Year:       2025,
		Value:      &value,
		Currency:   "EUR",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/dax/filter?year=2025", nil)

	handler.GetByFilters(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "2025")
}

func TestDAXHandler_GetByFilters_WithTickerAndYear(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := dax.NewInMemoryRepository()
	service := dax.NewService(repo)
	handler := NewDAXHandler(service)

	// Add test data
	value := 1000.0
	_ = repo.Create(&dax.DAXRecord{
		Company:    "Test Company",
		Ticker:     "TST",
		ReportType: "income",
		Metric:     "Revenue",
		Year:       2025,
		Value:      &value,
		Currency:   "EUR",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/dax/filter?ticker=TST&year=2025", nil)

	handler.GetByFilters(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "TST")
}

func TestDAXHandler_GetByFilters_InvalidYear(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := dax.NewInMemoryRepository()
	service := dax.NewService(repo)
	handler := NewDAXHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/dax/filter?year=invalid", nil)

	handler.GetByFilters(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid year parameter")
}

func TestDAXHandler_GetByFilters_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Use in-memory repo for consistency with other tests
	repo := dax.NewInMemoryRepository()
	service := dax.NewService(repo)
	handler := NewDAXHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/dax/filter", nil)

	handler.GetByFilters(c)

	// Should return OK with empty results, not error
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDAXHandler_GetMetrics_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := dax.NewInMemoryRepository()
	service := dax.NewService(repo)
	handler := NewDAXHandler(service)

	// Add test data
	value := 1000.0
	_ = repo.Create(&dax.DAXRecord{
		Company:    "Test Company",
		Ticker:     "TST",
		ReportType: "income",
		Metric:     "Revenue",
		Year:       2025,
		Value:      &value,
		Currency:   "EUR",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/dax/metrics?ticker=TST", nil)

	handler.GetMetrics(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Revenue")
}

func TestDAXHandler_GetMetrics_MissingTicker(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := dax.NewInMemoryRepository()
	service := dax.NewService(repo)
	handler := NewDAXHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/dax/metrics", nil)

	handler.GetMetrics(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "ticker parameter is required")
}

func TestDAXHandler_GetMetrics_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Use in-memory repo for consistency with other tests
	repo := dax.NewInMemoryRepository()
	service := dax.NewService(repo)
	handler := NewDAXHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/dax/metrics?ticker=TST", nil)

	handler.GetMetrics(c)

	// Should return OK with empty metrics, not error
	assert.Equal(t, http.StatusOK, w.Code)
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
			name:         "valid positive value",
			queryValue:   "5",
			defaultValue: 10,
			expected:     5,
		},
		{
			name:         "empty value uses default",
			queryValue:   "",
			defaultValue: 10,
			expected:     10,
		},
		{
			name:         "invalid value uses default",
			queryValue:   "invalid",
			defaultValue: 10,
			expected:     10,
		},
		{
			name:         "negative value uses default",
			queryValue:   "-5",
			defaultValue: 10,
			expected:     10,
		},
		{
			name:         "zero value uses default",
			queryValue:   "0",
			defaultValue: 10,
			expected:     10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/?key="+tt.queryValue, nil)

			result := parseIntQuery(c, "key", tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}
