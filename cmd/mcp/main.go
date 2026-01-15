// MCP Server for Weather Data
// Provides current weather information via Model Context Protocol
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/f00b455/blank-go/internal/version"
)

// JSON-RPC structures
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      any         `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MCP Protocol structures
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResult struct {
	ProtocolVersion string            `json:"protocolVersion"`
	ServerInfo      ServerInfo        `json:"serverInfo"`
	Capabilities    map[string]any    `json:"capabilities"`
}

type Tool struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	InputSchema InputSchema `json:"inputSchema"`
}

type InputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties,omitempty"`
	Required   []string            `json:"required,omitempty"`
}

type Property struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type ToolsListResult struct {
	Tools []Tool `json:"tools"`
}

type CallToolParams struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments,omitempty"`
}

type TextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type CallToolResult struct {
	Content []TextContent `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

// Weather data structures
type OpenMeteoResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
	Current   struct {
		Time             string  `json:"time"`
		Temperature2m    float64 `json:"temperature_2m"`
		RelativeHumidity int     `json:"relative_humidity_2m"`
		WeatherCode      int     `json:"weather_code"`
		WindSpeed10m     float64 `json:"wind_speed_10m"`
	} `json:"current"`
}

type GeocodingResponse struct {
	Results []struct {
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Country   string  `json:"country"`
	} `json:"results"`
}

// Weather code descriptions
var weatherCodes = map[int]string{
	0:  "Clear sky",
	1:  "Mainly clear",
	2:  "Partly cloudy",
	3:  "Overcast",
	45: "Foggy",
	48: "Depositing rime fog",
	51: "Light drizzle",
	53: "Moderate drizzle",
	55: "Dense drizzle",
	61: "Slight rain",
	63: "Moderate rain",
	65: "Heavy rain",
	71: "Slight snow",
	73: "Moderate snow",
	75: "Heavy snow",
	80: "Slight rain showers",
	81: "Moderate rain showers",
	82: "Violent rain showers",
	95: "Thunderstorm",
}

func getWeatherDescription(code int) string {
	if desc, ok := weatherCodes[code]; ok {
		return desc
	}
	return "Unknown"
}

func fetchWeather(city string) (string, error) {
	// First, geocode the city
	geoURL := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1", city)

	geoResp, err := http.Get(geoURL)
	if err != nil {
		return "", fmt.Errorf("geocoding failed: %w", err)
	}
	defer func() {
		_ = geoResp.Body.Close()
	}()

	var geoData GeocodingResponse
	if err := json.NewDecoder(geoResp.Body).Decode(&geoData); err != nil {
		return "", fmt.Errorf("failed to parse geocoding response: %w", err)
	}

	if len(geoData.Results) == 0 {
		return "", fmt.Errorf("city not found: %s", city)
	}

	location := geoData.Results[0]

	// Fetch weather data
	weatherURL := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&current=temperature_2m,relative_humidity_2m,weather_code,wind_speed_10m&timezone=auto",
		location.Latitude, location.Longitude,
	)

	weatherResp, err := http.Get(weatherURL)
	if err != nil {
		return "", fmt.Errorf("weather fetch failed: %w", err)
	}
	defer func() {
		_ = weatherResp.Body.Close()
	}()

	var weatherData OpenMeteoResponse
	if err := json.NewDecoder(weatherResp.Body).Decode(&weatherData); err != nil {
		return "", fmt.Errorf("failed to parse weather response: %w", err)
	}

	// Format the result
	result := fmt.Sprintf(`Weather for %s, %s
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🌡️  Temperature: %.1f°C
💧 Humidity: %d%%
💨 Wind Speed: %.1f km/h
🌤️  Conditions: %s
⏰ Updated: %s
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📍 Coordinates: %.4f, %.4f
🌐 Timezone: %s`,
		location.Name, location.Country,
		weatherData.Current.Temperature2m,
		weatherData.Current.RelativeHumidity,
		weatherData.Current.WindSpeed10m,
		getWeatherDescription(weatherData.Current.WeatherCode),
		weatherData.Current.Time,
		location.Latitude, location.Longitude,
		weatherData.Timezone,
	)

	return result, nil
}

func handleRequest(req JSONRPCRequest) JSONRPCResponse {
	switch req.Method {
	case "initialize":
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: InitializeResult{
				ProtocolVersion: "2024-11-05",
				ServerInfo: ServerInfo{
					Name:    "weather-mcp",
					Version: version.Version,
				},
				Capabilities: map[string]any{
					"tools": map[string]any{},
				},
			},
		}

	case "notifications/initialized":
		// No response needed for notifications
		return JSONRPCResponse{}

	case "tools/list":
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: ToolsListResult{
				Tools: []Tool{
					{
						Name:        "get_weather",
						Description: "Get current weather for a city. Returns temperature, humidity, wind speed, and conditions.",
						InputSchema: InputSchema{
							Type: "object",
							Properties: map[string]Property{
								"city": {
									Type:        "string",
									Description: "City name (e.g., 'Berlin', 'Munich', 'Hamburg')",
								},
							},
							Required: []string{"city"},
						},
					},
				},
			},
		}

	case "tools/call":
		var params CallToolParams
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &RPCError{Code: -32602, Message: "Invalid params"},
			}
		}

		if params.Name == "get_weather" {
			city, ok := params.Arguments["city"].(string)
			if !ok || city == "" {
				return JSONRPCResponse{
					JSONRPC: "2.0",
					ID:      req.ID,
					Result: CallToolResult{
						Content: []TextContent{{Type: "text", Text: "Error: city parameter is required"}},
						IsError: true,
					},
				}
			}

			weather, err := fetchWeather(city)
			if err != nil {
				return JSONRPCResponse{
					JSONRPC: "2.0",
					ID:      req.ID,
					Result: CallToolResult{
						Content: []TextContent{{Type: "text", Text: fmt.Sprintf("Error: %v", err)}},
						IsError: true,
					},
				}
			}

			return JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result: CallToolResult{
					Content: []TextContent{{Type: "text", Text: weather}},
				},
			}
		}

		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &RPCError{Code: -32601, Message: fmt.Sprintf("Unknown tool: %s", params.Name)},
		}

	default:
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &RPCError{Code: -32601, Message: fmt.Sprintf("Method not found: %s", req.Method)},
		}
	}
}

func main() {
	fmt.Fprintf(os.Stderr, "Weather MCP Server started at %s\n", time.Now().Format(time.RFC3339))

	reader := bufio.NewReader(os.Stdin)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "Read error: %v\n", err)
			continue
		}

		var req JSONRPCRequest
		if err := json.Unmarshal(line, &req); err != nil {
			fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
			continue
		}

		resp := handleRequest(req)

		// Don't send response for notifications
		if resp.ID == nil && resp.Result == nil && resp.Error == nil {
			continue
		}

		respBytes, err := json.Marshal(resp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Marshal error: %v\n", err)
			continue
		}

		fmt.Println(string(respBytes))
	}
}
