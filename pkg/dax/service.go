package dax

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Service provides business logic for DAX operations
type Service struct {
	repo Repository
}

// NewService creates a new DAX service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// ImportCSV imports DAX data from CSV content
func (s *Service) ImportCSV(reader io.Reader) (*ImportResponse, error) {
	csvReader := csv.NewReader(reader)

	// Read header
	header, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Validate header
	if err := validateHeader(header); err != nil {
		return nil, err
	}

	// Parse rows
	records := []DAXRecord{}
	rowNum := 1

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV row %d: %w", rowNum, err)
		}

		record, err := parseCSVRow(row)
		if err != nil {
			return nil, fmt.Errorf("invalid data at row %d: %w", rowNum, err)
		}

		records = append(records, record)
		rowNum++
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no records found in CSV")
	}

	// Bulk insert with upsert
	if err := s.repo.BulkUpsert(records); err != nil {
		return nil, fmt.Errorf("failed to import records: %w", err)
	}

	return &ImportResponse{
		RecordsImported: len(records),
		Message:         fmt.Sprintf("Successfully imported %d records", len(records)),
	}, nil
}

// GetAll retrieves all DAX records with pagination
func (s *Service) GetAll(page, limit int) (*PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	records, total, err := s.repo.FindAll(page, limit)
	if err != nil {
		return nil, err
	}

	totalPages := (total + limit - 1) / limit

	return &PaginatedResponse{
		Data: records,
		Pagination: PaginationMeta{
			Page:       page,
			Limit:      limit,
			TotalCount: total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetByFilters retrieves DAX records filtered by ticker and/or year
func (s *Service) GetByFilters(ticker string, year *int, page, limit int) (*PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	records, total, err := s.repo.FindByFilters(ticker, year, page, limit)
	if err != nil {
		return nil, err
	}

	totalPages := (total + limit - 1) / limit

	return &PaginatedResponse{
		Data: records,
		Pagination: PaginationMeta{
			Page:       page,
			Limit:      limit,
			TotalCount: total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetMetrics retrieves available metrics for a ticker
func (s *Service) GetMetrics(ticker string) (*MetricsResponse, error) {
	if ticker == "" {
		return nil, fmt.Errorf("ticker is required")
	}

	metrics, err := s.repo.GetMetrics(ticker)
	if err != nil {
		return nil, err
	}

	return &MetricsResponse{
		Ticker:  ticker,
		Metrics: metrics,
	}, nil
}

// validateHeader checks if CSV has all required fields
func validateHeader(header []string) error {
	required := map[string]bool{
		"company":     false,
		"ticker":      false,
		"report_type": false,
		"metric":      false,
		"year":        false,
		"value":       false,
		"currency":    false,
	}

	for _, col := range header {
		col = strings.TrimSpace(strings.ToLower(col))
		if _, exists := required[col]; exists {
			required[col] = true
		}
	}

	missing := []string{}
	for field, found := range required {
		if !found {
			missing = append(missing, field)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missing, ", "))
	}

	return nil
}

// parseCSVRow parses a CSV row into a DAXRecord
func parseCSVRow(row []string) (DAXRecord, error) {
	if len(row) < 7 {
		return DAXRecord{}, fmt.Errorf("insufficient columns")
	}

	year, err := strconv.Atoi(strings.TrimSpace(row[4]))
	if err != nil {
		return DAXRecord{}, fmt.Errorf("invalid year: %w", err)
	}

	value, err := strconv.ParseFloat(strings.TrimSpace(row[5]), 64)
	if err != nil {
		return DAXRecord{}, fmt.Errorf("invalid value: %w", err)
	}

	return DAXRecord{
		Company:    strings.TrimSpace(row[0]),
		Ticker:     strings.TrimSpace(row[1]),
		ReportType: strings.TrimSpace(row[2]),
		Metric:     strings.TrimSpace(row[3]),
		Year:       year,
		Value:      &value,
		Currency:   strings.TrimSpace(row[6]),
	}, nil
}
