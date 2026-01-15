package weather

import (
	"fmt"
	"strconv"
)

const (
	minDays = 1
	maxDays = 7
	minLat  = -90.0
	maxLat  = 90.0
	minLon  = -180.0
	maxLon  = 180.0
)

// WeatherClient defines the interface for weather data retrieval
type WeatherClient interface {
	GetCurrentWeather(lat, lon float64) (*WeatherResponse, error)
	GetForecast(lat, lon float64, days int) (*ForecastResponse, error)
	GeocodeCity(cityName string) (*GeocodingResult, error)
}

// Service provides weather business logic
type Service struct {
	client WeatherClient
}

// NewService creates a new weather service
func NewService(client WeatherClient) *Service {
	return &Service{
		client: client,
	}
}

// GetCurrentWeatherByCoords retrieves current weather by coordinates
func (s *Service) GetCurrentWeatherByCoords(latStr, lonStr string) (*WeatherResponse, error) {
	lat, lon, err := parseAndValidateCoords(latStr, lonStr)
	if err != nil {
		return nil, err
	}

	return s.client.GetCurrentWeather(lat, lon)
}

// GetForecastByCoords retrieves weather forecast by coordinates
func (s *Service) GetForecastByCoords(latStr, lonStr string, days int) (*ForecastResponse, error) {
	lat, lon, err := parseAndValidateCoords(latStr, lonStr)
	if err != nil {
		return nil, err
	}

	if err := validateDays(days); err != nil {
		return nil, err
	}

	return s.client.GetForecast(lat, lon, days)
}

// GetWeatherByCity retrieves current weather by city name
func (s *Service) GetWeatherByCity(cityName string) (*WeatherResponse, error) {
	if cityName == "" {
		return nil, fmt.Errorf("city name is required")
	}

	geocode, err := s.client.GeocodeCity(cityName)
	if err != nil {
		return nil, err
	}

	weather, err := s.client.GetCurrentWeather(geocode.Latitude, geocode.Longitude)
	if err != nil {
		return nil, err
	}

	weather.Location.City = geocode.Name
	weather.Location.Timezone = geocode.Timezone

	return weather, nil
}

// parseAndValidateCoords parses and validates latitude and longitude
func parseAndValidateCoords(latStr, lonStr string) (float64, float64, error) {
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid latitude")
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid longitude")
	}

	if err := validateLatitude(lat); err != nil {
		return 0, 0, err
	}

	if err := validateLongitude(lon); err != nil {
		return 0, 0, err
	}

	return lat, lon, nil
}

// validateLatitude checks if latitude is within valid range
func validateLatitude(lat float64) error {
	if lat < minLat || lat > maxLat {
		return fmt.Errorf("latitude out of range")
	}
	return nil
}

// validateLongitude checks if longitude is within valid range
func validateLongitude(lon float64) error {
	if lon < minLon || lon > maxLon {
		return fmt.Errorf("longitude out of range")
	}
	return nil
}

// validateDays checks if days parameter is within valid range
func validateDays(days int) error {
	if days < minDays || days > maxDays {
		return fmt.Errorf("days must be between 1 and 7")
	}
	return nil
}
