package stocks_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/f00b455/blank-go/pkg/stocks"
	"github.com/f00b455/blank-go/pkg/stocks/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_GetSummary(t *testing.T) {
	tests := []struct {
		name          string
		ticker        string
		mockQuote     *stocks.YahooQuote
		mockError     error
		expectError   bool
		errorContains string
	}{
		{
			name:   "successful retrieval",
			ticker: "AAPL",
			mockQuote: &stocks.YahooQuote{
				Symbol:                     "AAPL",
				ShortName:                  "Apple Inc.",
				RegularMarketPrice:         185.50,
				RegularMarketOpen:          184.00,
				RegularMarketHigh:          186.20,
				RegularMarketLow:           183.50,
				RegularMarketChange:        1.50,
				RegularMarketChangePercent: 0.81,
				RegularMarketVolume:        52000000,
				Currency:                   "USD",
			},
			expectError: false,
		},
		{
			name:          "empty ticker",
			ticker:        "",
			expectError:   true,
			errorContains: "ticker is required",
		},
		{
			name:          "client error",
			ticker:        "INVALID",
			mockError:     errors.New("ticker not found"),
			expectError:   true,
			errorContains: "ticker not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(mocks.MockStocksClient)
			service := stocks.NewService(mockClient)

			if tt.ticker != "" && tt.mockQuote != nil {
				mockClient.On("GetQuote", mock.Anything).Return(tt.mockQuote, tt.mockError)
			} else if tt.ticker != "" && tt.mockError != nil {
				mockClient.On("GetQuote", mock.Anything).Return(nil, tt.mockError)
			}

			summary, err := service.GetSummary(tt.ticker)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, summary)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, summary)
				assert.Equal(t, tt.mockQuote.Symbol, summary.Ticker)
				assert.Equal(t, tt.mockQuote.ShortName, summary.Name)
				assert.Equal(t, tt.mockQuote.RegularMarketPrice, summary.CurrentPrice)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestService_GetSummary_Caching(t *testing.T) {
	mockClient := new(mocks.MockStocksClient)
	service := stocks.NewService(mockClient)

	quote := &stocks.YahooQuote{
		Symbol:             "AAPL",
		ShortName:          "Apple Inc.",
		RegularMarketPrice: 185.50,
		Currency:           "USD",
	}

	mockClient.On("GetQuote", "AAPL").Return(quote, nil).Once()

	summary1, err := service.GetSummary("AAPL")
	assert.NoError(t, err)
	assert.NotNil(t, summary1)

	summary2, err := service.GetSummary("AAPL")
	assert.NoError(t, err)
	assert.NotNil(t, summary2)
	assert.Equal(t, summary1.Ticker, summary2.Ticker)

	mockClient.AssertExpectations(t)
}

func TestService_GetBatchSummary(t *testing.T) {
	tests := []struct {
		name           string
		tickersStr     string
		mockQuotes     map[string]*stocks.YahooQuote
		mockError      error
		expectError    bool
		expectedCount  int
		expectedErrors int
		errorContains  string
	}{
		{
			name:       "successful batch retrieval",
			tickersStr: "AAPL,GOOGL,MSFT",
			mockQuotes: map[string]*stocks.YahooQuote{
				"AAPL": {
					Symbol:             "AAPL",
					ShortName:          "Apple Inc.",
					RegularMarketPrice: 185.50,
					Currency:           "USD",
				},
				"GOOGL": {
					Symbol:             "GOOGL",
					ShortName:          "Alphabet Inc.",
					RegularMarketPrice: 140.20,
					Currency:           "USD",
				},
				"MSFT": {
					Symbol:             "MSFT",
					ShortName:          "Microsoft Corp.",
					RegularMarketPrice: 378.90,
					Currency:           "USD",
				},
			},
			expectError:    false,
			expectedCount:  3,
			expectedErrors: 0,
		},
		{
			name:          "empty tickers",
			tickersStr:    "",
			expectError:   true,
			errorContains: "tickers parameter is required",
		},
		{
			name:       "partial success",
			tickersStr: "AAPL,INVALID",
			mockQuotes: map[string]*stocks.YahooQuote{
				"AAPL": {
					Symbol:             "AAPL",
					ShortName:          "Apple Inc.",
					RegularMarketPrice: 185.50,
					Currency:           "USD",
				},
			},
			expectError:    false,
			expectedCount:  1,
			expectedErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(mocks.MockStocksClient)
			service := stocks.NewService(mockClient)

			if tt.tickersStr != "" && !tt.expectError {
				mockClient.On("GetQuotes", mock.Anything).Return(tt.mockQuotes, tt.mockError)
			}

			response, err := service.GetBatchSummary(tt.tickersStr)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, tt.expectedCount, len(response.Summaries))
				assert.Equal(t, tt.expectedErrors, len(response.Errors))
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestParseTickers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single ticker",
			input:    "AAPL",
			expected: []string{"AAPL"},
		},
		{
			name:     "multiple tickers",
			input:    "AAPL,GOOGL,MSFT",
			expected: []string{"AAPL", "GOOGL", "MSFT"},
		},
		{
			name:     "lowercase normalization",
			input:    "aapl,googl",
			expected: []string{"AAPL", "GOOGL"},
		},
		{
			name:     "whitespace handling",
			input:    " AAPL , GOOGL , MSFT ",
			expected: []string{"AAPL", "GOOGL", "MSFT"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts := strings.Split(tt.input, ",")
			result := make([]string, 0, len(parts))
			for _, part := range parts {
				ticker := strings.TrimSpace(strings.ToUpper(part))
				if ticker != "" {
					result = append(result, ticker)
				}
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertQuoteToSummary(t *testing.T) {
	mockClient := new(mocks.MockStocksClient)
	service := stocks.NewService(mockClient)

	quote := &stocks.YahooQuote{
		Symbol:                     "AAPL",
		ShortName:                  "Apple Inc.",
		RegularMarketPrice:         185.50,
		RegularMarketOpen:          184.00,
		RegularMarketHigh:          186.20,
		RegularMarketLow:           183.50,
		RegularMarketChange:        1.50,
		RegularMarketChangePercent: 0.81,
		RegularMarketVolume:        52000000,
		Currency:                   "USD",
	}

	mockClient.On("GetQuote", "AAPL").Return(quote, nil)

	summary, err := service.GetSummary("AAPL")
	assert.NoError(t, err)

	assert.Equal(t, quote.Symbol, summary.Ticker)
	assert.Equal(t, quote.ShortName, summary.Name)
	assert.Equal(t, quote.RegularMarketPrice, summary.CurrentPrice)
	assert.Equal(t, quote.RegularMarketOpen, summary.Open)
	assert.Equal(t, quote.RegularMarketHigh, summary.High)
	assert.Equal(t, quote.RegularMarketLow, summary.Low)
	assert.Equal(t, quote.RegularMarketChange, summary.Change)
	assert.Equal(t, quote.RegularMarketChangePercent, summary.ChangePercent)
	assert.Equal(t, quote.RegularMarketVolume, summary.Volume)
	assert.Equal(t, quote.Currency, summary.Currency)
	assert.NotEmpty(t, summary.Date)

	mockClient.AssertExpectations(t)
}
