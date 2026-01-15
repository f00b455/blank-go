package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	openMeteoBaseURL      = "https://api.open-meteo.com/v1/forecast"
	geocodingBaseURL      = "https://geocoding-api.open-meteo.com/v1/search"
	defaultRequestTimeout = 10 * time.Second
)

// HTTPClient interface for mocking HTTP requests
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// Client handles communication with Open-Meteo API
type Client struct {
	httpClient HTTPClient
	timeout    time.Duration
}

// NewClient creates a new Open-Meteo API client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: defaultRequestTimeout,
		},
		timeout: defaultRequestTimeout,
	}
}

// NewClientWithHTTP creates a client with custom HTTP client
func NewClientWithHTTP(httpClient HTTPClient) *Client {
	return &Client{
		httpClient: httpClient,
		timeout:    defaultRequestTimeout,
	}
}

// openMeteoCurrentResponse represents the API response for current weather
type openMeteoCurrentResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
	Current   struct {
		Temperature float64 `json:"temperature_2m"`
		Humidity    int     `json:"relative_humidity_2m"`
		WindSpeed   float64 `json:"wind_speed_10m"`
		WeatherCode int     `json:"weather_code"`
	} `json:"current"`
}

// openMeteoForecastResponse represents the API response for forecast
type openMeteoForecastResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
	Daily     struct {
		Time                     []string  `json:"time"`
		TemperatureMax           []float64 `json:"temperature_2m_max"`
		TemperatureMin           []float64 `json:"temperature_2m_min"`
		PrecipitationProbability []int     `json:"precipitation_probability_max"`
		WeatherCode              []int     `json:"weather_code"`
	} `json:"daily"`
}

// geocodingResponse represents the API response for geocoding
type geocodingResponse struct {
	Results []struct {
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Timezone  string  `json:"timezone"`
	} `json:"results"`
}

// GetCurrentWeather fetches current weather data
func (c *Client) GetCurrentWeather(lat, lon float64) (*WeatherResponse, error) {
	params := url.Values{}
	params.Set("latitude", formatFloat(lat))
	params.Set("longitude", formatFloat(lon))
	params.Set("current", "temperature_2m,relative_humidity_2m,wind_speed_10m,weather_code")

	apiURL := fmt.Sprintf("%s?%s", openMeteoBaseURL, params.Encode())

	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp openMeteoCurrentResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &WeatherResponse{
		Location: Location{
			Latitude:  apiResp.Latitude,
			Longitude: apiResp.Longitude,
			Timezone:  apiResp.Timezone,
		},
		Current: CurrentWeather{
			Temperature:        apiResp.Current.Temperature,
			Humidity:           apiResp.Current.Humidity,
			WindSpeed:          apiResp.Current.WindSpeed,
			WeatherCode:        apiResp.Current.WeatherCode,
			WeatherDescription: GetWeatherDescription(apiResp.Current.WeatherCode),
		},
		Units: Units{
			Temperature: "°C",
			WindSpeed:   "km/h",
			Humidity:    "%",
		},
	}, nil
}

// GetForecast fetches weather forecast data
func (c *Client) GetForecast(lat, lon float64, days int) (*ForecastResponse, error) {
	params := url.Values{}
	params.Set("latitude", formatFloat(lat))
	params.Set("longitude", formatFloat(lon))
	params.Set("daily", "temperature_2m_max,temperature_2m_min,precipitation_probability_max,weather_code")
	params.Set("forecast_days", strconv.Itoa(days))

	apiURL := fmt.Sprintf("%s?%s", openMeteoBaseURL, params.Encode())

	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch forecast data: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp openMeteoForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	forecast := make([]ForecastDay, len(apiResp.Daily.Time))
	for i := range apiResp.Daily.Time {
		forecast[i] = ForecastDay{
			Date:                     apiResp.Daily.Time[i],
			TemperatureMax:           apiResp.Daily.TemperatureMax[i],
			TemperatureMin:           apiResp.Daily.TemperatureMin[i],
			PrecipitationProbability: apiResp.Daily.PrecipitationProbability[i],
			WeatherCode:              apiResp.Daily.WeatherCode[i],
			WeatherDescription:       GetWeatherDescription(apiResp.Daily.WeatherCode[i]),
		}
	}

	return &ForecastResponse{
		Location: Location{
			Latitude:  apiResp.Latitude,
			Longitude: apiResp.Longitude,
			Timezone:  apiResp.Timezone,
		},
		Forecast: forecast,
	}, nil
}

// GeocodeCity converts city name to coordinates
func (c *Client) GeocodeCity(cityName string) (*GeocodingResult, error) {
	params := url.Values{}
	params.Set("name", cityName)
	params.Set("count", "1")
	params.Set("language", "en")
	params.Set("format", "json")

	apiURL := fmt.Sprintf("%s?%s", geocodingBaseURL, params.Encode())

	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to geocode city: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp geocodingResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(apiResp.Results) == 0 {
		return nil, fmt.Errorf("city not found")
	}

	result := apiResp.Results[0]
	return &GeocodingResult{
		Name:      result.Name,
		Latitude:  result.Latitude,
		Longitude: result.Longitude,
		Timezone:  result.Timezone,
	}, nil
}

// formatFloat formats a float64 to string with appropriate precision
func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
