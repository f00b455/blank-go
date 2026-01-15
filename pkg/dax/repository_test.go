package dax

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewInMemoryRepository(t *testing.T) {
	repo := NewInMemoryRepository()

	assert.NotNil(t, repo)
	assert.NotNil(t, repo.records)
	assert.Equal(t, 0, len(repo.records))
}

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

	// Verify record was stored
	count, _ := repo.Count()
	assert.Equal(t, 1, count)
}

func TestInMemoryRepository_Create_WithExistingID(t *testing.T) {
	repo := NewInMemoryRepository()
	id := uuid.New()
	value := 1000.0

	record := &DAXRecord{
		ID:         id,
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
	assert.Equal(t, id, record.ID)
}

func TestInMemoryRepository_BulkUpsert_Insert(t *testing.T) {
	repo := NewInMemoryRepository()
	value1 := 1000.0
	value2 := 2000.0

	records := []DAXRecord{
		{
			Company:    "Company 1",
			Ticker:     "TST1",
			ReportType: "income",
			Metric:     "Revenue",
			Year:       2025,
			Value:      &value1,
			Currency:   "EUR",
		},
		{
			Company:    "Company 2",
			Ticker:     "TST2",
			ReportType: "income",
			Metric:     "Revenue",
			Year:       2025,
			Value:      &value2,
			Currency:   "EUR",
		},
	}

	err := repo.BulkUpsert(records)

	assert.NoError(t, err)
	count, _ := repo.Count()
	assert.Equal(t, 2, count)
}

func TestInMemoryRepository_BulkUpsert_Update(t *testing.T) {
	repo := NewInMemoryRepository()
	value1 := 1000.0
	value2 := 2000.0

	// Insert initial record
	records := []DAXRecord{
		{
			Company:    "Test Company",
			Ticker:     "TST",
			ReportType: "income",
			Metric:     "Revenue",
			Year:       2025,
			Value:      &value1,
			Currency:   "EUR",
		},
	}
	_ = repo.BulkUpsert(records)

	// Upsert with same unique key but different value
	updatedRecords := []DAXRecord{
		{
			Company:    "Test Company",
			Ticker:     "TST",
			ReportType: "income",
			Metric:     "Revenue",
			Year:       2025,
			Value:      &value2,
			Currency:   "EUR",
		},
	}
	err := repo.BulkUpsert(updatedRecords)

	assert.NoError(t, err)
	count, _ := repo.Count()
	assert.Equal(t, 1, count) // Should still be 1 record (upserted)

	// Verify the value was updated
	allRecords, _, _ := repo.FindAll(1, 10)
	assert.Equal(t, value2, *allRecords[0].Value)
}

func TestInMemoryRepository_BulkUpsert_Empty(t *testing.T) {
	repo := NewInMemoryRepository()

	err := repo.BulkUpsert([]DAXRecord{})

	assert.NoError(t, err)
	count, _ := repo.Count()
	assert.Equal(t, 0, count)
}

func TestInMemoryRepository_FindAll(t *testing.T) {
	repo := NewInMemoryRepository()
	value := 1000.0

	// Insert test records
	records := []DAXRecord{
		{Company: "A", Ticker: "AAA", Metric: "M1", Year: 2025, Value: &value, Currency: "EUR", ReportType: "income"},
		{Company: "B", Ticker: "BBB", Metric: "M2", Year: 2024, Value: &value, Currency: "EUR", ReportType: "income"},
		{Company: "C", Ticker: "CCC", Metric: "M3", Year: 2023, Value: &value, Currency: "EUR", ReportType: "income"},
	}
	_ = repo.BulkUpsert(records)

	// Test pagination
	results, total, err := repo.FindAll(1, 2)

	assert.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Equal(t, 2, len(results))
	// Should be sorted by year DESC
	assert.Equal(t, 2025, results[0].Year)
}

func TestInMemoryRepository_FindAll_Sorting(t *testing.T) {
	repo := NewInMemoryRepository()
	value := 1000.0

	// Insert records in random order
	records := []DAXRecord{
		{Company: "C", Ticker: "CCC", Metric: "M1", Year: 2023, Value: &value, Currency: "EUR", ReportType: "income"},
		{Company: "A", Ticker: "AAA", Metric: "M2", Year: 2025, Value: &value, Currency: "EUR", ReportType: "income"},
		{Company: "B", Ticker: "BBB", Metric: "M3", Year: 2024, Value: &value, Currency: "EUR", ReportType: "income"},
	}
	_ = repo.BulkUpsert(records)

	results, _, err := repo.FindAll(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 3, len(results))
	// Verify sorting: year DESC, ticker ASC, metric ASC
	assert.Equal(t, 2025, results[0].Year)
	assert.Equal(t, 2024, results[1].Year)
	assert.Equal(t, 2023, results[2].Year)
}

func TestInMemoryRepository_FindAll_EmptyRepository(t *testing.T) {
	repo := NewInMemoryRepository()

	results, total, err := repo.FindAll(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Equal(t, 0, len(results))
}

func TestInMemoryRepository_FindAll_PaginationOutOfBounds(t *testing.T) {
	repo := NewInMemoryRepository()
	value := 1000.0

	records := []DAXRecord{
		{Company: "A", Ticker: "AAA", Metric: "M1", Year: 2025, Value: &value, Currency: "EUR", ReportType: "income"},
	}
	_ = repo.BulkUpsert(records)

	// Request page beyond available data
	results, total, err := repo.FindAll(5, 10)

	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Equal(t, 0, len(results))
}

func TestInMemoryRepository_FindByFilters_Ticker(t *testing.T) {
	repo := NewInMemoryRepository()
	value := 1000.0

	records := []DAXRecord{
		{Company: "A", Ticker: "AAA", Metric: "M1", Year: 2025, Value: &value, Currency: "EUR", ReportType: "income"},
		{Company: "B", Ticker: "BBB", Metric: "M2", Year: 2025, Value: &value, Currency: "EUR", ReportType: "income"},
		{Company: "C", Ticker: "AAA", Metric: "M3", Year: 2024, Value: &value, Currency: "EUR", ReportType: "income"},
	}
	_ = repo.BulkUpsert(records)

	results, total, err := repo.FindByFilters("AAA", nil, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Equal(t, 2, len(results))
	for _, r := range results {
		assert.Equal(t, "AAA", r.Ticker)
	}
}

func TestInMemoryRepository_FindByFilters_Year(t *testing.T) {
	repo := NewInMemoryRepository()
	value := 1000.0
	year := 2025

	records := []DAXRecord{
		{Company: "A", Ticker: "AAA", Metric: "M1", Year: 2025, Value: &value, Currency: "EUR", ReportType: "income"},
		{Company: "B", Ticker: "BBB", Metric: "M2", Year: 2024, Value: &value, Currency: "EUR", ReportType: "income"},
		{Company: "C", Ticker: "CCC", Metric: "M3", Year: 2025, Value: &value, Currency: "EUR", ReportType: "income"},
	}
	_ = repo.BulkUpsert(records)

	results, total, err := repo.FindByFilters("", &year, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Equal(t, 2, len(results))
	for _, r := range results {
		assert.Equal(t, 2025, r.Year)
	}
}

func TestInMemoryRepository_FindByFilters_TickerAndYear(t *testing.T) {
	repo := NewInMemoryRepository()
	value := 1000.0
	year := 2025

	records := []DAXRecord{
		{Company: "A", Ticker: "AAA", Metric: "M1", Year: 2025, Value: &value, Currency: "EUR", ReportType: "income"},
		{Company: "B", Ticker: "AAA", Metric: "M2", Year: 2024, Value: &value, Currency: "EUR", ReportType: "income"},
		{Company: "C", Ticker: "BBB", Metric: "M3", Year: 2025, Value: &value, Currency: "EUR", ReportType: "income"},
	}
	_ = repo.BulkUpsert(records)

	results, total, err := repo.FindByFilters("AAA", &year, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, "AAA", results[0].Ticker)
	assert.Equal(t, 2025, results[0].Year)
}

func TestInMemoryRepository_FindByFilters_NoMatch(t *testing.T) {
	repo := NewInMemoryRepository()
	value := 1000.0

	records := []DAXRecord{
		{Company: "A", Ticker: "AAA", Metric: "M1", Year: 2025, Value: &value, Currency: "EUR", ReportType: "income"},
	}
	_ = repo.BulkUpsert(records)

	results, total, err := repo.FindByFilters("XYZ", nil, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Equal(t, 0, len(results))
}

func TestInMemoryRepository_GetMetrics(t *testing.T) {
	repo := NewInMemoryRepository()
	value := 1000.0

	records := []DAXRecord{
		{Company: "A", Ticker: "AAA", Metric: "Revenue", Year: 2025, Value: &value, Currency: "EUR", ReportType: "income"},
		{Company: "A", Ticker: "AAA", Metric: "EBITDA", Year: 2025, Value: &value, Currency: "EUR", ReportType: "income"},
		{Company: "A", Ticker: "AAA", Metric: "Revenue", Year: 2024, Value: &value, Currency: "EUR", ReportType: "income"},
		{Company: "B", Ticker: "BBB", Metric: "Profit", Year: 2025, Value: &value, Currency: "EUR", ReportType: "income"},
	}
	_ = repo.BulkUpsert(records)

	metrics, err := repo.GetMetrics("AAA")

	assert.NoError(t, err)
	assert.Equal(t, 2, len(metrics))
	// Should be sorted
	assert.Contains(t, metrics, "Revenue")
	assert.Contains(t, metrics, "EBITDA")
}

func TestInMemoryRepository_GetMetrics_NoRecords(t *testing.T) {
	repo := NewInMemoryRepository()

	metrics, err := repo.GetMetrics("XYZ")

	assert.NoError(t, err)
	assert.Equal(t, 0, len(metrics))
}

func TestInMemoryRepository_DeleteAll(t *testing.T) {
	repo := NewInMemoryRepository()
	value := 1000.0

	records := []DAXRecord{
		{Company: "A", Ticker: "AAA", Metric: "M1", Year: 2025, Value: &value, Currency: "EUR", ReportType: "income"},
		{Company: "B", Ticker: "BBB", Metric: "M2", Year: 2025, Value: &value, Currency: "EUR", ReportType: "income"},
	}
	_ = repo.BulkUpsert(records)

	err := repo.DeleteAll()

	assert.NoError(t, err)
	count, _ := repo.Count()
	assert.Equal(t, 0, count)
}

func TestInMemoryRepository_Count(t *testing.T) {
	repo := NewInMemoryRepository()
	value := 1000.0

	// Empty repository
	count, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, 0, count)

	// After adding records
	records := []DAXRecord{
		{Company: "A", Ticker: "AAA", Metric: "M1", Year: 2025, Value: &value, Currency: "EUR", ReportType: "income"},
		{Company: "B", Ticker: "BBB", Metric: "M2", Year: 2025, Value: &value, Currency: "EUR", ReportType: "income"},
		{Company: "C", Ticker: "CCC", Metric: "M3", Year: 2025, Value: &value, Currency: "EUR", ReportType: "income"},
	}
	_ = repo.BulkUpsert(records)

	count, err = repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, 3, count)
}

func TestNewPostgresRepository(t *testing.T) {
	// Test that NewPostgresRepository doesn't panic with nil
	// (In real usage, this would be called with a real GORM DB)
	repo := NewPostgresRepository(nil)
	assert.NotNil(t, repo)
	assert.Nil(t, repo.db)
}
