# Issue: #3
# URL: https://github.com/f00b455/blank-go/issues/3
@pkg(dax) @issue-3
Feature: DAX Financial Data Import
  As a financial analyst
  I want to import DAX financial data via CSV upload
  So that I can store and analyze historical financial metrics

  Background:
    Given the DAX API is available
    And the PostgreSQL database is clean

  @happy-path @import
  Scenario: Import CSV file with financial data
    When I upload a CSV file with the following content:
      """
      company,ticker,report_type,metric,year,value,currency
      Siemens AG,SIE,income,EBITDA,2025,15859000000.0,EUR
      Siemens AG,SIE,income,Net Income,2025,9620000000.0,EUR
      """
    Then the response status should be 200
    And the response should indicate 2 records imported
    And the database should contain 2 DAX records

  @happy-path @import
  Scenario: Import handles duplicate records with UPSERT
    Given the following DAX record exists:
      | company    | ticker | report_type | metric    | year | value         | currency |
      | Siemens AG | SIE    | income      | EBITDA    | 2025 | 15000000000.0 | EUR      |
    When I upload a CSV file with the following content:
      """
      company,ticker,report_type,metric,year,value,currency
      Siemens AG,SIE,income,EBITDA,2025,15859000000.0,EUR
      """
    Then the response status should be 200
    And the database should contain 1 DAX records
    And the EBITDA value for SIE 2025 should be 15859000000.0

  @validation @import
  Scenario: Reject CSV with missing required fields
    When I upload a CSV file with the following content:
      """
      company,ticker,metric,year,value
      Siemens AG,SIE,EBITDA,2025,15859000000.0
      """
    Then the response status should be 400
    And the error response should contain "missing required fields"

  @validation @import
  Scenario: Reject CSV with invalid year
    When I upload a CSV file with the following content:
      """
      company,ticker,report_type,metric,year,value,currency
      Siemens AG,SIE,income,EBITDA,invalid,15859000000.0,EUR
      """
    Then the response status should be 400
    And the error response should contain "invalid year"

  @happy-path @import
  Scenario: Import large CSV file with bulk insert
    When I upload a CSV file with 1000 records
    Then the response status should be 200
    And the response should indicate 1000 records imported
    And the database should contain 1000 DAX records
