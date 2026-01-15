package features

import (
	"testing"

	"github.com/cucumber/godog"
)

func TestTaskFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"task-crud.feature", "task-listing.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run task feature tests")
	}
}

func TestDAXFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeDAXScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"dax-import.feature", "dax-retrieval.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run DAX feature tests")
	}
}

func TestWeatherFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeWeatherScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"weather-current.feature", "weather-forecast.feature", "weather-city.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run Weather feature tests")
	}
}
