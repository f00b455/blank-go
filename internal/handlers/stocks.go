package handlers

import (
	"net/http"

	"github.com/f00b455/blank-go/pkg/stocks"
	"github.com/gin-gonic/gin"
)

// StocksService defines the interface for stocks business logic
type StocksService interface {
	GetSummary(ticker string) (*stocks.StockSummary, error)
	GetBatchSummary(tickers string) (*stocks.BatchResponse, error)
}

// StocksHandler handles stock market HTTP requests
type StocksHandler struct {
	service StocksService
}

// NewStocksHandler creates a new stocks handler
func NewStocksHandler(service StocksService) *StocksHandler {
	return &StocksHandler{
		service: service,
	}
}

// GetStockSummary handles GET /api/v1/stocks/:ticker/summary
func (h *StocksHandler) GetStockSummary(c *gin.Context) {
	ticker := c.Param("ticker")

	summary, err := h.service.GetSummary(ticker)
	if err != nil {
		handleStockError(c, err)
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetBatchSummary handles GET /api/v1/stocks/summary?tickers=AAPL,GOOGL,MSFT
func (h *StocksHandler) GetBatchSummary(c *gin.Context) {
	tickers := c.Query("tickers")

	response, err := h.service.GetBatchSummary(tickers)
	if err != nil {
		handleStockError(c, err)
		return
	}

	// Determine status code based on errors
	statusCode := http.StatusOK
	if len(response.Errors) > 0 && len(response.Summaries) > 0 {
		// Partial success
		statusCode = http.StatusMultiStatus
	} else if len(response.Errors) > 0 && len(response.Summaries) == 0 {
		// All failed
		statusCode = http.StatusNotFound
	}

	c.JSON(statusCode, response)
}

// handleStockError maps service errors to HTTP responses
func handleStockError(c *gin.Context, err error) {
	errMsg := err.Error()

	switch errMsg {
	case "ticker is required":
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
	case "ticker not found":
		c.JSON(http.StatusNotFound, gin.H{"error": errMsg})
	case "tickers parameter is required":
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
	case "at least one valid ticker is required":
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
