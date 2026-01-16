package dax_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/f00b455/blank-go/pkg/dax"
	"github.com/f00b455/blank-go/pkg/dax/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewService(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	assert.NotNil(t, service)
}

func TestImportCSV_Success(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	csvContent := `company,ticker,report_type,metric,year,value,currency
Siemens AG,SIE,income,EBITDA,2025,15859000000.0,EUR
SAP SE,SAP,income,Net Income,2025,8500000000.0,EUR`

	mockRepo.On("BulkUpsert", mock.AnythingOfType("[]dax.DAXRecord")).Return(nil)

	reader := bytes.NewBufferString(csvContent)
	response, err := service.ImportCSV(reader)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, response.RecordsImported)
	assert.Contains(t, response.Message, "Successfully imported 2 records")
}

func TestImportCSV_MissingRequiredFields(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	csvContent := `company,ticker,metric,year,value
Siemens AG,SIE,EBITDA,2025,15859000000.0`

	reader := bytes.NewBufferString(csvContent)
	response, err := service.ImportCSV(reader)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "missing required fields")
}

func TestImportCSV_InvalidYear(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	csvContent := `company,ticker,report_type,metric,year,value,currency
Siemens AG,SIE,income,EBITDA,invalid,15859000000.0,EUR`

	reader := bytes.NewBufferString(csvContent)
	response, err := service.ImportCSV(reader)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "invalid year")
}

func TestImportCSV_InvalidValue(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	csvContent := `company,ticker,report_type,metric,year,value,currency
Siemens AG,SIE,income,EBITDA,2025,not-a-number,EUR`

	reader := bytes.NewBufferString(csvContent)
	response, err := service.ImportCSV(reader)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "invalid value")
}

func TestImportCSV_EmptyCSV(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	csvContent := `company,ticker,report_type,metric,year,value,currency`

	reader := bytes.NewBufferString(csvContent)
	response, err := service.ImportCSV(reader)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "no records found")
}

func TestImportCSV_InsufficientColumns(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	csvContent := `company,ticker,report_type,metric,year,value,currency
Siemens AG,SIE,income`

	reader := bytes.NewBufferString(csvContent)
	response, err := service.ImportCSV(reader)

	assert.Error(t, err)
	assert.Nil(t, response)
	// CSV parser detects wrong number of fields before our validation
	assert.Contains(t, err.Error(), "wrong number of fields")
}

func TestImportCSV_BulkUpsertError(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	csvContent := `company,ticker,report_type,metric,year,value,currency
Siemens AG,SIE,income,EBITDA,2025,15859000000.0,EUR`

	mockRepo.On("BulkUpsert", mock.AnythingOfType("[]dax.DAXRecord")).
		Return(errors.New("database error"))

	reader := bytes.NewBufferString(csvContent)
	response, err := service.ImportCSV(reader)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to import records")
}

func TestImportCSV_EmptyReader(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	reader := bytes.NewBufferString("")
	response, err := service.ImportCSV(reader)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to read CSV header")
}

func TestGetAll_Success(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	expectedRecords := []dax.DAXRecord{
		{Ticker: "SIE", Year: 2025},
		{Ticker: "SAP", Year: 2025},
	}

	mockRepo.On("FindAll", 1, 10).Return(expectedRecords, 2, nil)

	response, err := service.GetAll(1, 10)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, len(response.Data))
	assert.Equal(t, 2, response.Pagination.TotalCount)
	assert.Equal(t, 1, response.Pagination.TotalPages)
}

func TestGetAll_PageLessThanOne(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	expectedRecords := []dax.DAXRecord{{Ticker: "SIE", Year: 2025}}

	// When page < 1, it should default to 1
	mockRepo.On("FindAll", 1, 10).Return(expectedRecords, 1, nil)

	response, err := service.GetAll(0, 10)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 1, response.Pagination.Page)
}

func TestGetAll_NegativePage(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	expectedRecords := []dax.DAXRecord{{Ticker: "SIE", Year: 2025}}

	mockRepo.On("FindAll", 1, 10).Return(expectedRecords, 1, nil)

	response, err := service.GetAll(-5, 10)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 1, response.Pagination.Page)
}

func TestGetAll_LimitLessThanOne(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	expectedRecords := []dax.DAXRecord{{Ticker: "SIE", Year: 2025}}

	// When limit < 1, it should default to 10
	mockRepo.On("FindAll", 1, 10).Return(expectedRecords, 1, nil)

	response, err := service.GetAll(1, 0)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 10, response.Pagination.Limit)
}

func TestGetAll_LimitGreaterThan100(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	expectedRecords := []dax.DAXRecord{{Ticker: "SIE", Year: 2025}}

	// When limit > 100, it should default to 10
	mockRepo.On("FindAll", 1, 10).Return(expectedRecords, 1, nil)

	response, err := service.GetAll(1, 150)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 10, response.Pagination.Limit)
}

func TestGetAll_RepositoryError(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	mockRepo.On("FindAll", 1, 10).Return([]dax.DAXRecord{}, 0, errors.New("database error"))

	response, err := service.GetAll(1, 10)

	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestGetAll_TotalPagesCalculation(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	expectedRecords := []dax.DAXRecord{{Ticker: "SIE", Year: 2025}}

	// 25 total records with limit 10 = 3 pages
	mockRepo.On("FindAll", 1, 10).Return(expectedRecords, 25, nil)

	response, err := service.GetAll(1, 10)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 3, response.Pagination.TotalPages)
}

func TestGetByFilters_WithTickerAndYear(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	year := 2025
	expectedRecords := []dax.DAXRecord{
		{Ticker: "SIE", Year: 2025},
	}

	mockRepo.On("FindByFilters", "SIE", &year, 1, 10).
		Return(expectedRecords, 1, nil)

	response, err := service.GetByFilters("SIE", &year, 1, 10)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 1, len(response.Data))
}

func TestGetByFilters_WithOnlyTicker(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	expectedRecords := []dax.DAXRecord{
		{Ticker: "SIE", Year: 2024},
		{Ticker: "SIE", Year: 2025},
	}

	mockRepo.On("FindByFilters", "SIE", (*int)(nil), 1, 10).
		Return(expectedRecords, 2, nil)

	response, err := service.GetByFilters("SIE", nil, 1, 10)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, len(response.Data))
}

func TestGetByFilters_PageLessThanOne(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	expectedRecords := []dax.DAXRecord{{Ticker: "SIE", Year: 2025}}

	mockRepo.On("FindByFilters", "SIE", (*int)(nil), 1, 10).
		Return(expectedRecords, 1, nil)

	response, err := service.GetByFilters("SIE", nil, 0, 10)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 1, response.Pagination.Page)
}

func TestGetByFilters_NegativePage(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	expectedRecords := []dax.DAXRecord{{Ticker: "SIE", Year: 2025}}

	mockRepo.On("FindByFilters", "SIE", (*int)(nil), 1, 10).
		Return(expectedRecords, 1, nil)

	response, err := service.GetByFilters("SIE", nil, -10, 10)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 1, response.Pagination.Page)
}

func TestGetByFilters_LimitLessThanOne(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	expectedRecords := []dax.DAXRecord{{Ticker: "SIE", Year: 2025}}

	mockRepo.On("FindByFilters", "SIE", (*int)(nil), 1, 10).
		Return(expectedRecords, 1, nil)

	response, err := service.GetByFilters("SIE", nil, 1, 0)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 10, response.Pagination.Limit)
}

func TestGetByFilters_LimitGreaterThan100(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	expectedRecords := []dax.DAXRecord{{Ticker: "SIE", Year: 2025}}

	mockRepo.On("FindByFilters", "SIE", (*int)(nil), 1, 10).
		Return(expectedRecords, 1, nil)

	response, err := service.GetByFilters("SIE", nil, 1, 200)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 10, response.Pagination.Limit)
}

func TestGetByFilters_RepositoryError(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	mockRepo.On("FindByFilters", "SIE", (*int)(nil), 1, 10).
		Return([]dax.DAXRecord{}, 0, errors.New("database error"))

	response, err := service.GetByFilters("SIE", nil, 1, 10)

	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestGetByFilters_TotalPagesCalculation(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	expectedRecords := []dax.DAXRecord{{Ticker: "SIE", Year: 2025}}

	// 45 total records with limit 10 = 5 pages
	mockRepo.On("FindByFilters", "SIE", (*int)(nil), 1, 10).
		Return(expectedRecords, 45, nil)

	response, err := service.GetByFilters("SIE", nil, 1, 10)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 5, response.Pagination.TotalPages)
}

func TestGetMetrics_Success(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	expectedMetrics := []string{"EBITDA", "Net Income"}
	mockRepo.On("GetMetrics", "SIE").Return(expectedMetrics, nil)

	response, err := service.GetMetrics("SIE")

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "SIE", response.Ticker)
	assert.Equal(t, 2, len(response.Metrics))
}

func TestGetMetrics_EmptyTicker(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	response, err := service.GetMetrics("")

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "ticker is required")
}

func TestGetMetrics_RepositoryError(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	mockRepo.On("GetMetrics", "SIE").Return([]string(nil), errors.New("database error"))

	response, err := service.GetMetrics("SIE")

	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestGetMetrics_NilMetricsReturnsEmptySlice(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	service := dax.NewService(mockRepo)

	// Repository returns nil metrics (no data found)
	mockRepo.On("GetMetrics", "UNKNOWN").Return([]string(nil), nil)

	response, err := service.GetMetrics("UNKNOWN")

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "UNKNOWN", response.Ticker)
	assert.NotNil(t, response.Metrics)
	assert.Empty(t, response.Metrics)
}
