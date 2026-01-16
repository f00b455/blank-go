package features

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin"

	"github.com/f00b455/blank-go/internal/handlers"
	"github.com/f00b455/blank-go/pkg/weather"
)

type weatherFeatureContext struct {
	router       *gin.Engine
	service      *weather.Service
	handler      *handlers.WeatherHandler
	mockClient   *MockWeatherClient
	response     *httptest.ResponseRecorder
	lastResponse map[string]interface{}
}

// MockWeatherClient is a mock for testing
type MockWeatherClient struct {
	currentWeatherFunc func(lat, lon float64) (*weather.WeatherResponse, error)
	forecastFunc       func(lat, lon float64, days int) (*weather.ForecastResponse, error)
	geocodeFunc        func(cityName string) (*weather.GeocodingResult, error)
}

func (m *MockWeatherClient) GetCurrentWeather(lat, lon float64) (*weather.WeatherResponse, error) {
	if m.currentWeatherFunc != nil {
		return m.currentWeatherFunc(lat, lon)
	}
	return &weather.WeatherResponse{
		Location: weather.Location{
			Latitude:  lat,
			Longitude: lon,
			Timezone:  "Europe/Berlin",
		},
		Current: weather.CurrentWeather{
			Temperature:        15.2,
			Humidity:           65,
			WindSpeed:          12.5,
			WeatherCode:        2,
			WeatherDescription: "Partly cloudy",
		},
		Units: weather.Units{
			Temperature: "°C",
			WindSpeed:   "km/h",
			Humidity:    "%",
		},
	}, nil
}

func (m *MockWeatherClient) GetForecast(lat, lon float64, days int) (*weather.ForecastResponse, error) {
	if m.forecastFunc != nil {
		return m.forecastFunc(lat, lon, days)
	}

	forecast := make([]weather.ForecastDay, days)
	for i := 0; i < days; i++ {
		forecast[i] = weather.ForecastDay{
			Date:                     fmt.Sprintf("2025-01-%02d", 15+i),
			TemperatureMax:           18.5,
			TemperatureMin:           8.2,
			PrecipitationProbability: 20,
			WeatherCode:              2,
			WeatherDescription:       "Partly cloudy",
		}
	}

	return &weather.ForecastResponse{
		Location: weather.Location{
			Latitude:  lat,
			Longitude: lon,
			Timezone:  "Europe/Berlin",
		},
		Forecast: forecast,
	}, nil
}

func (m *MockWeatherClient) GeocodeCity(cityName string) (*weather.GeocodingResult, error) {
	if m.geocodeFunc != nil {
		return m.geocodeFunc(cityName)
	}

	if cityName == "NonExistentCity123456" {
		return nil, fmt.Errorf("city not found")
	}

	return &weather.GeocodingResult{
		Name:      cityName,
		Latitude:  52.52,
		Longitude: 13.41,
		Timezone:  "Europe/Berlin",
	}, nil
}

func (ctx *weatherFeatureContext) reset() {
	gin.SetMode(gin.TestMode)
	ctx.mockClient = &MockWeatherClient{}
	ctx.service = weather.NewService(ctx.mockClient)
	ctx.handler = handlers.NewWeatherHandler(ctx.service)
	ctx.setupRouter()

	ctx.response = nil
	ctx.lastResponse = nil
}

func (ctx *weatherFeatureContext) setupRouter() {
	ctx.router = gin.New()
	api := ctx.router.Group("/api/v1")
	{
		api.GET("/weather", ctx.handler.GetCurrentWeather)
		api.GET("/weather/forecast", ctx.handler.GetForecast)
		api.GET("/weather/cities/:city", ctx.handler.GetWeatherByCity)
	}
}

func (ctx *weatherFeatureContext) theWeatherAPIIsAvailable() error {
	return nil
}

func (ctx *weatherFeatureContext) theGeocodingAPIIsAvailable() error {
	return nil
}

func (ctx *weatherFeatureContext) iRequestCurrentWeatherFor(lat, lon string) error {
	url := fmt.Sprintf("/api/v1/weather?lat=%s&lon=%s", lat, lon)
	req, _ := http.NewRequest("GET", url, nil)
	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)
	return ctx.parseResponse()
}

func (ctx *weatherFeatureContext) iRequestWeatherForecastForDays(lat, lon, days string) error {
	url := fmt.Sprintf("/api/v1/weather/forecast?lat=%s&lon=%s&days=%s", lat, lon, days)
	req, _ := http.NewRequest("GET", url, nil)
	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)
	return ctx.parseResponse()
}

func (ctx *weatherFeatureContext) iRequestWeatherForecastWithoutDays(lat, lon string) error {
	url := fmt.Sprintf("/api/v1/weather/forecast?lat=%s&lon=%s", lat, lon)
	req, _ := http.NewRequest("GET", url, nil)
	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)
	return ctx.parseResponse()
}

func (ctx *weatherFeatureContext) iRequestWeatherForCity(cityName string) error {
	url := fmt.Sprintf("/api/v1/weather/cities/%s", cityName)
	req, _ := http.NewRequest("GET", url, nil)
	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)
	return ctx.parseResponse()
}

func (ctx *weatherFeatureContext) parseResponse() error {
	if ctx.response.Code != http.StatusOK {
		return nil
	}
	return json.Unmarshal(ctx.response.Body.Bytes(), &ctx.lastResponse)
}

func (ctx *weatherFeatureContext) theResponseStatusShouldBe(expectedCode int) error {
	if ctx.response.Code != expectedCode {
		return fmt.Errorf("expected status %d, got %d", expectedCode, ctx.response.Code)
	}
	return nil
}

func (ctx *weatherFeatureContext) theResponseShouldContainLocationData() error {
	if _, ok := ctx.lastResponse["location"]; !ok {
		return fmt.Errorf("response missing location data")
	}
	return nil
}

func (ctx *weatherFeatureContext) theResponseShouldContainCurrentWeatherData() error {
	if _, ok := ctx.lastResponse["current"]; !ok {
		return fmt.Errorf("response missing current weather data")
	}
	return nil
}

func (ctx *weatherFeatureContext) theCurrentWeatherShouldIncludeTemperature() error {
	current := ctx.lastResponse["current"].(map[string]interface{})
	if _, ok := current["temperature"]; !ok {
		return fmt.Errorf("current weather missing temperature")
	}
	return nil
}

func (ctx *weatherFeatureContext) theCurrentWeatherShouldIncludeHumidity() error {
	current := ctx.lastResponse["current"].(map[string]interface{})
	if _, ok := current["humidity"]; !ok {
		return fmt.Errorf("current weather missing humidity")
	}
	return nil
}

func (ctx *weatherFeatureContext) theCurrentWeatherShouldIncludeWindSpeed() error {
	current := ctx.lastResponse["current"].(map[string]interface{})
	if _, ok := current["wind_speed"]; !ok {
		return fmt.Errorf("current weather missing wind_speed")
	}
	return nil
}

func (ctx *weatherFeatureContext) theCurrentWeatherShouldIncludeWeatherDescription() error {
	current := ctx.lastResponse["current"].(map[string]interface{})
	if _, ok := current["weather_description"]; !ok {
		return fmt.Errorf("current weather missing weather_description")
	}
	return nil
}

func (ctx *weatherFeatureContext) theUnitsShouldSpecifyForTemperature(unit string) error {
	units := ctx.lastResponse["units"].(map[string]interface{})
	if units["temperature"] != unit {
		return fmt.Errorf("expected temperature unit %s, got %s", unit, units["temperature"])
	}
	return nil
}

func (ctx *weatherFeatureContext) theUnitsShouldSpecifyForWindSpeed(unit string) error {
	units := ctx.lastResponse["units"].(map[string]interface{})
	if units["wind_speed"] != unit {
		return fmt.Errorf("expected wind_speed unit %s, got %s", unit, units["wind_speed"])
	}
	return nil
}

func (ctx *weatherFeatureContext) theUnitsShouldSpecifyForHumidity(unit string) error {
	units := ctx.lastResponse["units"].(map[string]interface{})
	if units["humidity"] != unit {
		return fmt.Errorf("expected humidity unit %s, got %s", unit, units["humidity"])
	}
	return nil
}

func (ctx *weatherFeatureContext) theErrorMessageShouldIndicate(message string) error {
	body := ctx.response.Body.String()
	if !strings.Contains(strings.ToLower(body), strings.ToLower(message)) {
		return fmt.Errorf("expected error message containing '%s', got '%s'", message, body)
	}
	return nil
}

func (ctx *weatherFeatureContext) theResponseShouldContainForecastDays(days int) error {
	forecast, ok := ctx.lastResponse["forecast"].([]interface{})
	if !ok {
		return fmt.Errorf("response missing forecast array")
	}
	if len(forecast) != days {
		return fmt.Errorf("expected %d forecast days, got %d", days, len(forecast))
	}
	return nil
}

func (ctx *weatherFeatureContext) eachForecastDayShouldIncludeDate() error {
	forecast := ctx.lastResponse["forecast"].([]interface{})
	for i, day := range forecast {
		dayMap := day.(map[string]interface{})
		if _, ok := dayMap["date"]; !ok {
			return fmt.Errorf("forecast day %d missing date", i)
		}
	}
	return nil
}

func (ctx *weatherFeatureContext) eachForecastDayShouldIncludeMaxTemperature() error {
	forecast := ctx.lastResponse["forecast"].([]interface{})
	for i, day := range forecast {
		dayMap := day.(map[string]interface{})
		if _, ok := dayMap["temperature_max"]; !ok {
			return fmt.Errorf("forecast day %d missing temperature_max", i)
		}
	}
	return nil
}

func (ctx *weatherFeatureContext) eachForecastDayShouldIncludeMinTemperature() error {
	forecast := ctx.lastResponse["forecast"].([]interface{})
	for i, day := range forecast {
		dayMap := day.(map[string]interface{})
		if _, ok := dayMap["temperature_min"]; !ok {
			return fmt.Errorf("forecast day %d missing temperature_min", i)
		}
	}
	return nil
}

func (ctx *weatherFeatureContext) eachForecastDayShouldIncludePrecipitationProbability() error {
	forecast := ctx.lastResponse["forecast"].([]interface{})
	for i, day := range forecast {
		dayMap := day.(map[string]interface{})
		if _, ok := dayMap["precipitation_probability"]; !ok {
			return fmt.Errorf("forecast day %d missing precipitation_probability", i)
		}
	}
	return nil
}

func (ctx *weatherFeatureContext) eachForecastDayShouldIncludeWeatherDescription() error {
	forecast := ctx.lastResponse["forecast"].([]interface{})
	for i, day := range forecast {
		dayMap := day.(map[string]interface{})
		if _, ok := dayMap["weather_description"]; !ok {
			return fmt.Errorf("forecast day %d missing weather_description", i)
		}
	}
	return nil
}

func (ctx *weatherFeatureContext) theLocationShouldIncludeCityName(expectedCity string) error {
	location := ctx.lastResponse["location"].(map[string]interface{})
	city, ok := location["city"].(string)
	if !ok {
		return fmt.Errorf("location missing city name")
	}
	if city != expectedCity {
		return fmt.Errorf("expected city %s, got %s", expectedCity, city)
	}
	return nil
}

func (ctx *weatherFeatureContext) theLocationShouldIncludeCoordinates() error {
	location := ctx.lastResponse["location"].(map[string]interface{})
	if _, ok := location["latitude"]; !ok {
		return fmt.Errorf("location missing latitude")
	}
	if _, ok := location["longitude"]; !ok {
		return fmt.Errorf("location missing longitude")
	}
	return nil
}

func (ctx *weatherFeatureContext) theLocationShouldIncludeTimezone() error {
	location := ctx.lastResponse["location"].(map[string]interface{})
	if _, ok := location["timezone"]; !ok {
		return fmt.Errorf("location missing timezone")
	}
	return nil
}

func (ctx *weatherFeatureContext) theLocationCityNameShouldContain(expectedCity string) error {
	location := ctx.lastResponse["location"].(map[string]interface{})
	city, ok := location["city"].(string)
	if !ok {
		return fmt.Errorf("location missing city name")
	}
	if !strings.Contains(city, expectedCity) {
		return fmt.Errorf("expected city containing %s, got %s", expectedCity, city)
	}
	return nil
}

func InitializeWeatherScenario(sc *godog.ScenarioContext) {
	ctx := &weatherFeatureContext{}

	sc.Before(func(c context.Context, sc *godog.Scenario) (context.Context, error) {
		ctx.reset()
		return c, nil
	})

	// Background steps
	sc.Step(`^the Weather API is available$`, ctx.theWeatherAPIIsAvailable)
	sc.Step(`^the Geocoding API is available$`, ctx.theGeocodingAPIIsAvailable)

	// Current weather steps
	sc.Step(`^I request current weather for latitude "([^"]*)" and longitude "([^"]*)"$`, ctx.iRequestCurrentWeatherFor)

	// Forecast steps
	sc.Step(`^I request weather forecast for latitude "([^"]*)" and longitude "([^"]*)" for (\d+) days$`, ctx.iRequestWeatherForecastForDays)
	sc.Step(`^I request weather forecast for latitude "([^"]*)" and longitude "([^"]*)" without specifying days$`, ctx.iRequestWeatherForecastWithoutDays)

	// City weather steps
	sc.Step(`^I request weather for city "([^"]*)"$`, ctx.iRequestWeatherForCity)

	// Response validation steps
	sc.Step(`^the response status should be (\d+)$`, ctx.theResponseStatusShouldBe)
	sc.Step(`^the response should contain location data$`, ctx.theResponseShouldContainLocationData)
	sc.Step(`^the response should contain current weather data$`, ctx.theResponseShouldContainCurrentWeatherData)
	sc.Step(`^the current weather should include temperature$`, ctx.theCurrentWeatherShouldIncludeTemperature)
	sc.Step(`^the current weather should include humidity$`, ctx.theCurrentWeatherShouldIncludeHumidity)
	sc.Step(`^the current weather should include wind speed$`, ctx.theCurrentWeatherShouldIncludeWindSpeed)
	sc.Step(`^the current weather should include weather description$`, ctx.theCurrentWeatherShouldIncludeWeatherDescription)
	sc.Step(`^the units should specify "([^"]*)" for temperature$`, ctx.theUnitsShouldSpecifyForTemperature)
	sc.Step(`^the units should specify "([^"]*)" for wind speed$`, ctx.theUnitsShouldSpecifyForWindSpeed)
	sc.Step(`^the units should specify "([^"]*)" for humidity$`, ctx.theUnitsShouldSpecifyForHumidity)
	sc.Step(`^the error message should indicate "([^"]*)"$`, ctx.theErrorMessageShouldIndicate)

	// Forecast validation steps
	sc.Step(`^the response should contain (\d+) forecast days$`, ctx.theResponseShouldContainForecastDays)
	sc.Step(`^each forecast day should include date$`, ctx.eachForecastDayShouldIncludeDate)
	sc.Step(`^each forecast day should include max temperature$`, ctx.eachForecastDayShouldIncludeMaxTemperature)
	sc.Step(`^each forecast day should include min temperature$`, ctx.eachForecastDayShouldIncludeMinTemperature)
	sc.Step(`^each forecast day should include precipitation probability$`, ctx.eachForecastDayShouldIncludePrecipitationProbability)
	sc.Step(`^each forecast day should include weather description$`, ctx.eachForecastDayShouldIncludeWeatherDescription)

	// City weather validation steps
	sc.Step(`^the location should include city name "([^"]*)"$`, ctx.theLocationShouldIncludeCityName)
	sc.Step(`^the location should include coordinates$`, ctx.theLocationShouldIncludeCoordinates)
	sc.Step(`^the location should include timezone$`, ctx.theLocationShouldIncludeTimezone)
	sc.Step(`^the location city name should contain "([^"]*)"$`, ctx.theLocationCityNameShouldContain)
}
