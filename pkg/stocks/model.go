package stocks

import "time"

// StockSummary represents daily stock market summary
type StockSummary struct {
	Ticker        string    `json:"ticker"`
	Name          string    `json:"name"`
	Date          string    `json:"date"`
	CurrentPrice  float64   `json:"current_price"`
	Open          float64   `json:"open"`
	High          float64   `json:"high"`
	Low           float64   `json:"low"`
	Change        float64   `json:"change"`
	ChangePercent float64   `json:"change_percent"`
	Volume        int64     `json:"volume"`
	Currency      string    `json:"currency"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// BatchResponse represents a batch of stock summaries
type BatchResponse struct {
	Summaries []StockSummary `json:"summaries"`
	Errors    []BatchError   `json:"errors,omitempty"`
}

// BatchError represents an error for a specific ticker in batch request
type BatchError struct {
	Ticker  string `json:"ticker"`
	Message string `json:"message"`
}

// YahooQuote represents the raw quote data from Yahoo Finance
type YahooQuote struct {
	Symbol             string  `json:"symbol"`
	ShortName          string  `json:"shortName"`
	RegularMarketPrice float64 `json:"regularMarketPrice"`
	RegularMarketOpen  float64 `json:"regularMarketOpen"`
	RegularMarketHigh  float64 `json:"regularMarketDayHigh"`
	RegularMarketLow   float64 `json:"regularMarketDayLow"`
	RegularMarketVolume int64   `json:"regularMarketVolume"`
	Currency           string  `json:"currency"`
	RegularMarketChange float64 `json:"regularMarketChange"`
	RegularMarketChangePercent float64 `json:"regularMarketChangePercent"`
}
