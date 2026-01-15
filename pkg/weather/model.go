package weather

// Location represents geographic location information
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	City      string  `json:"city,omitempty"`
	Timezone  string  `json:"timezone,omitempty"`
}

// CurrentWeather represents current weather conditions
type CurrentWeather struct {
	Temperature        float64 `json:"temperature"`
	Humidity           int     `json:"humidity"`
	WindSpeed          float64 `json:"wind_speed"`
	WeatherCode        int     `json:"weather_code"`
	WeatherDescription string  `json:"weather_description"`
}

// Units represents measurement units for weather data
type Units struct {
	Temperature string `json:"temperature"`
	WindSpeed   string `json:"wind_speed"`
	Humidity    string `json:"humidity"`
}

// WeatherResponse represents the complete weather response
type WeatherResponse struct {
	Location Location       `json:"location"`
	Current  CurrentWeather `json:"current"`
	Units    Units          `json:"units"`
}

// ForecastDay represents a single day forecast
type ForecastDay struct {
	Date                     string  `json:"date"`
	TemperatureMax           float64 `json:"temperature_max"`
	TemperatureMin           float64 `json:"temperature_min"`
	PrecipitationProbability int     `json:"precipitation_probability"`
	WeatherCode              int     `json:"weather_code"`
	WeatherDescription       string  `json:"weather_description"`
}

// ForecastResponse represents the complete forecast response
type ForecastResponse struct {
	Location Location      `json:"location"`
	Forecast []ForecastDay `json:"forecast"`
}

// GeocodingResult represents a geocoding API result
type GeocodingResult struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
}

// WeatherCodeDescriptions maps Open-Meteo weather codes to descriptions
var WeatherCodeDescriptions = map[int]string{
	0:  "Clear sky",
	1:  "Mainly clear",
	2:  "Partly cloudy",
	3:  "Overcast",
	45: "Foggy",
	48: "Depositing rime fog",
	51: "Light drizzle",
	53: "Moderate drizzle",
	55: "Dense drizzle",
	56: "Light freezing drizzle",
	57: "Dense freezing drizzle",
	61: "Slight rain",
	63: "Moderate rain",
	65: "Heavy rain",
	66: "Light freezing rain",
	67: "Heavy freezing rain",
	71: "Slight snow fall",
	73: "Moderate snow fall",
	75: "Heavy snow fall",
	77: "Snow grains",
	80: "Slight rain showers",
	81: "Moderate rain showers",
	82: "Violent rain showers",
	85: "Slight snow showers",
	86: "Heavy snow showers",
	95: "Thunderstorm",
	96: "Thunderstorm with slight hail",
	99: "Thunderstorm with heavy hail",
}

// GetWeatherDescription returns the description for a weather code
func GetWeatherDescription(code int) string {
	if desc, ok := WeatherCodeDescriptions[code]; ok {
		return desc
	}
	return "Unknown"
}
