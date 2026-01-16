package dax

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	return gormDB, mock
}

func TestNewPostgresRepository(t *testing.T) {
	gormDB, _ := setupMockDB(t)
	repo := NewPostgresRepository(gormDB)

	assert.NotNil(t, repo)
	assert.Equal(t, gormDB, repo.db)
}

func TestPostgresRepository_Create(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := NewPostgresRepository(gormDB)

	value := 1000.0
	record := &DAXRecord{
		Company:    "Test AG",
		Ticker:     "TST",
		ReportType: "income",
		Metric:     "EBITDA",
		Year:       2025,
		Value:      &value,
		Currency:   "EUR",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "dax"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(record)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_Create_WithExistingID(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := NewPostgresRepository(gormDB)

	existingID := uuid.New()
	value := 1000.0
	record := &DAXRecord{
		ID:         existingID,
		Company:    "Test AG",
		Ticker:     "TST",
		ReportType: "income",
		Metric:     "EBITDA",
		Year:       2025,
		Value:      &value,
		Currency:   "EUR",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "dax"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(record)

	assert.NoError(t, err)
	assert.Equal(t, existingID, record.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_BulkUpsert_EmptyRecords(t *testing.T) {
	gormDB, _ := setupMockDB(t)
	repo := NewPostgresRepository(gormDB)

	err := repo.BulkUpsert([]DAXRecord{})

	assert.NoError(t, err)
}

func TestPostgresRepository_BulkUpsert(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := NewPostgresRepository(gormDB)

	value := 1000.0
	records := []DAXRecord{
		{Company: "Test AG", Ticker: "TST", ReportType: "income", Metric: "EBITDA", Year: 2025, Value: &value, Currency: "EUR"},
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "dax"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.BulkUpsert(records)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_FindAll(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := NewPostgresRepository(gormDB)

	id := uuid.New()
	now := time.Now()
	value := 1000.0

	// Count query
	mock.ExpectQuery(`SELECT count\(\*\) FROM "dax"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Select query
	mock.ExpectQuery(`SELECT \* FROM "dax"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "company", "ticker", "report_type", "metric", "year", "value", "currency", "created_at", "updated_at"}).
			AddRow(id, "Test AG", "TST", "income", "EBITDA", 2025, value, "EUR", now, now))

	records, total, err := repo.FindAll(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, records, 1)
	assert.Equal(t, "TST", records[0].Ticker)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_FindByFilters_WithTicker(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := NewPostgresRepository(gormDB)

	id := uuid.New()
	now := time.Now()
	value := 1000.0

	// Count query with ticker filter
	mock.ExpectQuery(`SELECT count\(\*\) FROM "dax" WHERE ticker`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Select query with ticker filter (includes ORDER BY, LIMIT, OFFSET)
	mock.ExpectQuery(`SELECT \* FROM "dax" WHERE ticker`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "company", "ticker", "report_type", "metric", "year", "value", "currency", "created_at", "updated_at"}).
			AddRow(id, "Test AG", "TST", "income", "EBITDA", 2025, value, "EUR", now, now))

	records, total, err := repo.FindByFilters("TST", nil, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, records, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_FindByFilters_WithYear(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := NewPostgresRepository(gormDB)

	id := uuid.New()
	now := time.Now()
	value := 1000.0
	year := 2025

	// Count query with year filter
	mock.ExpectQuery(`SELECT count\(\*\) FROM "dax" WHERE year`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Select query with year filter (includes ORDER BY, LIMIT, OFFSET)
	mock.ExpectQuery(`SELECT \* FROM "dax" WHERE year`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "company", "ticker", "report_type", "metric", "year", "value", "currency", "created_at", "updated_at"}).
			AddRow(id, "Test AG", "TST", "income", "EBITDA", 2025, value, "EUR", now, now))

	records, total, err := repo.FindByFilters("", &year, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, records, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_FindByFilters_WithTickerAndYear(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := NewPostgresRepository(gormDB)

	id := uuid.New()
	now := time.Now()
	value := 1000.0
	year := 2025

	// Count query with both filters
	mock.ExpectQuery(`SELECT count\(\*\) FROM "dax" WHERE ticker.*AND year`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Select query with both filters (includes ORDER BY, LIMIT, OFFSET)
	mock.ExpectQuery(`SELECT \* FROM "dax" WHERE ticker.*AND year`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "company", "ticker", "report_type", "metric", "year", "value", "currency", "created_at", "updated_at"}).
			AddRow(id, "Test AG", "TST", "income", "EBITDA", 2025, value, "EUR", now, now))

	records, total, err := repo.FindByFilters("TST", &year, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, records, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetMetrics(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := NewPostgresRepository(gormDB)

	mock.ExpectQuery(`SELECT DISTINCT "metric" FROM "dax" WHERE ticker = \$1`).
		WithArgs("TST").
		WillReturnRows(sqlmock.NewRows([]string{"metric"}).AddRow("EBITDA").AddRow("Revenue"))

	metrics, err := repo.GetMetrics("TST")

	assert.NoError(t, err)
	assert.Len(t, metrics, 2)
	assert.Contains(t, metrics, "EBITDA")
	assert.Contains(t, metrics, "Revenue")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_DeleteAll(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := NewPostgresRepository(gormDB)

	mock.ExpectExec(`DELETE FROM dax`).
		WillReturnResult(sqlmock.NewResult(0, 5))

	err := repo.DeleteAll()

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_Count(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := NewPostgresRepository(gormDB)

	mock.ExpectQuery(`SELECT count\(\*\) FROM "dax"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(42))

	count, err := repo.Count()

	assert.NoError(t, err)
	assert.Equal(t, 42, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}
