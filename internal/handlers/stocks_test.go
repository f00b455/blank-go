package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/f00b455/blank-go/internal/handlers/mocks"
	"github.com/f00b455/blank-go/pkg/stocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestStocksHandler_GetStockSummary(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		ticker         string
		mockSummary    *stocks.StockSummary
		mockError      error
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:   "successful retrieval",
			ticker: "AAPL",
			mockSummary: &stocks.StockSummary{
				Ticker:        "AAPL",
				Name:          "Apple Inc.",
				Date:          "2026-01-16",
				CurrentPrice:  185.50,
				Open:          184.00,
				High:          186.20,
				Low:           183.50,
				Change:        1.50,
				ChangePercent: 0.81,
				Volume:        52000000,
				Currency:      "USD",
				UpdatedAt:     time.Now(),
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "AAPL", body["ticker"])
				assert.Equal(t, "Apple Inc.", body["name"])
				assert.Equal(t, 185.50, body["current_price"])
			},
		},
		{
			name:           "empty ticker",
			ticker:         "",
			mockError:      errors.New("ticker is required"),
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body["error"], "ticker is required")
			},
		},
		{
			name:           "ticker not found",
			ticker:         "INVALID",
			mockError:      errors.New("ticker not found"),
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body["error"], "ticker not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.MockStocksService)
			handler := NewStocksHandler(mockService)

			if tt.mockSummary != nil {
				mockService.On("GetSummary", tt.ticker).Return(tt.mockSummary, tt.mockError)
			} else {
				mockService.On("GetSummary", tt.ticker).Return(nil, tt.mockError)
			}

			router := gin.New()
			router.GET("/stocks/:ticker/summary", handler.GetStockSummary)

			req := httptest.NewRequest(http.MethodGet, "/stocks/"+tt.ticker+"/summary", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestStocksHandler_GetBatchSummary(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		tickers        string
		mockResponse   *stocks.BatchResponse
		mockError      error
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:    "successful batch retrieval",
			tickers: "AAPL,GOOGL,MSFT",
			mockResponse: &stocks.BatchResponse{
				Summaries: []stocks.StockSummary{
					{
						Ticker:       "AAPL",
						Name:         "Apple Inc.",
						CurrentPrice: 185.50,
						Currency:     "USD",
					},
					{
						Ticker:       "GOOGL",
						Name:         "Alphabet Inc.",
						CurrentPrice: 140.20,
						Currency:     "USD",
					},
					{
						Ticker:       "MSFT",
						Name:         "Microsoft Corp.",
						CurrentPrice: 378.90,
						Currency:     "USD",
					},
				},
				Errors: []stocks.BatchError{},
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				summaries := body["summaries"].([]interface{})
				assert.Equal(t, 3, len(summaries))
			},
		},
		{
			name:           "empty tickers",
			tickers:        "",
			mockError:      errors.New("tickers parameter is required"),
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body["error"], "tickers parameter is required")
			},
		},
		{
			name:    "partial success",
			tickers: "AAPL,INVALID",
			mockResponse: &stocks.BatchResponse{
				Summaries: []stocks.StockSummary{
					{
						Ticker:       "AAPL",
						Name:         "Apple Inc.",
						CurrentPrice: 185.50,
						Currency:     "USD",
					},
				},
				Errors: []stocks.BatchError{
					{
						Ticker:  "INVALID",
						Message: "ticker not found",
					},
				},
			},
			expectedStatus: http.StatusMultiStatus,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				summaries := body["summaries"].([]interface{})
				errors := body["errors"].([]interface{})
				assert.Equal(t, 1, len(summaries))
				assert.Equal(t, 1, len(errors))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.MockStocksService)
			handler := NewStocksHandler(mockService)

			if tt.mockResponse != nil {
				mockService.On("GetBatchSummary", tt.tickers).Return(tt.mockResponse, tt.mockError)
			} else {
				mockService.On("GetBatchSummary", tt.tickers).Return(nil, tt.mockError)
			}

			router := gin.New()
			router.GET("/stocks/summary", handler.GetBatchSummary)

			req := httptest.NewRequest(http.MethodGet, "/stocks/summary?tickers="+tt.tickers, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestHandleStockError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		error          error
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "ticker required",
			error:          errors.New("ticker is required"),
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "ticker is required",
		},
		{
			name:           "ticker not found",
			error:          errors.New("ticker not found"),
			expectedStatus: http.StatusNotFound,
			expectedMsg:    "ticker not found",
		},
		{
			name:           "tickers parameter required",
			error:          errors.New("tickers parameter is required"),
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "tickers parameter is required",
		},
		{
			name:           "internal server error",
			error:          errors.New("unknown error"),
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			handleStockError(c, tt.error)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Contains(t, response["error"], tt.expectedMsg)
		})
	}
}
