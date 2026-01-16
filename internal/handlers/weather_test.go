package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/f00b455/blank-go/internal/handlers/mocks"
	"github.com/f00b455/blank-go/pkg/weather"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWeatherHandler(t *testing.T) {
	mockService := mocks.NewMockWeatherService(t)
	handler := NewWeatherHandler(mockService)

	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.service)
}

func TestGetCurrentWeather(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		lat            string
		lon            string
		mockResponse   *weather.WeatherResponse
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful request with valid coordinates",
			lat:  "52.52",
			lon:  "13.41",
			mockResponse: &weather.WeatherResponse{
				Location: weather.Location{
					Latitude:  52.52,
					Longitude: 13.41,
					City:      "Berlin",
				},
				Current: weather.CurrentWeather{
					Temperature: 20.5,
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing latitude parameter",
			lat:            "",
			lon:            "13.41",
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "latitude and longitude are required",
			},
		},
		{
			name:           "missing longitude parameter",
			lat:            "52.52",
			lon:            "",
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "latitude and longitude are required",
			},
		},
		{
			name:           "missing both parameters",
			lat:            "",
			lon:            "",
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "latitude and longitude are required",
			},
		},
		{
			name:           "service returns invalid latitude error",
			lat:            "invalid",
			lon:            "13.41",
			mockResponse:   nil,
			mockError:      errors.New("invalid latitude"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "service returns internal error",
			lat:            "52.52",
			lon:            "13.41",
			mockResponse:   nil,
			mockError:      errors.New("service unavailable"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewMockWeatherService(t)
			handler := NewWeatherHandler(mockService)

			// Setup mock expectations only if we expect the service to be called
			if tt.lat != "" && tt.lon != "" {
				mockService.On("GetCurrentWeatherByCoords", tt.lat, tt.lon).
					Return(tt.mockResponse, tt.mockError)
			}

			// Setup request
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req, _ := http.NewRequest("GET", "/api/v1/weather?lat="+tt.lat+"&lon="+tt.lon, nil)
			c.Request = req

			// Execute
			handler.GetCurrentWeather(c)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedBody["error"], response["error"])
			}
		})
	}
}

func TestGetForecast(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		lat            string
		lon            string
		days           string
		mockResponse   *weather.ForecastResponse
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful request with default days",
			lat:  "52.52",
			lon:  "13.41",
			days: "",
			mockResponse: &weather.ForecastResponse{
				Location: weather.Location{
					Latitude:  52.52,
					Longitude: 13.41,
					City:      "Berlin",
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name: "successful request with custom days",
			lat:  "52.52",
			lon:  "13.41",
			days: "5",
			mockResponse: &weather.ForecastResponse{
				Location: weather.Location{
					Latitude:  52.52,
					Longitude: 13.41,
					City:      "Berlin",
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing latitude parameter",
			lat:            "",
			lon:            "13.41",
			days:           "7",
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "latitude and longitude are required",
			},
		},
		{
			name:           "missing longitude parameter",
			lat:            "52.52",
			lon:            "",
			days:           "7",
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "latitude and longitude are required",
			},
		},
		{
			name:           "invalid days parameter",
			lat:            "52.52",
			lon:            "13.41",
			days:           "invalid",
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid days parameter",
			},
		},
		{
			name:           "service returns error",
			lat:            "52.52",
			lon:            "13.41",
			days:           "7",
			mockResponse:   nil,
			mockError:      errors.New("service error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewMockWeatherService(t)
			handler := NewWeatherHandler(mockService)

			// Setup mock expectations only if we expect the service to be called
			if tt.lat != "" && tt.lon != "" && tt.days != "invalid" {
				expectedDays := 7
				if tt.days != "" {
					expectedDays, _ = strconv.Atoi(tt.days)
				}
				mockService.On("GetForecastByCoords", tt.lat, tt.lon, expectedDays).
					Return(tt.mockResponse, tt.mockError)
			}

			// Setup request
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			url := "/api/v1/weather/forecast?lat=" + tt.lat + "&lon=" + tt.lon
			if tt.days != "" {
				url += "&days=" + tt.days
			}
			req, _ := http.NewRequest("GET", url, nil)
			c.Request = req

			// Execute
			handler.GetForecast(c)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedBody["error"], response["error"])
			}
		})
	}
}

func TestGetWeatherByCity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		city           string
		mockResponse   *weather.WeatherResponse
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful request with valid city",
			city: "Berlin",
			mockResponse: &weather.WeatherResponse{
				Location: weather.Location{
					Latitude:  52.52,
					Longitude: 13.41,
					City:      "Berlin",
				},
				Current: weather.CurrentWeather{
					Temperature: 20.5,
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "city not found",
			city:           "InvalidCity",
			mockResponse:   nil,
			mockError:      errors.New("city not found"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "service returns internal error",
			city:           "Berlin",
			mockResponse:   nil,
			mockError:      errors.New("service unavailable"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewMockWeatherService(t)
			handler := NewWeatherHandler(mockService)

			// Setup mock expectations
			mockService.On("GetWeatherByCity", tt.city).
				Return(tt.mockResponse, tt.mockError)

			// Setup request
			w := httptest.NewRecorder()
			c, router := gin.CreateTestContext(w)
			router.GET("/api/v1/weather/cities/:city", handler.GetWeatherByCity)
			req, _ := http.NewRequest("GET", "/api/v1/weather/cities/"+tt.city, nil)
			c.Request = req
			c.Params = gin.Params{{Key: "city", Value: tt.city}}

			// Execute
			handler.GetWeatherByCity(c)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedBody["error"], response["error"])
			}
		})
	}
}

func TestDetermineStatusCode(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{
			name:           "city not found error",
			err:            errors.New("city not found"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid latitude error",
			err:            errors.New("invalid latitude"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid longitude error",
			err:            errors.New("invalid longitude"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "latitude out of range error",
			err:            errors.New("latitude out of range"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "longitude out of range error",
			err:            errors.New("longitude out of range"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "days must be error",
			err:            errors.New("days must be between 1 and 14"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "city name is required error",
			err:            errors.New("city name is required"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "generic error",
			err:            errors.New("some other error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := determineStatusCode(tt.err)
			assert.Equal(t, tt.expectedStatus, status)
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
			name:     "substring exists at start",
			s:        "hello world",
			substr:   "hello",
			expected: true,
		},
		{
			name:     "substring exists in middle",
			s:        "hello world",
			substr:   "lo wo",
			expected: true,
		},
		{
			name:     "substring exists at end",
			s:        "hello world",
			substr:   "world",
			expected: true,
		},
		{
			name:     "substring does not exist",
			s:        "hello world",
			substr:   "xyz",
			expected: false,
		},
		{
			name:     "exact match",
			s:        "hello",
			substr:   "hello",
			expected: true,
		},
		{
			name:     "empty substring",
			s:        "hello",
			substr:   "",
			expected: true,
		},
		{
			name:     "substring longer than string",
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

func TestContainsAt(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		substr   string
		expected bool
	}{
		{
			name:     "substring at beginning",
			s:        "hello world",
			substr:   "hello",
			expected: true,
		},
		{
			name:     "substring in middle",
			s:        "hello world",
			substr:   "lo wo",
			expected: true,
		},
		{
			name:     "substring at end",
			s:        "hello world",
			substr:   "world",
			expected: true,
		},
		{
			name:     "substring not found",
			s:        "hello world",
			substr:   "xyz",
			expected: false,
		},
		{
			name:     "empty string",
			s:        "",
			substr:   "test",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsAt(tt.s, tt.substr)
			assert.Equal(t, tt.expected, result)
		})
	}
}
