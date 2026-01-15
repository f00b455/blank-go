# Issue: #9
# URL: https://github.com/f00b455/blank-go/issues/9
@pkg(weather) @issue-9
Feature: Weather Forecast Retrieval
  As an API user
  I want to retrieve weather forecast data
  So that I can display future weather predictions in my application

  Background:
    Given the Weather API is available

  @happy-path @forecast
  Scenario: Get 7-day forecast by coordinates
    When I request weather forecast for latitude "52.52" and longitude "13.41" for 7 days
    Then the response status should be 200
    And the response should contain location data
    And the response should contain 7 forecast days
    And each forecast day should include date
    And each forecast day should include max temperature
    And each forecast day should include min temperature
    And each forecast day should include precipitation probability
    And each forecast day should include weather description

  @happy-path @forecast
  Scenario: Get 1-day forecast by coordinates
    When I request weather forecast for latitude "52.52" and longitude "13.41" for 1 days
    Then the response status should be 200
    And the response should contain 1 forecast days

  @happy-path @forecast
  Scenario: Default forecast days is 7
    When I request weather forecast for latitude "52.52" and longitude "13.41" without specifying days
    Then the response status should be 200
    And the response should contain 7 forecast days

  @error-handling
  Scenario: Invalid number of forecast days
    When I request weather forecast for latitude "52.52" and longitude "13.41" for 0 days
    Then the response status should be 400
    And the error message should indicate "days must be between 1 and 7"

  @error-handling
  Scenario: Too many forecast days requested
    When I request weather forecast for latitude "52.52" and longitude "13.41" for 8 days
    Then the response status should be 400
    And the error message should indicate "days must be between 1 and 7"

  @error-handling
  Scenario: Invalid latitude for forecast
    When I request weather forecast for latitude "invalid" and longitude "13.41" for 7 days
    Then the response status should be 400
    And the error message should indicate "invalid latitude"

  @error-handling
  Scenario: Invalid longitude for forecast
    When I request weather forecast for latitude "52.52" and longitude "invalid" for 7 days
    Then the response status should be 400
    And the error message should indicate "invalid longitude"
