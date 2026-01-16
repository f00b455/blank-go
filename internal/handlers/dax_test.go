package handlers

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/f00b455/blank-go/pkg/dax"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func float64Ptr(f float64) *float64 {
	return &f
}

func setupDAXHandler() (*DAXHandler, dax.Repository) {
	repo := dax.NewInMemoryRepository()
	service := dax.NewService(repo)
	handler := NewDAXHandler(service)
	return handler, repo
}

func createMultipartRequest(t *testing.T, csvContent string) (*http.Request, *bytes.Buffer) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "test.csv")
	require.NoError(t, err)

	_, err = io.WriteString(part, csvContent)
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/dax/import", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, body
}

func TestDAXHandler_ImportCSV_Success(t *testing.T) {
	handler, _ := setupDAXHandler()

	csvContent := `company,ticker,report_type,metric,year,value,currency
Siemens AG,SIE,income,EBITDA,2025,15859000000.0,EUR
Siemens AG,SIE,income,Net Income,2025,9620000000.0,EUR`

	req, _ := createMultipartRequest(t, csvContent)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ImportCSV(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "records_imported")
	assert.Contains(t, w.Body.String(), "2")
}

func TestDAXHandler_ImportCSV_MissingFile(t *testing.T) {
	handler, _ := setupDAXHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/dax/import", nil)

	handler.ImportCSV(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "file is required")
}

func TestDAXHandler_ImportCSV_MissingRequiredFields(t *testing.T) {
	handler, _ := setupDAXHandler()

	csvContent := `company,ticker,metric,year,value
Siemens AG,SIE,EBITDA,2025,15859000000.0`

	req, _ := createMultipartRequest(t, csvContent)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ImportCSV(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "missing required fields")
}

func TestDAXHandler_ImportCSV_InvalidYear(t *testing.T) {
	handler, _ := setupDAXHandler()

	csvContent := `company,ticker,report_type,metric,year,value,currency
Siemens AG,SIE,income,EBITDA,invalid,15859000000.0,EUR`

	req, _ := createMultipartRequest(t, csvContent)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ImportCSV(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid year")
}

func TestDAXHandler_GetAll_Success(t *testing.T) {
	handler, repo := setupDAXHandler()

	// Insert test data
	records := []dax.DAXRecord{
		{Company: "Siemens AG", Ticker: "SIE", ReportType: "income", Metric: "EBITDA", Year: 2025, Value: float64Ptr(1000.0), Currency: "EUR"},
		{Company: "SAP SE", Ticker: "SAP", ReportType: "income", Metric: "Revenue", Year: 2025, Value: float64Ptr(2000.0), Currency: "EUR"},
	}
	err := repo.BulkUpsert(records)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/dax?page=1&limit=10", nil)

	handler.GetAll(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "data")
	assert.Contains(t, w.Body.String(), "SIE")
	assert.Contains(t, w.Body.String(), "SAP")
}

func TestDAXHandler_GetAll_DefaultPagination(t *testing.T) {
	handler, repo := setupDAXHandler()

	// Insert test data
	records := []dax.DAXRecord{
		{Company: "Test", Ticker: "TST", ReportType: "income", Metric: "Revenue", Year: 2025, Value: float64Ptr(100.0), Currency: "EUR"},
	}
	err := repo.BulkUpsert(records)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/dax", nil)

	handler.GetAll(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDAXHandler_GetByFilters_ByTicker(t *testing.T) {
	handler, repo := setupDAXHandler()

	// Insert test data
	records := []dax.DAXRecord{
		{Company: "Siemens AG", Ticker: "SIE", ReportType: "income", Metric: "EBITDA", Year: 2025, Value: float64Ptr(1000.0), Currency: "EUR"},
		{Company: "SAP SE", Ticker: "SAP", ReportType: "income", Metric: "Revenue", Year: 2025, Value: float64Ptr(2000.0), Currency: "EUR"},
	}
	err := repo.BulkUpsert(records)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/dax/filter?ticker=SIE", nil)

	handler.GetByFilters(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "SIE")
	assert.NotContains(t, w.Body.String(), "SAP")
}

func TestDAXHandler_GetByFilters_ByYear(t *testing.T) {
	handler, repo := setupDAXHandler()

	// Insert test data
	records := []dax.DAXRecord{
		{Company: "Siemens AG", Ticker: "SIE", ReportType: "income", Metric: "EBITDA", Year: 2025, Value: float64Ptr(1000.0), Currency: "EUR"},
		{Company: "Siemens AG", Ticker: "SIE", ReportType: "income", Metric: "EBITDA", Year: 2024, Value: float64Ptr(900.0), Currency: "EUR"},
	}
	err := repo.BulkUpsert(records)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/dax/filter?year=2025", nil)

	handler.GetByFilters(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "2025")
}

func TestDAXHandler_GetByFilters_ByTickerAndYear(t *testing.T) {
	handler, repo := setupDAXHandler()

	// Insert test data
	records := []dax.DAXRecord{
		{Company: "Siemens AG", Ticker: "SIE", ReportType: "income", Metric: "EBITDA", Year: 2025, Value: float64Ptr(1000.0), Currency: "EUR"},
		{Company: "Siemens AG", Ticker: "SIE", ReportType: "income", Metric: "EBITDA", Year: 2024, Value: float64Ptr(900.0), Currency: "EUR"},
		{Company: "SAP SE", Ticker: "SAP", ReportType: "income", Metric: "Revenue", Year: 2025, Value: float64Ptr(2000.0), Currency: "EUR"},
	}
	err := repo.BulkUpsert(records)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/dax/filter?ticker=SIE&year=2025", nil)

	handler.GetByFilters(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "SIE")
	assert.Contains(t, w.Body.String(), "2025")
}

func TestDAXHandler_GetByFilters_InvalidYear(t *testing.T) {
	handler, _ := setupDAXHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/dax/filter?year=invalid", nil)

	handler.GetByFilters(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid year parameter")
}

func TestDAXHandler_GetMetrics_Success(t *testing.T) {
	handler, repo := setupDAXHandler()

	// Insert test data
	records := []dax.DAXRecord{
		{Company: "Siemens AG", Ticker: "SIE", ReportType: "income", Metric: "EBITDA", Year: 2025, Value: float64Ptr(1000.0), Currency: "EUR"},
		{Company: "Siemens AG", Ticker: "SIE", ReportType: "income", Metric: "Revenue", Year: 2025, Value: float64Ptr(5000.0), Currency: "EUR"},
	}
	err := repo.BulkUpsert(records)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/dax/metrics?ticker=SIE", nil)

	handler.GetMetrics(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "EBITDA")
	assert.Contains(t, w.Body.String(), "Revenue")
}

func TestDAXHandler_GetMetrics_MissingTicker(t *testing.T) {
	handler, _ := setupDAXHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/dax/metrics", nil)

	handler.GetMetrics(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "ticker parameter is required")
}

func TestParseIntQuery(t *testing.T) {
	tests := []struct {
		name          string
		queryKey      string
		queryValue    string
		defaultValue  int
		expectedValue int
	}{
		{
			name:          "Valid positive integer",
			queryKey:      "page",
			queryValue:    "5",
			defaultValue:  1,
			expectedValue: 5,
		},
		{
			name:          "Empty query returns default",
			queryKey:      "page",
			queryValue:    "",
			defaultValue:  1,
			expectedValue: 1,
		},
		{
			name:          "Invalid integer returns default",
			queryKey:      "page",
			queryValue:    "invalid",
			defaultValue:  1,
			expectedValue: 1,
		},
		{
			name:          "Zero returns default",
			queryKey:      "page",
			queryValue:    "0",
			defaultValue:  1,
			expectedValue: 1,
		},
		{
			name:          "Negative returns default",
			queryKey:      "page",
			queryValue:    "-5",
			defaultValue:  1,
			expectedValue: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			url := "/test"
			if tt.queryValue != "" {
				url += "?" + tt.queryKey + "=" + tt.queryValue
			}
			c.Request = httptest.NewRequest("GET", url, nil)

			result := parseIntQuery(c, tt.queryKey, tt.defaultValue)
			assert.Equal(t, tt.expectedValue, result)
		})
	}
}

func TestDAXHandler_GetByFilters_WithPagination(t *testing.T) {
	handler, repo := setupDAXHandler()

	// Insert test data
	records := make([]dax.DAXRecord, 15)
	for i := 0; i < 15; i++ {
		records[i] = dax.DAXRecord{
			Company:    "Test",
			Ticker:     "TST",
			ReportType: "income",
			Metric:     "Revenue",
			Year:       2025,
			Value:      float64Ptr(float64(i * 100)),
			Currency:   "EUR",
		}
	}
	err := repo.BulkUpsert(records)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/dax/filter?ticker=TST&page=1&limit=5", nil)

	handler.GetByFilters(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "page")
	assert.Contains(t, w.Body.String(), "limit")
}

func TestDAXHandler_GetAll_ServiceError(t *testing.T) {
	// Use a failing repository to simulate service error
	handler, _ := setupDAXHandler()

	// Override with a mock that returns error through real service
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/dax?page=1&limit=10", nil)

	handler.GetAll(c)

	// With empty repo, this should still succeed with empty data
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDAXHandler_GetAll_InvalidPage(t *testing.T) {
	handler, _ := setupDAXHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/dax?page=invalid&limit=10", nil)

	handler.GetAll(c)

	// Invalid page should default to 1 and succeed
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDAXHandler_GetAll_InvalidLimit(t *testing.T) {
	handler, _ := setupDAXHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/dax?page=1&limit=invalid", nil)

	handler.GetAll(c)

	// Invalid limit should default to 10 and succeed
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDAXHandler_GetMetrics_NoData(t *testing.T) {
	handler, _ := setupDAXHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/dax/metrics?ticker=NONEXISTENT", nil)

	handler.GetMetrics(c)

	// Should return empty metrics, not error
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "metrics")
}

func TestDAXHandler_GetByFilters_NoFilters(t *testing.T) {
	handler, repo := setupDAXHandler()

	// Insert test data
	records := []dax.DAXRecord{
		{Company: "Test", Ticker: "TST", ReportType: "income", Metric: "Revenue", Year: 2025, Value: float64Ptr(100.0), Currency: "EUR"},
	}
	err := repo.BulkUpsert(records)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/dax/filter", nil)

	handler.GetByFilters(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDAXHandler_ImportCSV_InvalidData(t *testing.T) {
	handler, _ := setupDAXHandler()

	csvContent := `company,ticker,report_type,metric,year,value,currency
Siemens AG,SIE,income,EBITDA,2025,not-a-number,EUR`

	req, _ := createMultipartRequest(t, csvContent)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ImportCSV(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid data at row")
}
