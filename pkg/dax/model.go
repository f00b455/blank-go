package dax

import (
	"time"

	"github.com/google/uuid"
)

// DAXRecord represents a financial data record for a DAX company
type DAXRecord struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Company    string    `json:"company" gorm:"type:varchar(255);not null"`
	Ticker     string    `json:"ticker" gorm:"type:varchar(10);not null;index:idx_ticker_year"`
	ReportType string    `json:"report_type" gorm:"type:varchar(50)"`
	Metric     string    `json:"metric" gorm:"type:varchar(100);not null"`
	Year       int       `json:"year" gorm:"not null;index:idx_ticker_year"`
	Value      *float64  `json:"value" gorm:"type:decimal(20,2)"`
	Currency   string    `json:"currency" gorm:"type:varchar(3);default:'EUR'"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName sets the table name for GORM
func (DAXRecord) TableName() string {
	return "dax"
}

// CSVRow represents a row from the CSV import file
type CSVRow struct {
	Company    string  `json:"company"`
	Ticker     string  `json:"ticker"`
	ReportType string  `json:"report_type"`
	Metric     string  `json:"metric"`
	Year       int     `json:"year"`
	Value      float64 `json:"value"`
	Currency   string  `json:"currency"`
}

// ImportResponse represents the response from an import operation
type ImportResponse struct {
	RecordsImported int    `json:"records_imported"`
	Message         string `json:"message"`
}

// PaginatedResponse represents a paginated list of DAX records
type PaginatedResponse struct {
	Data       []DAXRecord    `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// PaginationMeta contains pagination metadata
type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalCount int `json:"total_count"`
	TotalPages int `json:"total_pages"`
}

// MetricsResponse contains available metrics for a ticker
type MetricsResponse struct {
	Ticker  string   `json:"ticker"`
	Metrics []string `json:"metrics"`
}
