package weather

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWeatherClient is a mock implementation of WeatherClient
type MockWeatherClient struct {
	mock.Mock
}

func (m *MockWeatherClient) GetCurrentWeather(lat, lon float64) (*WeatherResponse, error) {
	args := m.Called(lat, lon)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*WeatherResponse), args.Error(1)
}

func (m *MockWeatherClient) GetForecast(lat, lon float64, days int) (*ForecastResponse, error) {
	args := m.Called(lat, lon, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ForecastResponse), args.Error(1)
}

func (m *MockWeatherClient) GeocodeCity(cityName string) (*GeocodingResult, error) {
	args := m.Called(cityName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GeocodingResult), args.Error(1)
}

func TestGetCurrentWeatherByCoords_Success(t *testing.T) {
	mockClient := new(MockWeatherClient)
	service := NewService(mockClient)

	expectedResp := &WeatherResponse{
		Location: Location{Latitude: 52.52, Longitude: 13.41},
		Current:  CurrentWeather{Temperature: 15.2},
	}

	mockClient.On("GetCurrentWeather", 52.52, 13.41).Return(expectedResp, nil)

	result, err := service.GetCurrentWeatherByCoords("52.52", "13.41")

	assert.NoError(t, err)
	assert.Equal(t, expectedResp, result)
	mockClient.AssertExpectations(t)
}

func TestGetCurrentWeatherByCoords_InvalidLatitude(t *testing.T) {
	mockClient := new(MockWeatherClient)
	service := NewService(mockClient)

	_, err := service.GetCurrentWeatherByCoords("invalid", "13.41")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid latitude")
}

func TestGetCurrentWeatherByCoords_InvalidLongitude(t *testing.T) {
	mockClient := new(MockWeatherClient)
	service := NewService(mockClient)

	_, err := service.GetCurrentWeatherByCoords("52.52", "invalid")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid longitude")
}

func TestGetCurrentWeatherByCoords_LatitudeOutOfRange(t *testing.T) {
	mockClient := new(MockWeatherClient)
	service := NewService(mockClient)

	_, err := service.GetCurrentWeatherByCoords("91.0", "13.41")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "latitude out of range")
}

func TestGetCurrentWeatherByCoords_LongitudeOutOfRange(t *testing.T) {
	mockClient := new(MockWeatherClient)
	service := NewService(mockClient)

	_, err := service.GetCurrentWeatherByCoords("52.52", "181.0")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "longitude out of range")
}

func TestGetForecastByCoords_Success(t *testing.T) {
	mockClient := new(MockWeatherClient)
	service := NewService(mockClient)

	expectedResp := &ForecastResponse{
		Location: Location{Latitude: 52.52, Longitude: 13.41},
		Forecast: []ForecastDay{{Date: "2025-01-15"}},
	}

	mockClient.On("GetForecast", 52.52, 13.41, 7).Return(expectedResp, nil)

	result, err := service.GetForecastByCoords("52.52", "13.41", 7)

	assert.NoError(t, err)
	assert.Equal(t, expectedResp, result)
	mockClient.AssertExpectations(t)
}

func TestGetForecastByCoords_InvalidDays(t *testing.T) {
	tests := []struct {
		name string
		days int
	}{
		{"Zero days", 0},
		{"Negative days", -1},
		{"Too many days", 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockWeatherClient)
			service := NewService(mockClient)

			_, err := service.GetForecastByCoords("52.52", "13.41", tt.days)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "days must be between 1 and 7")
		})
	}
}

func TestGetWeatherByCity_Success(t *testing.T) {
	mockClient := new(MockWeatherClient)
	service := NewService(mockClient)

	geocodeResp := &GeocodingResult{
		Name:      "Berlin",
		Latitude:  52.52,
		Longitude: 13.41,
		Timezone:  "Europe/Berlin",
	}

	weatherResp := &WeatherResponse{
		Location: Location{Latitude: 52.52, Longitude: 13.41},
		Current:  CurrentWeather{Temperature: 15.2},
	}

	mockClient.On("GeocodeCity", "Berlin").Return(geocodeResp, nil)
	mockClient.On("GetCurrentWeather", 52.52, 13.41).Return(weatherResp, nil)

	result, err := service.GetWeatherByCity("Berlin")

	assert.NoError(t, err)
	assert.Equal(t, "Berlin", result.Location.City)
	assert.Equal(t, "Europe/Berlin", result.Location.Timezone)
	mockClient.AssertExpectations(t)
}

func TestGetWeatherByCity_EmptyName(t *testing.T) {
	mockClient := new(MockWeatherClient)
	service := NewService(mockClient)

	_, err := service.GetWeatherByCity("")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "city name is required")
}

func TestGetWeatherByCity_CityNotFound(t *testing.T) {
	mockClient := new(MockWeatherClient)
	service := NewService(mockClient)

	mockClient.On("GeocodeCity", "NonExistent").Return(nil, fmt.Errorf("city not found"))

	_, err := service.GetWeatherByCity("NonExistent")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "city not found")
	mockClient.AssertExpectations(t)
}

func TestValidateLatitude(t *testing.T) {
	tests := []struct {
		name    string
		lat     float64
		wantErr bool
	}{
		{"Valid latitude", 52.52, false},
		{"Min latitude", -90.0, false},
		{"Max latitude", 90.0, false},
		{"Below min", -91.0, true},
		{"Above max", 91.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLatitude(tt.lat)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateLongitude(t *testing.T) {
	tests := []struct {
		name    string
		lon     float64
		wantErr bool
	}{
		{"Valid longitude", 13.41, false},
		{"Min longitude", -180.0, false},
		{"Max longitude", 180.0, false},
		{"Below min", -181.0, true},
		{"Above max", 181.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLongitude(tt.lon)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetWeatherDescription(t *testing.T) {
	tests := []struct {
		code     int
		expected string
	}{
		{0, "Clear sky"},
		{2, "Partly cloudy"},
		{61, "Slight rain"},
		{95, "Thunderstorm"},
		{999, "Unknown"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Code %d", tt.code), func(t *testing.T) {
			result := GetWeatherDescription(tt.code)
			assert.Equal(t, tt.expected, result)
		})
	}
}
