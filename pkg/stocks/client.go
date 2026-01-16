package stocks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	yahooFinanceURL = "https://query1.finance.yahoo.com/v7/finance/quote"
	defaultTimeout  = 10 * time.Second
)

// StocksClient defines the interface for stock market data retrieval
type StocksClient interface {
	GetQuote(ticker string) (*YahooQuote, error)
	GetQuotes(tickers []string) (map[string]*YahooQuote, error)
}

// Client implements StocksClient using Yahoo Finance API
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new stocks client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		baseURL: yahooFinanceURL,
	}
}

// yahooResponse represents the Yahoo Finance API response
type yahooResponse struct {
	QuoteResponse struct {
		Result []YahooQuote `json:"result"`
		Error  interface{}  `json:"error"`
	} `json:"quoteResponse"`
}

// GetQuote retrieves a single stock quote
func (c *Client) GetQuote(ticker string) (*YahooQuote, error) {
	if ticker == "" {
		return nil, fmt.Errorf("ticker is required")
	}

	quotes, err := c.GetQuotes([]string{ticker})
	if err != nil {
		return nil, err
	}

	quote, ok := quotes[ticker]
	if !ok {
		return nil, fmt.Errorf("ticker not found")
	}

	return quote, nil
}

// GetQuotes retrieves multiple stock quotes
func (c *Client) GetQuotes(tickers []string) (map[string]*YahooQuote, error) {
	if len(tickers) == 0 {
		return nil, fmt.Errorf("at least one ticker is required")
	}

	// Build URL with query parameters
	params := url.Values{}
	params.Add("symbols", strings.Join(tickers, ","))

	requestURL := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())

	// Make HTTP request
	resp, err := c.httpClient.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch quotes: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse JSON response
	var yahooResp yahooResponse
	if err := json.Unmarshal(body, &yahooResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to map
	quotes := make(map[string]*YahooQuote)
	for i := range yahooResp.QuoteResponse.Result {
		quote := &yahooResp.QuoteResponse.Result[i]
		quotes[quote.Symbol] = quote
	}

	return quotes, nil
}
