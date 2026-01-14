package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/f00b455/blank-go/pkg/dax"
	"github.com/gin-gonic/gin"
)

// DAXHandler handles DAX-related HTTP requests
type DAXHandler struct {
	service *dax.Service
}

// NewDAXHandler creates a new DAX handler
func NewDAXHandler(service *dax.Service) *DAXHandler {
	return &DAXHandler{service: service}
}

// ImportCSV handles CSV file upload and import
func (h *DAXHandler) ImportCSV(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "file is required",
		})
		return
	}

	// Open uploaded file
	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to open file",
		})
		return
	}
	defer func() { _ = openedFile.Close() }()

	// Import CSV
	response, err := h.service.ImportCSV(openedFile)
	if err != nil {
		status := http.StatusInternalServerError
		errMsg := err.Error()
		// Check if it's a validation error (missing fields or invalid data)
		if strings.Contains(errMsg, "missing required fields") ||
			strings.Contains(errMsg, "invalid year") ||
			strings.Contains(errMsg, "invalid data at row") {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{
			"error": errMsg,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetAll retrieves all DAX records with pagination
func (h *DAXHandler) GetAll(c *gin.Context) {
	page := parseIntQuery(c, "page", 1)
	limit := parseIntQuery(c, "limit", 10)

	response, err := h.service.GetAll(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetByFilters retrieves DAX records filtered by ticker and/or year
func (h *DAXHandler) GetByFilters(c *gin.Context) {
	ticker := c.Query("ticker")
	yearStr := c.Query("year")

	var year *int
	if yearStr != "" {
		y, err := strconv.Atoi(yearStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid year parameter",
			})
			return
		}
		year = &y
	}

	page := parseIntQuery(c, "page", 1)
	limit := parseIntQuery(c, "limit", 10)

	response, err := h.service.GetByFilters(ticker, year, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetMetrics retrieves available metrics for a ticker
func (h *DAXHandler) GetMetrics(c *gin.Context) {
	ticker := c.Query("ticker")
	if ticker == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ticker parameter is required",
		})
		return
	}

	response, err := h.service.GetMetrics(ticker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// parseIntQuery parses an integer query parameter with a default value
func parseIntQuery(c *gin.Context, key string, defaultValue int) int {
	valueStr := c.Query(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil || value < 1 {
		return defaultValue
	}

	return value
}
