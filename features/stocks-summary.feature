# Issue: #19
# URL: https://github.com/f00b455/blank-go/issues/19
@pkg(stocks) @issue-19
Feature: Stock Market Daily Summary
  As an API user
  I want to retrieve daily stock market summaries
  So that I can display current stock information in my application

  Background:
    Given the Yahoo Finance API is available

  @happy-path @stock-summary
  Scenario: Get daily summary for a single stock ticker
    When I request stock summary for ticker "AAPL"
    Then the response status should be 200
    And the response should contain ticker "AAPL"
    And the response should contain current price
    And the response should contain open price
    And the response should contain high price
    And the response should contain low price
    And the response should contain change value
    And the response should contain change percentage
    And the response should contain volume
    And the response should contain currency

  @happy-path @stock-summary
  Scenario: Stock summary includes company name
    When I request stock summary for ticker "AAPL"
    Then the response status should be 200
    And the response should contain company name
    And the company name should not be empty

  @happy-path @stock-summary
  Scenario: Stock summary includes date
    When I request stock summary for ticker "AAPL"
    Then the response status should be 200
    And the response should contain date
    And the date should be in format "YYYY-MM-DD"

  @error-handling
  Scenario: Invalid stock ticker
    When I request stock summary for ticker "INVALID_TICKER_XYZ"
    Then the response status should be 404
    And the error message should indicate "ticker not found"

  @error-handling
  Scenario: Empty stock ticker
    When I request stock summary for ticker ""
    Then the response status should be 400
    And the error message should indicate "ticker is required"

  @caching
  Scenario: Cached stock data is returned for repeated requests
    When I request stock summary for ticker "AAPL"
    And I request stock summary for ticker "AAPL" again within cache TTL
    Then both responses should be identical
    And the second request should be served from cache
