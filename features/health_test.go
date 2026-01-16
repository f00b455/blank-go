package features

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"time"

	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin"

	"github.com/f00b455/blank-go/internal/handlers"
)

type healthFeatureContext struct {
	router       *gin.Engine
	response     *httptest.ResponseRecorder
	lastResponse map[string]interface{}
	startTime    time.Time
}

func (ctx *healthFeatureContext) reset() {
	gin.SetMode(gin.TestMode)
	ctx.startTime = time.Now()
	ctx.setupRouter()
	ctx.response = nil
	ctx.lastResponse = nil
}

func (ctx *healthFeatureContext) setupRouter() {
	ctx.router = gin.New()
	api := ctx.router.Group("/api/v1")
	{
		api.GET("/health/detailed", handlers.DetailedHealthCheck(ctx.startTime))
	}
}

func (ctx *healthFeatureContext) theAPIServerIsRunning() error {
	return nil
}

func (ctx *healthFeatureContext) iRequestDetailedHealthCheckAt(path string) error {
	req, _ := http.NewRequest("GET", path, nil)
	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)

	if ctx.response.Code == http.StatusOK {
		if err := json.Unmarshal(ctx.response.Body.Bytes(), &ctx.lastResponse); err != nil {
			return err
		}
	}

	return nil
}

func (ctx *healthFeatureContext) theResponseStatusShouldBe(expected int) error {
	if ctx.response.Code != expected {
		return fmt.Errorf("expected status %d, got %d", expected, ctx.response.Code)
	}
	return nil
}

func (ctx *healthFeatureContext) theResponseShouldContainFieldWithValue(field, value string) error {
	if ctx.lastResponse == nil {
		return fmt.Errorf("no response data available")
	}

	actualValue, ok := ctx.lastResponse[field]
	if !ok {
		return fmt.Errorf("field %q not found in response", field)
	}

	if actualValue != value {
		return fmt.Errorf("expected %q to be %q, got %v", field, value, actualValue)
	}

	return nil
}

func (ctx *healthFeatureContext) theResponseShouldContainField(field string) error {
	if ctx.lastResponse == nil {
		return fmt.Errorf("no response data available")
	}

	if _, ok := ctx.lastResponse[field]; !ok {
		return fmt.Errorf("field %q not found in response", field)
	}

	return nil
}

func (ctx *healthFeatureContext) theResponseShouldContainSystemMetricsWithFields(table *godog.Table) error {
	if ctx.lastResponse == nil {
		return fmt.Errorf("no response data available")
	}

	systemData, ok := ctx.lastResponse["system"]
	if !ok {
		return fmt.Errorf("system field not found in response")
	}

	system, ok := systemData.(map[string]interface{})
	if !ok {
		return fmt.Errorf("system field is not a map")
	}

	for i := 1; i < len(table.Rows); i++ {
		field := table.Rows[i].Cells[0].Value
		if _, ok := system[field]; !ok {
			return fmt.Errorf("system metric %q not found", field)
		}
	}

	return nil
}

func (ctx *healthFeatureContext) theResponseShouldContainChecksWithFieldWithValue(field, value string) error {
	if ctx.lastResponse == nil {
		return fmt.Errorf("no response data available")
	}

	checksData, ok := ctx.lastResponse["checks"]
	if !ok {
		return fmt.Errorf("checks field not found in response")
	}

	checks, ok := checksData.(map[string]interface{})
	if !ok {
		return fmt.Errorf("checks field is not a map")
	}

	actualValue, ok := checks[field]
	if !ok {
		return fmt.Errorf("check %q not found in checks", field)
	}

	if actualValue != value {
		return fmt.Errorf("expected check %q to be %q, got %v", field, value, actualValue)
	}

	return nil
}

func (ctx *healthFeatureContext) theSystemMetricShouldBeGreaterThan(metric string, threshold float64) error {
	if ctx.lastResponse == nil {
		return fmt.Errorf("no response data available")
	}

	systemData, ok := ctx.lastResponse["system"]
	if !ok {
		return fmt.Errorf("system field not found")
	}

	system, ok := systemData.(map[string]interface{})
	if !ok {
		return fmt.Errorf("system field is not a map")
	}

	value, ok := system[metric]
	if !ok {
		return fmt.Errorf("metric %q not found in system", metric)
	}

	var numValue float64
	switch v := value.(type) {
	case float64:
		numValue = v
	case int:
		numValue = float64(v)
	default:
		return fmt.Errorf("metric %q is not a number", metric)
	}

	if numValue <= threshold {
		return fmt.Errorf("expected %q to be greater than %f, got %f", metric, threshold, numValue)
	}

	return nil
}

func (ctx *healthFeatureContext) theSystemMetricShouldBeGreaterThanOrEqualTo(metric string, threshold float64) error {
	if ctx.lastResponse == nil {
		return fmt.Errorf("no response data available")
	}

	systemData, ok := ctx.lastResponse["system"]
	if !ok {
		return fmt.Errorf("system field not found")
	}

	system, ok := systemData.(map[string]interface{})
	if !ok {
		return fmt.Errorf("system field is not a map")
	}

	value, ok := system[metric]
	if !ok {
		return fmt.Errorf("metric %q not found in system", metric)
	}

	var numValue float64
	switch v := value.(type) {
	case float64:
		numValue = v
	case int:
		numValue = float64(v)
	default:
		return fmt.Errorf("metric %q is not a number", metric)
	}

	if numValue < threshold {
		return fmt.Errorf("expected %q to be >= %f, got %f", metric, threshold, numValue)
	}

	return nil
}

func (ctx *healthFeatureContext) theUptimeShouldBeGreaterThan(threshold float64) error {
	if ctx.lastResponse == nil {
		return fmt.Errorf("no response data available")
	}

	value, ok := ctx.lastResponse["uptime_seconds"]
	if !ok {
		return fmt.Errorf("uptime_seconds field not found")
	}

	var numValue float64
	switch v := value.(type) {
	case float64:
		numValue = v
	case int:
		numValue = float64(v)
	default:
		return fmt.Errorf("uptime_seconds is not a number")
	}

	if numValue <= threshold {
		return fmt.Errorf("expected uptime to be > %f, got %f", threshold, numValue)
	}

	return nil
}

func (ctx *healthFeatureContext) theResponseShouldContainVersionInformation() error {
	if ctx.lastResponse == nil {
		return fmt.Errorf("no response data available")
	}

	version, ok := ctx.lastResponse["version"]
	if !ok {
		return fmt.Errorf("version field not found in response")
	}

	versionStr, ok := version.(string)
	if !ok || versionStr == "" {
		return fmt.Errorf("version is empty or not a string")
	}

	return nil
}

func (ctx *healthFeatureContext) theVersionShouldMatchSemanticVersioningPattern() error {
	if ctx.lastResponse == nil {
		return fmt.Errorf("no response data available")
	}

	version, ok := ctx.lastResponse["version"]
	if !ok {
		return fmt.Errorf("version field not found")
	}

	versionStr, ok := version.(string)
	if !ok {
		return fmt.Errorf("version is not a string")
	}

	semverPattern := `^\d+\.\d+\.\d+(-[a-zA-Z0-9.-]+)?(\+[a-zA-Z0-9.-]+)?$`
	matched, err := regexp.MatchString(semverPattern, versionStr)
	if err != nil {
		return err
	}

	if !matched {
		// Allow "dev" or "unknown" for development
		if !strings.Contains(versionStr, "dev") && !strings.Contains(versionStr, "unknown") {
			return fmt.Errorf("version %q does not match semver pattern", versionStr)
		}
	}

	return nil
}

func InitializeHealthScenario(ctx *godog.ScenarioContext) {
	feature := &healthFeatureContext{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		feature.reset()
		return ctx, nil
	})

	ctx.Step(`^the API server is running$`, feature.theAPIServerIsRunning)
	ctx.Step(`^I request detailed health check at "([^"]*)"$`, feature.iRequestDetailedHealthCheckAt)
	ctx.Step(`^the response status should be (\d+)$`, feature.theResponseStatusShouldBe)
	ctx.Step(`^the response should contain field "([^"]*)" with value "([^"]*)"$`, feature.theResponseShouldContainFieldWithValue)
	ctx.Step(`^the response should contain field "([^"]*)"$`, feature.theResponseShouldContainField)
	ctx.Step(`^the response should contain system metrics with fields:$`, feature.theResponseShouldContainSystemMetricsWithFields)
	ctx.Step(`^the response should contain checks with field "([^"]*)" with value "([^"]*)"$`, feature.theResponseShouldContainChecksWithFieldWithValue)
	ctx.Step(`^the system metric "([^"]*)" should be greater than (\d+)$`, feature.theSystemMetricShouldBeGreaterThan)
	ctx.Step(`^the system metric "([^"]*)" should be greater than or equal to (\d+)$`, feature.theSystemMetricShouldBeGreaterThanOrEqualTo)
	ctx.Step(`^the uptime should be greater than (\d+)$`, feature.theUptimeShouldBeGreaterThan)
	ctx.Step(`^the response should contain version information$`, feature.theResponseShouldContainVersionInformation)
	ctx.Step(`^the version should match semantic versioning pattern$`, feature.theVersionShouldMatchSemanticVersioningPattern)
}
