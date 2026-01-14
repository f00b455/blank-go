# Issue: #1
# URL: https://github.com/f00b455/blank-go/issues/1
@pkg(task) @issue-1
Feature: Task Listing and Filtering
  As a developer/API user
  I want to list and filter tasks
  So that I can find and organize my tasks efficiently

  Background:
    Given the task API is available
    And the following tasks exist:
      | title          | status      | priority | tags          |
      | Task 1         | pending     | high     | work,urgent   |
      | Task 2         | in_progress | medium   | personal      |
      | Task 3         | completed   | low      | work          |
      | Task 4         | pending     | high     | urgent        |
      | Task 5         | in_progress | high     | work,personal |

  @happy-path @list
  Scenario: Get all tasks
    When I request all tasks
    Then the response status should be 200
    And the response should contain 5 tasks

  @pagination @list
  Scenario: Get tasks with pagination
    When I request tasks with limit 2 and offset 0
    Then the response status should be 200
    And the response should contain 2 tasks
    When I request tasks with limit 2 and offset 2
    Then the response status should be 200
    And the response should contain 2 tasks

  @filter @list
  Scenario: Filter tasks by status
    When I request tasks with status "pending"
    Then the response status should be 200
    And the response should contain 2 tasks
    And all tasks should have status "pending"

  @filter @list
  Scenario: Filter tasks by priority
    When I request tasks with priority "high"
    Then the response status should be 200
    And the response should contain 3 tasks
    And all tasks should have priority "high"

  @filter @list
  Scenario: Filter tasks by tag
    When I request tasks with tag "work"
    Then the response status should be 200
    And the response should contain 3 tasks
    And all tasks should have tag "work"

  @filter @list
  Scenario: Filter tasks by multiple criteria
    When I request tasks with status "in_progress" and priority "high"
    Then the response status should be 200
    And the response should contain 1 task
    And the task should have title "Task 5"

  @sorting @list
  Scenario: Sort tasks by created_at ascending
    When I request tasks sorted by "created_at" in "asc" order
    Then the response status should be 200
    And the tasks should be ordered by creation time ascending

  @sorting @list
  Scenario: Sort tasks by priority descending
    When I request tasks sorted by "priority" in "desc" order
    Then the response status should be 200
    And the first task should have priority "high"

  @edge-case @list
  Scenario: Get tasks when storage is empty
    Given the task storage is empty
    When I request all tasks
    Then the response status should be 200
    And the response should contain 0 tasks
