package dax

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func floatPtr(f float64) *float64 {
	return &f
}

func TestNewInMemoryRepository(t *testing.T) {
	repo := NewInMemoryRepository()
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.records)
}

func TestInMemoryRepository_Create(t *testing.T) {
	repo := NewInMemoryRepository()

	record := &DAXRecord{
		Company:    "Test GmbH",
		Ticker:     "TEST",
		Metric:     "Revenue",
		Year:       2024,
		ReportType: "Q4",
		Value:      floatPtr(100.0),
		Currency:   "EUR",
	}

	err := repo.Create(record)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, record.ID)

	// Verify record was stored
	count, err := repo.Count()
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestInMemoryRepository_Create_PreserveExistingID(t *testing.T) {
	repo := NewInMemoryRepository()

	existingID := uuid.New()
	record := &DAXRecord{
		ID:         existingID,
		Company:    "Test GmbH",
		Ticker:     "TEST",
		Metric:     "Revenue",
		Year:       2024,
		ReportType: "Q4",
		Value:      floatPtr(100.0),
		Currency:   "EUR",
	}

	err := repo.Create(record)
	require.NoError(t, err)
	assert.Equal(t, existingID, record.ID)
}

func TestInMemoryRepository_BulkUpsert(t *testing.T) {
	repo := NewInMemoryRepository()

	records := []DAXRecord{
		{
			Company:    "Company A",
			Ticker:     "AAA",
			Metric:     "Revenue",
			Year:       2024,
			ReportType: "Q4",
			Value:      floatPtr(100.0),
			Currency:   "EUR",
		},
		{
			Company:    "Company B",
			Ticker:     "BBB",
			Metric:     "Profit",
			Year:       2024,
			ReportType: "Q4",
			Value:      floatPtr(50.0),
			Currency:   "EUR",
		},
	}

	err := repo.BulkUpsert(records)
	require.NoError(t, err)

	count, err := repo.Count()
	require.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestInMemoryRepository_BulkUpsert_Updates(t *testing.T) {
	repo := NewInMemoryRepository()

	// Insert initial record
	initial := []DAXRecord{
		{
			Company:    "Company A",
			Ticker:     "AAA",
			Metric:     "Revenue",
			Year:       2024,
			ReportType: "Q4",
			Value:      floatPtr(100.0),
			Currency:   "EUR",
		},
	}

	err := repo.BulkUpsert(initial)
	require.NoError(t, err)

	// Update the same record (same unique key)
	updated := []DAXRecord{
		{
			Company:    "Company A",
			Ticker:     "AAA",
			Metric:     "Revenue",
			Year:       2024,
			ReportType: "Q4",
			Value:      floatPtr(200.0), // Different value
			Currency:   "EUR",
		},
	}

	err = repo.BulkUpsert(updated)
	require.NoError(t, err)

	// Should still have only 1 record (updated, not duplicated)
	count, err := repo.Count()
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	// Verify the value was updated
	records, _, err := repo.FindAll(1, 10)
	require.NoError(t, err)
	require.NotNil(t, records[0].Value)
	assert.Equal(t, 200.0, *records[0].Value)
}

func TestInMemoryRepository_FindAll(t *testing.T) {
	repo := NewInMemoryRepository()

	// Insert test data
	records := []DAXRecord{
		{Company: "A", Ticker: "AAA", Metric: "Revenue", Year: 2024, Value: floatPtr(100.0), Currency: "EUR"},
		{Company: "B", Ticker: "BBB", Metric: "Profit", Year: 2023, Value: floatPtr(50.0), Currency: "EUR"},
		{Company: "C", Ticker: "CCC", Metric: "EBIT", Year: 2024, Value: floatPtr(75.0), Currency: "EUR"},
	}

	err := repo.BulkUpsert(records)
	require.NoError(t, err)

	// Test FindAll with pagination
	results, total, err := repo.FindAll(1, 2)
	require.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Len(t, results, 2)
	// Should be sorted by year DESC, ticker ASC
	assert.Equal(t, 2024, results[0].Year)
}

func TestInMemoryRepository_FindAll_Pagination(t *testing.T) {
	repo := NewInMemoryRepository()

	// Insert test data
	records := []DAXRecord{
		{Company: "A", Ticker: "AAA", Metric: "Revenue", Year: 2024, Value: floatPtr(100.0), Currency: "EUR"},
		{Company: "B", Ticker: "BBB", Metric: "Profit", Year: 2023, Value: floatPtr(50.0), Currency: "EUR"},
		{Company: "C", Ticker: "CCC", Metric: "EBIT", Year: 2024, Value: floatPtr(75.0), Currency: "EUR"},
	}

	err := repo.BulkUpsert(records)
	require.NoError(t, err)

	// Test second page
	results, total, err := repo.FindAll(2, 2)
	require.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Len(t, results, 1)
}

func TestInMemoryRepository_FindAll_OffsetTooLarge(t *testing.T) {
	repo := NewInMemoryRepository()

	records := []DAXRecord{
		{Company: "A", Ticker: "AAA", Metric: "Revenue", Year: 2024, Value: floatPtr(100.0), Currency: "EUR"},
	}

	err := repo.BulkUpsert(records)
	require.NoError(t, err)

	// Request page beyond available data
	results, total, err := repo.FindAll(10, 10)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Empty(t, results)
}

func TestInMemoryRepository_FindByFilters(t *testing.T) {
	repo := NewInMemoryRepository()

	// Insert test data
	records := []DAXRecord{
		{Company: "A", Ticker: "AAA", Metric: "Revenue", Year: 2024, Value: floatPtr(100.0), Currency: "EUR"},
		{Company: "A", Ticker: "AAA", Metric: "Profit", Year: 2024, Value: floatPtr(50.0), Currency: "EUR"},
		{Company: "B", Ticker: "BBB", Metric: "Revenue", Year: 2023, Value: floatPtr(75.0), Currency: "EUR"},
	}

	err := repo.BulkUpsert(records)
	require.NoError(t, err)

	t.Run("filter by ticker", func(t *testing.T) {
		results, total, err := repo.FindByFilters("AAA", nil, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, 2, total)
		assert.Len(t, results, 2)
		for _, r := range results {
			assert.Equal(t, "AAA", r.Ticker)
		}
	})

	t.Run("filter by year", func(t *testing.T) {
		year := 2024
		results, total, err := repo.FindByFilters("", &year, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, 2, total)
		assert.Len(t, results, 2)
		for _, r := range results {
			assert.Equal(t, 2024, r.Year)
		}
	})

	t.Run("filter by ticker and year", func(t *testing.T) {
		year := 2024
		results, total, err := repo.FindByFilters("AAA", &year, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, 2, total)
		assert.Len(t, results, 2)
	})

	t.Run("filter with no matches", func(t *testing.T) {
		results, total, err := repo.FindByFilters("XXX", nil, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, 0, total)
		assert.Empty(t, results)
	})
}

func TestInMemoryRepository_FindByFilters_Pagination(t *testing.T) {
	repo := NewInMemoryRepository()

	records := []DAXRecord{
		{Company: "A", Ticker: "AAA", Metric: "Revenue", Year: 2024, Value: floatPtr(100.0), Currency: "EUR"},
		{Company: "A", Ticker: "AAA", Metric: "Profit", Year: 2023, Value: floatPtr(50.0), Currency: "EUR"},
		{Company: "A", Ticker: "AAA", Metric: "EBIT", Year: 2022, Value: floatPtr(75.0), Currency: "EUR"},
	}

	err := repo.BulkUpsert(records)
	require.NoError(t, err)

	// Test pagination with offset beyond available data
	results, total, err := repo.FindByFilters("AAA", nil, 10, 10)
	require.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Empty(t, results)
}

func TestInMemoryRepository_GetMetrics(t *testing.T) {
	repo := NewInMemoryRepository()

	records := []DAXRecord{
		{Company: "A", Ticker: "AAA", Metric: "Revenue", Year: 2024, Value: floatPtr(100.0), Currency: "EUR"},
		{Company: "A", Ticker: "AAA", Metric: "Profit", Year: 2024, Value: floatPtr(50.0), Currency: "EUR"},
		{Company: "A", Ticker: "AAA", Metric: "EBIT", Year: 2023, Value: floatPtr(75.0), Currency: "EUR"},
		{Company: "B", Ticker: "BBB", Metric: "Revenue", Year: 2024, Value: floatPtr(200.0), Currency: "EUR"},
	}

	err := repo.BulkUpsert(records)
	require.NoError(t, err)

	metrics, err := repo.GetMetrics("AAA")
	require.NoError(t, err)
	assert.Len(t, metrics, 3)
	assert.Contains(t, metrics, "Revenue")
	assert.Contains(t, metrics, "Profit")
	assert.Contains(t, metrics, "EBIT")
}

func TestInMemoryRepository_GetMetrics_NoMatches(t *testing.T) {
	repo := NewInMemoryRepository()

	records := []DAXRecord{
		{Company: "A", Ticker: "AAA", Metric: "Revenue", Year: 2024, Value: floatPtr(100.0), Currency: "EUR"},
	}

	err := repo.BulkUpsert(records)
	require.NoError(t, err)

	metrics, err := repo.GetMetrics("XXX")
	require.NoError(t, err)
	assert.Empty(t, metrics)
}

func TestInMemoryRepository_DeleteAll(t *testing.T) {
	repo := NewInMemoryRepository()

	records := []DAXRecord{
		{Company: "A", Ticker: "AAA", Metric: "Revenue", Year: 2024, Value: floatPtr(100.0), Currency: "EUR"},
		{Company: "B", Ticker: "BBB", Metric: "Profit", Year: 2024, Value: floatPtr(50.0), Currency: "EUR"},
	}

	err := repo.BulkUpsert(records)
	require.NoError(t, err)

	// Verify records exist
	count, err := repo.Count()
	require.NoError(t, err)
	assert.Equal(t, 2, count)

	// Delete all
	err = repo.DeleteAll()
	require.NoError(t, err)

	// Verify all records deleted
	count, err = repo.Count()
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestInMemoryRepository_Count(t *testing.T) {
	repo := NewInMemoryRepository()

	count, err := repo.Count()
	require.NoError(t, err)
	assert.Equal(t, 0, count)

	records := []DAXRecord{
		{Company: "A", Ticker: "AAA", Metric: "Revenue", Year: 2024, Value: floatPtr(100.0), Currency: "EUR"},
		{Company: "B", Ticker: "BBB", Metric: "Profit", Year: 2024, Value: floatPtr(50.0), Currency: "EUR"},
		{Company: "C", Ticker: "CCC", Metric: "EBIT", Year: 2024, Value: floatPtr(75.0), Currency: "EUR"},
	}

	err = repo.BulkUpsert(records)
	require.NoError(t, err)

	count, err = repo.Count()
	require.NoError(t, err)
	assert.Equal(t, 3, count)
}

func TestDAXRecord_TableName(t *testing.T) {
	record := DAXRecord{}
	assert.Equal(t, "dax", record.TableName())
}
