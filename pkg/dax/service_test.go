package dax

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(record *DAXRecord) error {
	args := m.Called(record)
	return args.Error(0)
}

func (m *MockRepository) BulkUpsert(records []DAXRecord) error {
	args := m.Called(records)
	return args.Error(0)
}

func (m *MockRepository) FindAll(page, limit int) ([]DAXRecord, int, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]DAXRecord), args.Int(1), args.Error(2)
}

func (m *MockRepository) FindByFilters(ticker string, year *int, page, limit int) ([]DAXRecord, int, error) {
	args := m.Called(ticker, year, page, limit)
	return args.Get(0).([]DAXRecord), args.Int(1), args.Error(2)
}

func (m *MockRepository) GetMetrics(ticker string) ([]string, error) {
	args := m.Called(ticker)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRepository) DeleteAll() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRepository) Count() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func TestImportCSV_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	csvContent := `company,ticker,report_type,metric,year,value,currency
Siemens AG,SIE,income,EBITDA,2025,15859000000.0,EUR
SAP SE,SAP,income,Net Income,2025,8500000000.0,EUR`

	mockRepo.On("BulkUpsert", mock.AnythingOfType("[]dax.DAXRecord")).Return(nil)

	reader := bytes.NewBufferString(csvContent)
	response, err := service.ImportCSV(reader)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, response.RecordsImported)
	mockRepo.AssertExpectations(t)
}

func TestImportCSV_MissingRequiredFields(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	csvContent := `company,ticker,metric,year,value
Siemens AG,SIE,EBITDA,2025,15859000000.0`

	reader := bytes.NewBufferString(csvContent)
	response, err := service.ImportCSV(reader)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "missing required fields")
}

func TestImportCSV_InvalidYear(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	csvContent := `company,ticker,report_type,metric,year,value,currency
Siemens AG,SIE,income,EBITDA,invalid,15859000000.0,EUR`

	reader := bytes.NewBufferString(csvContent)
	response, err := service.ImportCSV(reader)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "invalid year")
}

func TestGetAll_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	expectedRecords := []DAXRecord{
		{Ticker: "SIE", Year: 2025},
		{Ticker: "SAP", Year: 2025},
	}

	mockRepo.On("FindAll", 1, 10).Return(expectedRecords, 2, nil)

	response, err := service.GetAll(1, 10)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, len(response.Data))
	assert.Equal(t, 2, response.Pagination.TotalCount)
	mockRepo.AssertExpectations(t)
}

func TestGetByFilters_WithTickerAndYear(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	year := 2025
	expectedRecords := []DAXRecord{
		{Ticker: "SIE", Year: 2025},
	}

	mockRepo.On("FindByFilters", "SIE", &year, 1, 10).
		Return(expectedRecords, 1, nil)

	response, err := service.GetByFilters("SIE", &year, 1, 10)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 1, len(response.Data))
	mockRepo.AssertExpectations(t)
}

func TestGetMetrics_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	expectedMetrics := []string{"EBITDA", "Net Income"}
	mockRepo.On("GetMetrics", "SIE").Return(expectedMetrics, nil)

	response, err := service.GetMetrics("SIE")

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "SIE", response.Ticker)
	assert.Equal(t, 2, len(response.Metrics))
	mockRepo.AssertExpectations(t)
}

func TestGetMetrics_EmptyTicker(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	response, err := service.GetMetrics("")

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "ticker is required")
}
