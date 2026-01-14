# Issue: #3
# URL: https://github.com/f00b455/blank-go/issues/3
@pkg(dax) @issue-3
Feature: DAX Financial Data Retrieval
  As a financial analyst
  I want to query DAX financial data
  So that I can perform analysis on historical metrics

  Background:
    Given the DAX API is available
    And the following DAX records exist:
      | company       | ticker | report_type | metric       | year | value         | currency |
      | Siemens AG    | SIE    | income      | EBITDA       | 2025 | 15859000000.0 | EUR      |
      | Siemens AG    | SIE    | income      | Net Income   | 2025 | 9620000000.0  | EUR      |
      | Siemens AG    | SIE    | income      | EBITDA       | 2024 | 14500000000.0 | EUR      |
      | SAP SE        | SAP    | income      | EBITDA       | 2025 | 12000000000.0 | EUR      |
      | SAP SE        | SAP    | income      | Net Income   | 2025 | 8500000000.0  | EUR      |

  @happy-path @retrieval
  Scenario: Get all DAX records with pagination
    When I request all DAX records with page 1 and limit 10
    Then the response status should be 200
    And the response should contain 5 records
    And the response should include pagination metadata

  @happy-path @retrieval
  Scenario: Filter DAX records by ticker
    When I request DAX records for ticker "SIE"
    Then the response status should be 200
    And the response should contain 3 records
    And all records should have ticker "SIE"

  @happy-path @retrieval
  Scenario: Filter DAX records by ticker and year
    When I request DAX records for ticker "SIE" and year 2025
    Then the response status should be 200
    And the response should contain 2 records
    And all records should have ticker "SIE" and year 2025

  @happy-path @retrieval
  Scenario: Filter DAX records by year
    When I request DAX records for year 2025
    Then the response status should be 200
    And the response should contain 4 records
    And all records should have year 2025

  @happy-path @retrieval
  Scenario: Get available metrics for a ticker
    When I request available metrics for ticker "SIE"
    Then the response status should be 200
    And the response should contain metrics "EBITDA,Net Income"

  @happy-path @retrieval
  Scenario: Get available metrics for non-existent ticker
    When I request available metrics for ticker "UNKNOWN"
    Then the response status should be 200
    And the response should contain 0 metrics

  @happy-path @retrieval
  Scenario: Pagination works correctly
    When I request all DAX records with page 1 and limit 2
    Then the response status should be 200
    And the response should contain 2 records
    When I request all DAX records with page 2 and limit 2
    Then the response status should be 200
    And the response should contain 2 records
    When I request all DAX records with page 3 and limit 2
    Then the response status should be 200
    And the response should contain 1 records
