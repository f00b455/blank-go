package stocks

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	defaultCacheTTL = 5 * time.Minute
)

// Service provides stock market business logic with caching
type Service struct {
	client    StocksClient
	cache     *stockCache
	cacheTTL  time.Duration
}

// stockCache implements a simple in-memory cache with TTL
type stockCache struct {
	mu    sync.RWMutex
	items map[string]*cacheItem
}

type cacheItem struct {
	summary   *StockSummary
	expiresAt time.Time
}

// NewService creates a new stocks service with caching
func NewService(client StocksClient) *Service {
	return &Service{
		client:   client,
		cache:    newStockCache(),
		cacheTTL: defaultCacheTTL,
	}
}

// newStockCache creates a new cache instance
func newStockCache() *stockCache {
	return &stockCache{
		items: make(map[string]*cacheItem),
	}
}

// get retrieves an item from cache if not expired
func (c *stockCache) get(key string) (*StockSummary, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.expiresAt) {
		return nil, false
	}

	return item.summary, true
}

// set stores an item in cache with TTL
func (c *stockCache) set(key string, summary *StockSummary, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &cacheItem{
		summary:   summary,
		expiresAt: time.Now().Add(ttl),
	}
}

// GetSummary retrieves stock summary for a single ticker
func (s *Service) GetSummary(ticker string) (*StockSummary, error) {
	if ticker == "" {
		return nil, fmt.Errorf("ticker is required")
	}

	ticker = strings.ToUpper(ticker)

	// Check cache first
	if cached, found := s.cache.get(ticker); found {
		return cached, nil
	}

	// Fetch from API
	quote, err := s.client.GetQuote(ticker)
	if err != nil {
		return nil, err
	}

	summary := convertQuoteToSummary(quote)

	// Store in cache
	s.cache.set(ticker, summary, s.cacheTTL)

	return summary, nil
}

// GetBatchSummary retrieves stock summaries for multiple tickers
func (s *Service) GetBatchSummary(tickersStr string) (*BatchResponse, error) {
	if tickersStr == "" {
		return nil, fmt.Errorf("tickers parameter is required")
	}

	// Parse and normalize tickers
	tickers := parseTickers(tickersStr)
	if len(tickers) == 0 {
		return nil, fmt.Errorf("at least one valid ticker is required")
	}

	response := &BatchResponse{
		Summaries: make([]StockSummary, 0, len(tickers)),
		Errors:    make([]BatchError, 0),
	}

	// Check cache for each ticker
	uncachedTickers := make([]string, 0, len(tickers))
	for _, ticker := range tickers {
		if cached, found := s.cache.get(ticker); found {
			response.Summaries = append(response.Summaries, *cached)
		} else {
			uncachedTickers = append(uncachedTickers, ticker)
		}
	}

	// Fetch uncached tickers from API
	if len(uncachedTickers) > 0 {
		quotes, err := s.client.GetQuotes(uncachedTickers)
		if err != nil {
			return nil, err
		}

		// Process results
		for _, ticker := range uncachedTickers {
			quote, found := quotes[ticker]
			if !found {
				response.Errors = append(response.Errors, BatchError{
					Ticker:  ticker,
					Message: "ticker not found",
				})
				continue
			}

			summary := convertQuoteToSummary(quote)
			response.Summaries = append(response.Summaries, *summary)

			// Cache the result
			s.cache.set(ticker, summary, s.cacheTTL)
		}
	}

	return response, nil
}

// parseTickers splits and normalizes ticker string
func parseTickers(tickersStr string) []string {
	parts := strings.Split(tickersStr, ",")
	tickers := make([]string, 0, len(parts))

	for _, part := range parts {
		ticker := strings.TrimSpace(strings.ToUpper(part))
		if ticker != "" {
			tickers = append(tickers, ticker)
		}
	}

	return tickers
}

// convertQuoteToSummary converts Yahoo quote to stock summary
func convertQuoteToSummary(quote *YahooQuote) *StockSummary {
	now := time.Now()
	return &StockSummary{
		Ticker:        quote.Symbol,
		Name:          quote.ShortName,
		Date:          now.Format("2006-01-02"),
		CurrentPrice:  quote.RegularMarketPrice,
		Open:          quote.RegularMarketOpen,
		High:          quote.RegularMarketHigh,
		Low:           quote.RegularMarketLow,
		Change:        quote.RegularMarketChange,
		ChangePercent: quote.RegularMarketChangePercent,
		Volume:        quote.RegularMarketVolume,
		Currency:      quote.Currency,
		UpdatedAt:     now,
	}
}
