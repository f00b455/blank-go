package dax

import (
	"fmt"
	"sort"
	"sync"

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

// InMemoryRepository implements Repository using in-memory storage
type InMemoryRepository struct {
	mu      sync.RWMutex
	records map[string]*DAXRecord
}

// NewInMemoryRepository creates a new in-memory repository
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		records: make(map[string]*DAXRecord),
	}
}

// Create inserts a single DAX record
func (r *InMemoryRepository) Create(record *DAXRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if record.ID == uuid.Nil {
		record.ID = uuid.New()
	}

	r.records[record.ID.String()] = record
	return nil
}

// BulkUpsert performs bulk insert with upsert on conflict
func (r *InMemoryRepository) BulkUpsert(records []DAXRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range records {
		if records[i].ID == uuid.Nil {
			records[i].ID = uuid.New()
		}

		// Check for existing record with same unique key (company, ticker, metric, year)
		var existingID string
		for id, existing := range r.records {
			if existing.Company == records[i].Company &&
				existing.Ticker == records[i].Ticker &&
				existing.Metric == records[i].Metric &&
				existing.Year == records[i].Year {
				existingID = id
				break
			}
		}

		if existingID != "" {
			// Update existing record
			records[i].ID = r.records[existingID].ID
			records[i].CreatedAt = r.records[existingID].CreatedAt
		}

		r.records[records[i].ID.String()] = &records[i]
	}

	return nil
}

// FindAll retrieves all DAX records with pagination
func (r *InMemoryRepository) FindAll(page, limit int) ([]DAXRecord, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var allRecords []DAXRecord
	for _, record := range r.records {
		allRecords = append(allRecords, *record)
	}

	// Sort by year DESC, ticker ASC, metric ASC
	sort.Slice(allRecords, func(i, j int) bool {
		if allRecords[i].Year != allRecords[j].Year {
			return allRecords[i].Year > allRecords[j].Year
		}
		if allRecords[i].Ticker != allRecords[j].Ticker {
			return allRecords[i].Ticker < allRecords[j].Ticker
		}
		return allRecords[i].Metric < allRecords[j].Metric
	})

	totalCount := len(allRecords)
	offset := (page - 1) * limit

	if offset >= totalCount {
		return []DAXRecord{}, totalCount, nil
	}

	end := offset + limit
	if end > totalCount {
		end = totalCount
	}

	return allRecords[offset:end], totalCount, nil
}

// FindByFilters retrieves DAX records with optional filters and pagination
func (r *InMemoryRepository) FindByFilters(ticker string, year *int, page, limit int) ([]DAXRecord, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []DAXRecord
	for _, record := range r.records {
		if ticker != "" && record.Ticker != ticker {
			continue
		}
		if year != nil && record.Year != *year {
			continue
		}
		filtered = append(filtered, *record)
	}

	// Sort by year DESC, ticker ASC, metric ASC
	sort.Slice(filtered, func(i, j int) bool {
		if filtered[i].Year != filtered[j].Year {
			return filtered[i].Year > filtered[j].Year
		}
		if filtered[i].Ticker != filtered[j].Ticker {
			return filtered[i].Ticker < filtered[j].Ticker
		}
		return filtered[i].Metric < filtered[j].Metric
	})

	totalCount := len(filtered)
	offset := (page - 1) * limit

	if offset >= totalCount {
		return []DAXRecord{}, totalCount, nil
	}

	end := offset + limit
	if end > totalCount {
		end = totalCount
	}

	return filtered[offset:end], totalCount, nil
}

// GetMetrics retrieves unique metrics for a given ticker
func (r *InMemoryRepository) GetMetrics(ticker string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	metricsMap := make(map[string]bool)
	for _, record := range r.records {
		if record.Ticker == ticker {
			metricsMap[record.Metric] = true
		}
	}

	var metrics []string
	for metric := range metricsMap {
		metrics = append(metrics, metric)
	}

	sort.Strings(metrics)
	return metrics, nil
}

// DeleteAll removes all DAX records (for testing)
func (r *InMemoryRepository) DeleteAll() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.records = make(map[string]*DAXRecord)
	return nil
}

// Count returns the total number of records
func (r *InMemoryRepository) Count() (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.records), nil
}
