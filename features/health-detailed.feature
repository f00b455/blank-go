# Issue: #13
# URL: https://github.com/f00b455/blank-go/issues/13
@pkg(handlers) @issue-13
Feature: Health Check with System Metrics
  As an API user or DevOps Engineer
  I want an extended health check endpoint with system metrics
  So that I can monitor the API and server status

  Background:
    Given the API server is running

  @happy-path @health-check
  Scenario: Get detailed health check with system metrics
    When I request detailed health check at "/api/v1/health/detailed"
    Then the response status should be 200
    And the response should contain field "status" with value "healthy"
    And the response should contain field "timestamp"
    And the response should contain field "version"
    And the response should contain field "uptime_seconds"
    And the response should contain system metrics with fields:
      | field            |
      | go_version       |
      | goroutines       |
      | memory_alloc_mb  |
      | memory_sys_mb    |
      | gc_runs          |
    And the response should contain checks with field "api" with value "ok"

  @validation @health-check
  Scenario: System metrics values are valid
    When I request detailed health check at "/api/v1/health/detailed"
    Then the response status should be 200
    And the system metric "goroutines" should be greater than 0
    And the system metric "memory_alloc_mb" should be greater than 0
    And the system metric "memory_sys_mb" should be greater than 0
    And the system metric "gc_runs" should be greater than or equal to 0
    And the uptime should be greater than 0

  @monitoring @health-check
  Scenario: Health check response includes version information
    When I request detailed health check at "/api/v1/health/detailed"
    Then the response status should be 200
    And the response should contain version information
    And the version should match semantic versioning pattern
