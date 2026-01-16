package dax

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryRepository_Create(t *testing.T) {
	repo := NewInMemoryRepository()

	value := 1000.0
	record := &DAXRecord{
		Company:    "Test Company",
		Ticker:     "TST",
		ReportType: "income",
		Metric:     "Revenue",
		Year:       2025,
		Value:      &value,
		Currency:   "EUR",
	}

	err := repo.Create(record)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, record.ID)

	count, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestInMemoryRepository_Create_WithID(t *testing.T) {
	repo := NewInMemoryRepository()

	existingID := uuid.New()
	value := 1000.0
	record := &DAXRecord{
		ID:         existingID,
		Company:    "Test Company",
		Ticker:     "TST",
		ReportType: "income",
		Metric:     "Revenue",
		Year:       2025,
		Value:      &value,
		Currency:   "EUR",
	}

	err := repo.Create(record)
	assert.NoError(t, err)
	assert.Equal(t, existingID, record.ID)
}

func TestInMemoryRepository_BulkUpsert_Insert(t *testing.T) {
	repo := NewInMemoryRepository()

	value1 := 1000.0
	value2 := 2000.0
	records := []DAXRecord{
		{
			Company:    "Company A",
			Ticker:     "AAA",
			ReportType: "income",
			Metric:     "Revenue",
			Year:       2025,
			Value:      &value1,
			Currency:   "EUR",
		},
		{
			Company:    "Company B",
			Ticker:     "BBB",
			ReportType: "income",
			Metric:     "Revenue",
			Year:       2025,
			Value:      &value2,
			Currency:   "EUR",
		},
	}

	err := repo.BulkUpsert(records)
	assert.NoError(t, err)

	count, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestInMemoryRepository_BulkUpsert_Update(t *testing.T) {
	repo := NewInMemoryRepository()

	value1 := 1000.0
	records1 := []DAXRecord{
		{
			Company:    "Company A",
			Ticker:     "AAA",
			ReportType: "income",
			Metric:     "Revenue",
			Year:       2025,
			Value:      &value1,
			Currency:   "EUR",
		},
	}

	err := repo.BulkUpsert(records1)
	require.NoError(t, err)

	value2 := 2000.0
	records2 := []DAXRecord{
		{
			Company:    "Company A",
			Ticker:     "AAA",
			ReportType: "income",
			Metric:     "Revenue",
			Year:       2025,
			Value:      &value2,
			Currency:   "USD",
		},
	}

	err = repo.BulkUpsert(records2)
	assert.NoError(t, err)

	count, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	allRecords, total, err := repo.FindAll(1, 10)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Equal(t, "USD", allRecords[0].Currency)
	assert.Equal(t, value2, *allRecords[0].Value)
}

func TestInMemoryRepository_BulkUpsert_EmptySlice(t *testing.T) {
	repo := NewInMemoryRepository()

	err := repo.BulkUpsert([]DAXRecord{})
	assert.NoError(t, err)

	count, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestInMemoryRepository_FindAll(t *testing.T) {
	repo := NewInMemoryRepository()

	value1 := 1000.0
	value2 := 2000.0
	value3 := 3000.0
	records := []DAXRecord{
		{
			Company:    "Company A",
			Ticker:     "AAA",
			ReportType: "income",
			Metric:     "Revenue",
			Year:       2024,
			Value:      &value1,
			Currency:   "EUR",
		},
		{
			Company:    "Company B",
			Ticker:     "BBB",
			ReportType: "income",
			Metric:     "EBITDA",
			Year:       2025,
			Value:      &value2,
			Currency:   "EUR",
		},
		{
			Company:    "Company C",
			Ticker:     "AAA",
			ReportType: "income",
			Metric:     "EBITDA",
			Year:       2025,
			Value:      &value3,
			Currency:   "EUR",
		},
	}

	err := repo.BulkUpsert(records)
	require.NoError(t, err)

	result, total, err := repo.FindAll(1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Equal(t, 3, len(result))

	// Verify sorting: year DESC, ticker ASC, metric ASC
	assert.Equal(t, 2025, result[0].Year)
	assert.Equal(t, "AAA", result[0].Ticker)
	assert.Equal(t, 2025, result[1].Year)
	assert.Equal(t, "BBB", result[1].Ticker)
	assert.Equal(t, 2024, result[2].Year)
}

func TestInMemoryRepository_FindAll_Pagination(t *testing.T) {
	repo := NewInMemoryRepository()

	// Create 5 records
	for i := 1; i <= 5; i++ {
		value := float64(i * 1000)
		err := repo.Create(&DAXRecord{
			Company:    "Company",
			Ticker:     "TST",
			ReportType: "income",
			Metric:     "Metric",
			Year:       2020 + i,
			Value:      &value,
			Currency:   "EUR",
		})
		require.NoError(t, err)
	}

	// Page 1: 2 items
	result, total, err := repo.FindAll(1, 2)
	assert.NoError(t, err)
	assert.Equal(t, 5, total)
	assert.Equal(t, 2, len(result))

	// Page 2: 2 items
	result, total, err = repo.FindAll(2, 2)
	assert.NoError(t, err)
	assert.Equal(t, 5, total)
	assert.Equal(t, 2, len(result))

	// Page 3: 1 item
	result, total, err = repo.FindAll(3, 2)
	assert.NoError(t, err)
	assert.Equal(t, 5, total)
	assert.Equal(t, 1, len(result))

	// Page 4: 0 items (beyond range)
	result, total, err = repo.FindAll(4, 2)
	assert.NoError(t, err)
	assert.Equal(t, 5, total)
	assert.Equal(t, 0, len(result))
}

func TestInMemoryRepository_FindByFilters_TickerOnly(t *testing.T) {
	repo := NewInMemoryRepository()

	value1 := 1000.0
	value2 := 2000.0
	records := []DAXRecord{
		{
			Company:    "Company A",
			Ticker:     "AAA",
			ReportType: "income",
			Metric:     "Revenue",
			Year:       2025,
			Value:      &value1,
			Currency:   "EUR",
		},
		{
			Company:    "Company B",
			Ticker:     "BBB",
			ReportType: "income",
			Metric:     "Revenue",
			Year:       2025,
			Value:      &value2,
			Currency:   "EUR",
		},
	}

	err := repo.BulkUpsert(records)
	require.NoError(t, err)

	result, total, err := repo.FindByFilters("AAA", nil, 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "AAA", result[0].Ticker)
}

func TestInMemoryRepository_FindByFilters_YearOnly(t *testing.T) {
	repo := NewInMemoryRepository()

	value1 := 1000.0
	value2 := 2000.0
	records := []DAXRecord{
		{
			Company:    "Company A",
			Ticker:     "AAA",
			ReportType: "income",
			Metric:     "Revenue",
			Year:       2024,
			Value:      &value1,
			Currency:   "EUR",
		},
		{
			Company:    "Company B",
			Ticker:     "BBB",
			ReportType: "income",
			Metric:     "Revenue",
			Year:       2025,
			Value:      &value2,
			Currency:   "EUR",
		},
	}

	err := repo.BulkUpsert(records)
	require.NoError(t, err)

	year := 2025
	result, total, err := repo.FindByFilters("", &year, 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, 2025, result[0].Year)
}

func TestInMemoryRepository_FindByFilters_TickerAndYear(t *testing.T) {
	repo := NewInMemoryRepository()

	value1 := 1000.0
	value2 := 2000.0
	value3 := 3000.0
	records := []DAXRecord{
		{
			Company:    "Company A",
			Ticker:     "AAA",
			ReportType: "income",
			Metric:     "Revenue",
			Year:       2024,
			Value:      &value1,
			Currency:   "EUR",
		},
		{
			Company:    "Company A",
			Ticker:     "AAA",
			ReportType: "income",
			Metric:     "EBITDA",
			Year:       2025,
			Value:      &value2,
			Currency:   "EUR",
		},
		{
			Company:    "Company B",
			Ticker:     "BBB",
			ReportType: "income",
			Metric:     "Revenue",
			Year:       2025,
			Value:      &value3,
			Currency:   "EUR",
		},
	}

	err := repo.BulkUpsert(records)
	require.NoError(t, err)

	year := 2025
	result, total, err := repo.FindByFilters("AAA", &year, 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "AAA", result[0].Ticker)
	assert.Equal(t, 2025, result[0].Year)
}

func TestInMemoryRepository_FindByFilters_NoMatch(t *testing.T) {
	repo := NewInMemoryRepository()

	value := 1000.0
	records := []DAXRecord{
		{
			Company:    "Company A",
			Ticker:     "AAA",
			ReportType: "income",
			Metric:     "Revenue",
			Year:       2025,
			Value:      &value,
			Currency:   "EUR",
		},
	}

	err := repo.BulkUpsert(records)
	require.NoError(t, err)

	result, total, err := repo.FindByFilters("XXX", nil, 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Equal(t, 0, len(result))
}

func TestInMemoryRepository_GetMetrics(t *testing.T) {
	repo := NewInMemoryRepository()

	value1 := 1000.0
	value2 := 2000.0
	value3 := 3000.0
	records := []DAXRecord{
		{
			Company:    "Company A",
			Ticker:     "AAA",
			ReportType: "income",
			Metric:     "Revenue",
			Year:       2025,
			Value:      &value1,
			Currency:   "EUR",
		},
		{
			Company:    "Company A",
			Ticker:     "AAA",
			ReportType: "income",
			Metric:     "EBITDA",
			Year:       2025,
			Value:      &value2,
			Currency:   "EUR",
		},
		{
			Company:    "Company B",
			Ticker:     "BBB",
			ReportType: "income",
			Metric:     "Net Income",
			Year:       2025,
			Value:      &value3,
			Currency:   "EUR",
		},
	}

	err := repo.BulkUpsert(records)
	require.NoError(t, err)

	metrics, err := repo.GetMetrics("AAA")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(metrics))
	assert.Contains(t, metrics, "Revenue")
	assert.Contains(t, metrics, "EBITDA")

	// Verify sorted
	assert.True(t, metrics[0] < metrics[1])
}

func TestInMemoryRepository_GetMetrics_NoMatch(t *testing.T) {
	repo := NewInMemoryRepository()

	value := 1000.0
	records := []DAXRecord{
		{
			Company:    "Company A",
			Ticker:     "AAA",
			ReportType: "income",
			Metric:     "Revenue",
			Year:       2025,
			Value:      &value,
			Currency:   "EUR",
		},
	}

	err := repo.BulkUpsert(records)
	require.NoError(t, err)

	metrics, err := repo.GetMetrics("XXX")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(metrics))
}

func TestInMemoryRepository_DeleteAll(t *testing.T) {
	repo := NewInMemoryRepository()

	value := 1000.0
	records := []DAXRecord{
		{
			Company:    "Company A",
			Ticker:     "AAA",
			ReportType: "income",
			Metric:     "Revenue",
			Year:       2025,
			Value:      &value,
			Currency:   "EUR",
		},
	}

	err := repo.BulkUpsert(records)
	require.NoError(t, err)

	count, err := repo.Count()
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	err = repo.DeleteAll()
	assert.NoError(t, err)

	count, err = repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestInMemoryRepository_Count(t *testing.T) {
	repo := NewInMemoryRepository()

	count, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, 0, count)

	value1 := 1000.0
	value2 := 2000.0
	records := []DAXRecord{
		{
			Company:    "Company A",
			Ticker:     "AAA",
			ReportType: "income",
			Metric:     "Revenue",
			Year:       2025,
			Value:      &value1,
			Currency:   "EUR",
		},
		{
			Company:    "Company B",
			Ticker:     "BBB",
			ReportType: "income",
			Metric:     "Revenue",
			Year:       2025,
			Value:      &value2,
			Currency:   "EUR",
		},
	}

	err = repo.BulkUpsert(records)
	require.NoError(t, err)

	count, err = repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}
