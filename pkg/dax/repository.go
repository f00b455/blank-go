package dax

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Repository defines the interface for DAX data operations
type Repository interface {
	Create(record *DAXRecord) error
	BulkUpsert(records []DAXRecord) error
	FindAll(page, limit int) ([]DAXRecord, int, error)
	FindByFilters(ticker string, year *int, page, limit int) ([]DAXRecord, int, error)
	GetMetrics(ticker string) ([]string, error)
	DeleteAll() error
	Count() (int, error)
}

// PostgresRepository implements Repository using PostgreSQL
type PostgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Create inserts a single DAX record
func (r *PostgresRepository) Create(record *DAXRecord) error {
	if record.ID == uuid.Nil {
		record.ID = uuid.New()
	}
	return r.db.Create(record).Error
}

// BulkUpsert performs bulk insert with upsert on conflict
func (r *PostgresRepository) BulkUpsert(records []DAXRecord) error {
	if len(records) == 0 {
		return nil
	}

	// Assign UUIDs to records that don't have one
	for i := range records {
		if records[i].ID == uuid.Nil {
			records[i].ID = uuid.New()
		}
	}

	// Use Clauses for UPSERT behavior
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "company"},
			{Name: "ticker"},
			{Name: "metric"},
			{Name: "year"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"report_type",
			"value",
			"currency",
			"updated_at",
		}),
	}).Create(&records).Error
}

// FindAll retrieves all DAX records with pagination
func (r *PostgresRepository) FindAll(page, limit int) ([]DAXRecord, int, error) {
	var records []DAXRecord
	var totalCount int64

	if err := r.db.Model(&DAXRecord{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.Offset(offset).Limit(limit).
		Order("year DESC, ticker ASC, metric ASC").
		Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, int(totalCount), nil
}

// FindByFilters retrieves DAX records with optional filters and pagination
func (r *PostgresRepository) FindByFilters(ticker string, year *int, page, limit int) ([]DAXRecord, int, error) {
	var records []DAXRecord
	var totalCount int64

	query := r.db.Model(&DAXRecord{})

	if ticker != "" {
		query = query.Where("ticker = ?", ticker)
	}

	if year != nil {
		query = query.Where("year = ?", *year)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).
		Order("year DESC, ticker ASC, metric ASC").
		Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, int(totalCount), nil
}

// GetMetrics retrieves unique metrics for a given ticker
func (r *PostgresRepository) GetMetrics(ticker string) ([]string, error) {
	var metrics []string

	if err := r.db.Model(&DAXRecord{}).
		Where("ticker = ?", ticker).
		Distinct("metric").
		Pluck("metric", &metrics).Error; err != nil {
		return nil, err
	}

	return metrics, nil
}

// DeleteAll removes all DAX records (for testing)
func (r *PostgresRepository) DeleteAll() error {
	return r.db.Exec("DELETE FROM dax").Error
}

// Count returns the total number of records
func (r *PostgresRepository) Count() (int, error) {
	var count int64
	if err := r.db.Model(&DAXRecord{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// AutoMigrate creates/updates the database schema
func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&DAXRecord{}); err != nil {
		return fmt.Errorf("failed to migrate DAX schema: %w", err)
	}

	// Create unique composite index for upsert
	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_dax_unique
		ON dax (company, ticker, metric, year)
	`).Error; err != nil {
		return fmt.Errorf("failed to create unique index: %w", err)
	}

	return nil
}
