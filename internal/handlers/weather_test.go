package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/f00b455/blank-go/internal/handlers/mocks"
	"github.com/f00b455/blank-go/pkg/weather"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestWeatherHandler_GetCurrentWeather_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := mocks.NewMockWeatherService(t)
	handler := NewWeatherHandler(mockService)

	expectedResponse := &weather.WeatherResponse{
		Location: weather.Location{
			Latitude:  52.52,
			Longitude: 13.41,
			City:      "Berlin",
			Timezone:  "Europe/Berlin",
		},
		Current: weather.CurrentWeather{
			Temperature:        20.5,
			Humidity:           65,
			WindSpeed:          5.0,
			WeatherCode:        0,
			WeatherDescription: "Clear sky",
		},
		Units: weather.Units{
			Temperature: "°C",
			WindSpeed:   "km/h",
			Humidity:    "%",
		},
	}

	mockService.On("GetCurrentWeatherByCoords", "52.52", "13.41").
		Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/weather?lat=52.52&lon=13.41", nil)

	handler.GetCurrentWeather(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestWeatherHandler_GetCurrentWeather_MissingLatitude(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := mocks.NewMockWeatherService(t)
	handler := NewWeatherHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/weather?lon=13.41", nil)

	handler.GetCurrentWeather(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "latitude and longitude are required")
}

func TestWeatherHandler_GetCurrentWeather_MissingLongitude(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := mocks.NewMockWeatherService(t)
	handler := NewWeatherHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/weather?lat=52.52", nil)

	handler.GetCurrentWeather(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "latitude and longitude are required")
}

func TestWeatherHandler_GetCurrentWeather_InvalidCoordinates(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := mocks.NewMockWeatherService(t)
	handler := NewWeatherHandler(mockService)

	mockService.On("GetCurrentWeatherByCoords", "999", "13.41").
		Return(nil, errors.New("latitude out of range"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/weather?lat=999&lon=13.41", nil)

	handler.GetCurrentWeather(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "latitude out of range")
	mockService.AssertExpectations(t)
}

func TestWeatherHandler_GetCurrentWeather_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := mocks.NewMockWeatherService(t)
	handler := NewWeatherHandler(mockService)

	mockService.On("GetCurrentWeatherByCoords", "52.52", "13.41").
		Return(nil, errors.New("internal service error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/weather?lat=52.52&lon=13.41", nil)

	handler.GetCurrentWeather(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "internal service error")
	mockService.AssertExpectations(t)
}

func TestWeatherHandler_GetForecast_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := mocks.NewMockWeatherService(t)
	handler := NewWeatherHandler(mockService)

	expectedResponse := &weather.ForecastResponse{
		Location: weather.Location{
			Latitude:  52.52,
			Longitude: 13.41,
			City:      "Berlin",
			Timezone:  "Europe/Berlin",
		},
		Forecast: []weather.ForecastDay{
			{Date: "2025-01-01"},
		},
	}

	mockService.On("GetForecastByCoords", "52.52", "13.41", 7).
		Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/weather/forecast?lat=52.52&lon=13.41&days=7", nil)

	handler.GetForecast(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestWeatherHandler_GetForecast_DefaultDays(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := mocks.NewMockWeatherService(t)
	handler := NewWeatherHandler(mockService)

	expectedResponse := &weather.ForecastResponse{
		Location: weather.Location{
			Latitude:  52.52,
			Longitude: 13.41,
			City:      "Berlin",
			Timezone:  "Europe/Berlin",
		},
		Forecast: []weather.ForecastDay{},
	}

	mockService.On("GetForecastByCoords", "52.52", "13.41", 7).
		Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/weather/forecast?lat=52.52&lon=13.41", nil)

	handler.GetForecast(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestWeatherHandler_GetForecast_MissingLatitude(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := mocks.NewMockWeatherService(t)
	handler := NewWeatherHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/weather/forecast?lon=13.41", nil)

	handler.GetForecast(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "latitude and longitude are required")
}

func TestWeatherHandler_GetForecast_InvalidDays(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := mocks.NewMockWeatherService(t)
	handler := NewWeatherHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/weather/forecast?lat=52.52&lon=13.41&days=invalid", nil)

	handler.GetForecast(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid days parameter")
}

func TestWeatherHandler_GetForecast_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := mocks.NewMockWeatherService(t)
	handler := NewWeatherHandler(mockService)

	mockService.On("GetForecastByCoords", "52.52", "13.41", 5).
		Return(nil, errors.New("days must be between 1 and 16"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/weather/forecast?lat=52.52&lon=13.41&days=5", nil)

	handler.GetForecast(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "days must be")
	mockService.AssertExpectations(t)
}

func TestWeatherHandler_GetWeatherByCity_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := mocks.NewMockWeatherService(t)
	handler := NewWeatherHandler(mockService)

	expectedResponse := &weather.WeatherResponse{
		Location: weather.Location{
			Latitude:  52.52,
			Longitude: 13.41,
			City:      "Berlin",
			Timezone:  "Europe/Berlin",
		},
		Current: weather.CurrentWeather{
			Temperature:        15.0,
			Humidity:           70,
			WindSpeed:          10.0,
			WeatherCode:        2,
			WeatherDescription: "Partly cloudy",
		},
		Units: weather.Units{
			Temperature: "°C",
			WindSpeed:   "km/h",
			Humidity:    "%",
		},
	}

	mockService.On("GetWeatherByCity", "Berlin").
		Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "city", Value: "Berlin"}}
	c.Request = httptest.NewRequest("GET", "/api/v1/weather/cities/Berlin", nil)

	handler.GetWeatherByCity(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestWeatherHandler_GetWeatherByCity_EmptyCity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := mocks.NewMockWeatherService(t)
	handler := NewWeatherHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "city", Value: ""}}
	c.Request = httptest.NewRequest("GET", "/api/v1/weather/cities/", nil)

	handler.GetWeatherByCity(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "city name is required")
}

func TestWeatherHandler_GetWeatherByCity_CityNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := mocks.NewMockWeatherService(t)
	handler := NewWeatherHandler(mockService)

	mockService.On("GetWeatherByCity", "UnknownCity").
		Return(nil, errors.New("city not found"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "city", Value: "UnknownCity"}}
	c.Request = httptest.NewRequest("GET", "/api/v1/weather/cities/UnknownCity", nil)

	handler.GetWeatherByCity(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "city not found")
	mockService.AssertExpectations(t)
}

func TestDetermineStatusCode(t *testing.T) {
	tests := []struct {
		name     string
		errMsg   string
		expected int
	}{
		{
			name:     "city not found",
			errMsg:   "city not found",
			expected: http.StatusNotFound,
		},
		{
			name:     "invalid latitude",
			errMsg:   "invalid latitude",
			expected: http.StatusBadRequest,
		},
		{
			name:     "invalid longitude",
			errMsg:   "invalid longitude",
			expected: http.StatusBadRequest,
		},
		{
			name:     "latitude out of range",
			errMsg:   "latitude out of range",
			expected: http.StatusBadRequest,
		},
		{
			name:     "longitude out of range",
			errMsg:   "longitude out of range",
			expected: http.StatusBadRequest,
		},
		{
			name:     "days validation error",
			errMsg:   "days must be between 1 and 16",
			expected: http.StatusBadRequest,
		},
		{
			name:     "generic error",
			errMsg:   "some unknown error",
			expected: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errors.New(tt.errMsg)
			statusCode := determineStatusCode(err)
			assert.Equal(t, tt.expected, statusCode)
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		substr   string
		expected bool
	}{
		{
			name:     "exact match",
			s:        "hello",
			substr:   "hello",
			expected: true,
		},
		{
			name:     "substring at start",
			s:        "hello world",
			substr:   "hello",
			expected: true,
		},
		{
			name:     "substring at end",
			s:        "hello world",
			substr:   "world",
			expected: true,
		},
		{
			name:     "substring in middle",
			s:        "hello world",
			substr:   "lo wo",
			expected: true,
		},
		{
			name:     "not found",
			s:        "hello world",
			substr:   "xyz",
			expected: false,
		},
		{
			name:     "empty substring",
			s:        "hello",
			substr:   "",
			expected: true,
		},
		{
			name:     "substr longer than s",
			s:        "hi",
			substr:   "hello",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.s, tt.substr)
			assert.Equal(t, tt.expected, result)
		})
	}
}
