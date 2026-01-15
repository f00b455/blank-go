# Issue: #9
# URL: https://github.com/f00b455/blank-go/issues/9
@pkg(weather) @issue-9
Feature: Current Weather Retrieval
  As an API user
  I want to retrieve current weather data
  So that I can display real-time weather information in my application

  Background:
    Given the Weather API is available

  @happy-path @current-weather
  Scenario: Get current weather by coordinates
    When I request current weather for latitude "52.52" and longitude "13.41"
    Then the response status should be 200
    And the response should contain location data
    And the response should contain current weather data
    And the current weather should include temperature
    And the current weather should include humidity
    And the current weather should include wind speed
    And the current weather should include weather description

  @happy-path @current-weather
  Scenario: Response contains correct units
    When I request current weather for latitude "52.52" and longitude "13.41"
    Then the response status should be 200
    And the units should specify "°C" for temperature
    And the units should specify "km/h" for wind speed
    And the units should specify "%" for humidity

  @error-handling
  Scenario: Invalid latitude coordinate
    When I request current weather for latitude "invalid" and longitude "13.41"
    Then the response status should be 400
    And the error message should indicate "invalid latitude"

  @error-handling
  Scenario: Invalid longitude coordinate
    When I request current weather for latitude "52.52" and longitude "invalid"
    Then the response status should be 400
    And the error message should indicate "invalid longitude"

  @error-handling
  Scenario: Latitude out of range
    When I request current weather for latitude "91.0" and longitude "13.41"
    Then the response status should be 400
    And the error message should indicate "latitude out of range"

  @error-handling
  Scenario: Longitude out of range
    When I request current weather for latitude "52.52" and longitude "181.0"
    Then the response status should be 400
    And the error message should indicate "longitude out of range"
