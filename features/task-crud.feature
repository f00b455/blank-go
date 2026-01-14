# Issue: #1
# URL: https://github.com/f00b455/blank-go/issues/1
@pkg(task) @issue-1
Feature: Task CRUD Operations
  As a developer/API user
  I want to perform CRUD operations on tasks
  So that I can manage and track my tasks effectively

  Background:
    Given the task API is available
    And the task storage is empty

  @happy-path @create
  Scenario: Create a new task with required fields only
    When I create a task with the following details:
      | field | value           |
      | title | Complete report |
    Then the response status should be 201
    And the response should contain a task with title "Complete report"
    And the task should have status "pending"
    And the task should have priority "medium"
    And the task should have a valid UUID
    And the task should have created_at and updated_at timestamps

  @happy-path @create
  Scenario: Create a new task with all fields
    When I create a task with the following details:
      | field       | value                    |
      | title       | Review pull request      |
      | description | Check code quality       |
      | priority    | high                     |
      | status      | in_progress              |
      | due_date    | 2026-01-20T10:00:00Z     |
      | tags        | code-review,urgent       |
    Then the response status should be 201
    And the response should contain a task with title "Review pull request"
    And the task should have description "Check code quality"
    And the task should have priority "high"
    And the task should have status "in_progress"
    And the task should have due_date "2026-01-20T10:00:00Z"
    And the task should have tags "code-review,urgent"

  @validation @create
  Scenario: Fail to create task without title
    When I create a task with the following details:
      | field       | value          |
      | description | No title given |
    Then the response status should be 400
    And the error response should contain "title is required"

  @validation @create
  Scenario: Fail to create task with invalid priority
    When I create a task with the following details:
      | field    | value        |
      | title    | Test task    |
      | priority | super-urgent |
    Then the response status should be 400
    And the error response should contain "invalid priority"

  @validation @create
  Scenario: Fail to create task with invalid status
    When I create a task with the following details:
      | field  | value     |
      | title  | Test task |
      | status | done      |
    Then the response status should be 400
    And the error response should contain "invalid status"

  @happy-path @read
  Scenario: Get a single task by ID
    Given a task exists with title "Existing task"
    When I request the task by its ID
    Then the response status should be 200
    And the response should contain a task with title "Existing task"

  @error @read
  Scenario: Get non-existent task returns 404
    When I request a task with ID "non-existent-id"
    Then the response status should be 404
    And the error response should contain "task not found"

  @happy-path @update
  Scenario: Update task fields
    Given a task exists with title "Old title"
    When I update the task with the following details:
      | field       | value                 |
      | title       | Updated title         |
      | description | Updated description   |
      | priority    | high                  |
      | status      | completed             |
    Then the response status should be 200
    And the response should contain a task with title "Updated title"
    And the task should have description "Updated description"
    And the task should have priority "high"
    And the task should have status "completed"

  @happy-path @update
  Scenario: Partial update of task
    Given a task exists with title "Original" and priority "low"
    When I update the task with the following details:
      | field  | value      |
      | status | completed  |
    Then the response status should be 200
    And the response should contain a task with title "Original"
    And the task should have priority "low"
    And the task should have status "completed"

  @error @update
  Scenario: Update non-existent task returns 404
    When I update a task with ID "non-existent-id" with title "New title"
    Then the response status should be 404
    And the error response should contain "task not found"

  @happy-path @delete
  Scenario: Delete an existing task
    Given a task exists with title "To be deleted"
    When I delete the task by its ID
    Then the response status should be 200
    And the task should no longer exist

  @error @delete
  Scenario: Delete non-existent task returns 404
    When I delete a task with ID "non-existent-id"
    Then the response status should be 404
    And the error response should contain "task not found"
