package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/f00b455/blank-go/pkg/weather"
)

const defaultForecastDays = 7

// WeatherService defines the interface for weather operations
type WeatherService interface {
	GetCurrentWeatherByCoords(lat, lon string) (*weather.WeatherResponse, error)
	GetForecastByCoords(lat, lon string, days int) (*weather.ForecastResponse, error)
	GetWeatherByCity(cityName string) (*weather.WeatherResponse, error)
}

// WeatherHandler handles weather-related HTTP requests
type WeatherHandler struct {
	service WeatherService
}

// NewWeatherHandler creates a new weather handler
func NewWeatherHandler(service WeatherService) *WeatherHandler {
	return &WeatherHandler{
		service: service,
	}
}

// GetCurrentWeather handles GET /api/v1/weather
func (h *WeatherHandler) GetCurrentWeather(c *gin.Context) {
	lat := c.Query("lat")
	lon := c.Query("lon")

	if lat == "" || lon == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "latitude and longitude are required",
		})
		return
	}

	result, err := h.service.GetCurrentWeatherByCoords(lat, lon)
	if err != nil {
		statusCode := determineStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetForecast handles GET /api/v1/weather/forecast
func (h *WeatherHandler) GetForecast(c *gin.Context) {
	lat := c.Query("lat")
	lon := c.Query("lon")
	daysStr := c.DefaultQuery("days", strconv.Itoa(defaultForecastDays))

	if lat == "" || lon == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "latitude and longitude are required",
		})
		return
	}

	days, err := strconv.Atoi(daysStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid days parameter",
		})
		return
	}

	result, err := h.service.GetForecastByCoords(lat, lon, days)
	if err != nil {
		statusCode := determineStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetWeatherByCity handles GET /api/v1/weather/cities/:city
func (h *WeatherHandler) GetWeatherByCity(c *gin.Context) {
	city := c.Param("city")

	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "city name is required",
		})
		return
	}

	result, err := h.service.GetWeatherByCity(city)
	if err != nil {
		statusCode := determineStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// determineStatusCode maps error messages to HTTP status codes
func determineStatusCode(err error) int {
	errMsg := err.Error()

	switch {
	case contains(errMsg, "city not found"):
		return http.StatusNotFound
	case contains(errMsg, "invalid latitude"),
		contains(errMsg, "invalid longitude"),
		contains(errMsg, "latitude out of range"),
		contains(errMsg, "longitude out of range"),
		contains(errMsg, "days must be"),
		contains(errMsg, "city name is required"):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsAt(s, substr))
}

// containsAt checks if substr exists anywhere in s
func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
