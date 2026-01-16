# Issue: #19
# URL: https://github.com/f00b455/blank-go/issues/19
@pkg(stocks) @issue-19
Feature: Batch Stock Market Summary
  As an API user
  I want to retrieve daily summaries for multiple stocks at once
  So that I can efficiently display multiple stock information

  Background:
    Given the Yahoo Finance API is available

  @happy-path @batch-summary
  Scenario: Get daily summary for multiple stock tickers
    When I request batch stock summary for tickers "AAPL,GOOGL,MSFT"
    Then the response status should be 200
    And the response should contain 3 stock summaries
    And the response should include ticker "AAPL"
    And the response should include ticker "GOOGL"
    And the response should include ticker "MSFT"

  @happy-path @batch-summary
  Scenario: Single ticker in batch request
    When I request batch stock summary for tickers "AAPL"
    Then the response status should be 200
    And the response should contain 1 stock summary

  @error-handling
  Scenario: Empty tickers parameter
    When I request batch stock summary for tickers ""
    Then the response status should be 400
    And the error message should indicate "tickers parameter is required"

  @error-handling
  Scenario: Partial success with invalid ticker
    When I request batch stock summary for tickers "AAPL,INVALID_XYZ,MSFT"
    Then the response status should be 207
    And the response should contain 2 successful summaries
    And the response should contain 1 error
    And the error should indicate ticker "INVALID_XYZ" not found

  @rate-limiting
  Scenario: Batch request respects rate limits
    When I request batch stock summary for tickers "AAPL,GOOGL,MSFT,TSLA,AMZN"
    Then the response status should be 200
    And the request should not exceed API rate limits
    And all 5 stock summaries should be returned
