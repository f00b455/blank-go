package weather

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockHTTPClient is a mock implementation of HTTPClient
type MockHTTPClient struct {
	GetFunc func(url string) (*http.Response, error)
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	if m.GetFunc != nil {
		return m.GetFunc(url)
	}
	return nil, errors.New("GetFunc not implemented")
}

func TestNewClient(t *testing.T) {
	client := NewClient()
	assert.NotNil(t, client)
	assert.NotNil(t, client.httpClient)
	assert.Equal(t, defaultRequestTimeout, client.timeout)
}

func TestNewClientWithHTTP(t *testing.T) {
	mockClient := &MockHTTPClient{}
	client := NewClientWithHTTP(mockClient)
	assert.NotNil(t, client)
	assert.Equal(t, mockClient, client.httpClient)
	assert.Equal(t, defaultRequestTimeout, client.timeout)
}

func TestGetCurrentWeather(t *testing.T) {
	tests := []struct {
		name          string
		lat           float64
		lon           float64
		mockResponse  string
		mockStatus    int
		mockError     error
		expectError   bool
		errorContains string
	}{
		{
			name: "successful weather fetch",
			lat:  52.52,
			lon:  13.405,
			mockResponse: `{
				"latitude": 52.52,
				"longitude": 13.405,
				"timezone": "Europe/Berlin",
				"current": {
					"temperature_2m": 15.5,
					"relative_humidity_2m": 65,
					"wind_speed_10m": 12.3,
					"weather_code": 0
				}
			}`,
			mockStatus:  http.StatusOK,
			expectError: false,
		},
		{
			name:          "http client error",
			lat:           52.52,
			lon:           13.405,
			mockError:     errors.New("network error"),
			expectError:   true,
			errorContains: "failed to fetch weather data",
		},
		{
			name:          "API returns error status",
			lat:           52.52,
			lon:           13.405,
			mockResponse:  `{"error": "invalid coordinates"}`,
			mockStatus:    http.StatusBadRequest,
			expectError:   true,
			errorContains: "API returned status 400",
		},
		{
			name:          "invalid JSON response",
			lat:           52.52,
			lon:           13.405,
			mockResponse:  `{invalid json}`,
			mockStatus:    http.StatusOK,
			expectError:   true,
			errorContains: "failed to decode response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockHTTPClient{
				GetFunc: func(url string) (*http.Response, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return &http.Response{
						StatusCode: tt.mockStatus,
						Body:       io.NopCloser(strings.NewReader(tt.mockResponse)),
					}, nil
				},
			}

			client := NewClientWithHTTP(mockClient)
			result, err := client.GetCurrentWeather(tt.lat, tt.lon)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.lat, result.Location.Latitude)
				assert.Equal(t, tt.lon, result.Location.Longitude)
				assert.NotEmpty(t, result.Current.WeatherDescription)
			}
		})
	}
}

func TestGetForecast(t *testing.T) {
	tests := []struct {
		name          string
		lat           float64
		lon           float64
		days          int
		mockResponse  string
		mockStatus    int
		mockError     error
		expectError   bool
		errorContains string
	}{
		{
			name: "successful forecast fetch",
			lat:  52.52,
			lon:  13.405,
			days: 3,
			mockResponse: `{
				"latitude": 52.52,
				"longitude": 13.405,
				"timezone": "Europe/Berlin",
				"daily": {
					"time": ["2026-01-16", "2026-01-17", "2026-01-18"],
					"temperature_2m_max": [10.5, 12.3, 11.8],
					"temperature_2m_min": [5.2, 6.1, 5.8],
					"precipitation_probability_max": [20, 40, 10],
					"weather_code": [0, 1, 2]
				}
			}`,
			mockStatus:  http.StatusOK,
			expectError: false,
		},
		{
			name:          "http client error",
			lat:           52.52,
			lon:           13.405,
			days:          3,
			mockError:     errors.New("network error"),
			expectError:   true,
			errorContains: "failed to fetch forecast data",
		},
		{
			name:          "API returns error status",
			lat:           52.52,
			lon:           13.405,
			days:          3,
			mockResponse:  `{"error": "invalid parameters"}`,
			mockStatus:    http.StatusBadRequest,
			expectError:   true,
			errorContains: "API returned status 400",
		},
		{
			name:          "invalid JSON response",
			lat:           52.52,
			lon:           13.405,
			days:          3,
			mockResponse:  `{invalid}`,
			mockStatus:    http.StatusOK,
			expectError:   true,
			errorContains: "failed to decode response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockHTTPClient{
				GetFunc: func(url string) (*http.Response, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return &http.Response{
						StatusCode: tt.mockStatus,
						Body:       io.NopCloser(strings.NewReader(tt.mockResponse)),
					}, nil
				},
			}

			client := NewClientWithHTTP(mockClient)
			result, err := client.GetForecast(tt.lat, tt.lon, tt.days)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.lat, result.Location.Latitude)
				assert.Equal(t, tt.lon, result.Location.Longitude)
				assert.Len(t, result.Forecast, tt.days)
			}
		})
	}
}

func TestGeocodeCity(t *testing.T) {
	tests := []struct {
		name          string
		cityName      string
		mockResponse  string
		mockStatus    int
		mockError     error
		expectError   bool
		errorContains string
	}{
		{
			name:     "successful geocoding",
			cityName: "Berlin",
			mockResponse: `{
				"results": [{
					"name": "Berlin",
					"latitude": 52.52,
					"longitude": 13.405,
					"timezone": "Europe/Berlin"
				}]
			}`,
			mockStatus:  http.StatusOK,
			expectError: false,
		},
		{
			name:          "http client error",
			cityName:      "Berlin",
			mockError:     errors.New("network error"),
			expectError:   true,
			errorContains: "failed to geocode city",
		},
		{
			name:          "API returns error status",
			cityName:      "Berlin",
			mockResponse:  `{"error": "invalid request"}`,
			mockStatus:    http.StatusBadRequest,
			expectError:   true,
			errorContains: "API returned status 400",
		},
		{
			name:          "city not found",
			cityName:      "NonExistentCity",
			mockResponse:  `{"results": []}`,
			mockStatus:    http.StatusOK,
			expectError:   true,
			errorContains: "city not found",
		},
		{
			name:          "invalid JSON response",
			cityName:      "Berlin",
			mockResponse:  `{invalid}`,
			mockStatus:    http.StatusOK,
			expectError:   true,
			errorContains: "failed to decode response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockHTTPClient{
				GetFunc: func(url string) (*http.Response, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return &http.Response{
						StatusCode: tt.mockStatus,
						Body:       io.NopCloser(strings.NewReader(tt.mockResponse)),
					}, nil
				},
			}

			client := NewClientWithHTTP(mockClient)
			result, err := client.GeocodeCity(tt.cityName)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.Name)
				assert.NotZero(t, result.Latitude)
				assert.NotZero(t, result.Longitude)
			}
		})
	}
}

func TestFormatFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{
			name:     "positive float",
			input:    52.52,
			expected: "52.52",
		},
		{
			name:     "negative float",
			input:    -13.405,
			expected: "-13.405",
		},
		{
			name:     "zero",
			input:    0.0,
			expected: "0",
		},
		{
			name:     "integer value",
			input:    10.0,
			expected: "10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatFloat(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGetCurrentWeatherResponseBodyClose ensures the response body is properly closed
func TestGetCurrentWeatherResponseBodyClose(t *testing.T) {
	bodyClosed := false
	mockClient := &MockHTTPClient{
		GetFunc: func(url string) (*http.Response, error) {
			body := &mockReadCloser{
				Reader: bytes.NewReader([]byte(`{
					"latitude": 52.52,
					"longitude": 13.405,
					"timezone": "Europe/Berlin",
					"current": {
						"temperature_2m": 15.5,
						"relative_humidity_2m": 65,
						"wind_speed_10m": 12.3,
						"weather_code": 0
					}
				}`)),
				closeFunc: func() error {
					bodyClosed = true
					return nil
				},
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       body,
			}, nil
		},
	}

	client := NewClientWithHTTP(mockClient)
	_, err := client.GetCurrentWeather(52.52, 13.405)
	require.NoError(t, err)
	assert.True(t, bodyClosed, "Response body should be closed")
}

// mockReadCloser is a helper for testing body closing
type mockReadCloser struct {
	io.Reader
	closeFunc func() error
}

func (m *mockReadCloser) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	return nil
}
