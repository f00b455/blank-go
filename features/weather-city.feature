# Issue: #9
# URL: https://github.com/f00b455/blank-go/issues/9
@pkg(weather) @issue-9
Feature: Weather by City Name
  As an API user
  I want to retrieve weather data by city name
  So that I can get weather without knowing coordinates

  Background:
    Given the Weather API is available
    And the Geocoding API is available

  @happy-path @city-weather
  Scenario: Get weather for a known city
    When I request weather for city "Berlin"
    Then the response status should be 200
    And the response should contain location data
    And the location should include city name "Berlin"
    And the location should include coordinates
    And the location should include timezone
    And the response should contain current weather data

  @happy-path @city-weather
  Scenario: Get weather for city with multiple matches uses first result
    When I request weather for city "Paris"
    Then the response status should be 200
    And the response should contain location data
    And the location city name should contain "Paris"

  @error-handling
  Scenario: City not found
    When I request weather for city "NonExistentCity123456"
    Then the response status should be 404
    And the error message should indicate "city not found"

  @error-handling
  Scenario: Empty city name
    When I request weather for city ""
    Then the response status should be 404

  @error-handling
  Scenario: City name with special characters
    When I request weather for city "São Paulo"
    Then the response status should be 200
    And the location city name should contain "São Paulo"
