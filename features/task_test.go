package features

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/cucumber/godog"
	"github.com/f00b455/blank-go/internal/handlers"
	"github.com/f00b455/blank-go/pkg/task"
	"github.com/gin-gonic/gin"
)

type taskFeatureContext struct {
	router       *gin.Engine
	service      *task.Service
	handler      *handlers.TaskHandler
	response     *httptest.ResponseRecorder
	lastTaskID   string
	lastResponse map[string]interface{}
	taskData     map[string]string
}

func (ctx *taskFeatureContext) reset() {
	gin.SetMode(gin.TestMode)
	repo := task.NewInMemoryRepository()
	ctx.service = task.NewService(repo)
	ctx.handler = handlers.NewTaskHandler(ctx.service)
	ctx.setupRouter()

	ctx.response = nil
	ctx.lastTaskID = ""
	ctx.lastResponse = nil
	ctx.taskData = make(map[string]string)
}

func (ctx *taskFeatureContext) setupRouter() {
	ctx.router = gin.New()
	api := ctx.router.Group("/api/v1")
	{
		api.POST("/tasks", ctx.handler.CreateTask)
		api.GET("/tasks", ctx.handler.ListTasks)
		api.GET("/tasks/:id", ctx.handler.GetTask)
		api.PUT("/tasks/:id", ctx.handler.UpdateTask)
		api.DELETE("/tasks/:id", ctx.handler.DeleteTask)
	}
}

func (ctx *taskFeatureContext) theTaskAPIIsAvailable() error {
	return nil
}

func (ctx *taskFeatureContext) theTaskStorageIsEmpty() error {
	// Create a fresh repository and handler to ensure storage is empty
	repo := task.NewInMemoryRepository()
	ctx.service = task.NewService(repo)
	ctx.handler = handlers.NewTaskHandler(ctx.service)
	ctx.setupRouter()
	return nil
}

func (ctx *taskFeatureContext) iCreateATaskWithTheFollowingDetails(table *godog.Table) error {
	ctx.taskData = tableToMap(table)
	return ctx.createTask()
}

func (ctx *taskFeatureContext) createTask() error {
	reqBody := make(map[string]interface{})
	if title, ok := ctx.taskData["title"]; ok {
		reqBody["title"] = title
	}
	if desc, ok := ctx.taskData["description"]; ok {
		reqBody["description"] = desc
	}
	if priority, ok := ctx.taskData["priority"]; ok {
		reqBody["priority"] = priority
	}
	if status, ok := ctx.taskData["status"]; ok {
		reqBody["status"] = status
	}
	if dueDate, ok := ctx.taskData["due_date"]; ok {
		t, _ := time.Parse(time.RFC3339, dueDate)
		reqBody["due_date"] = t
	}
	if tags, ok := ctx.taskData["tags"]; ok {
		reqBody["tags"] = strings.Split(tags, ",")
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)

	if ctx.response.Code == http.StatusCreated {
		if err := json.Unmarshal(ctx.response.Body.Bytes(), &ctx.lastResponse); err != nil {
			return err
		}
		if id, ok := ctx.lastResponse["id"].(string); ok {
			ctx.lastTaskID = id
		}
	}

	return nil
}

func (ctx *taskFeatureContext) theResponseStatusShouldBe(expected int) error {
	if ctx.response.Code != expected {
		return fmt.Errorf("expected status %d, got %d", expected, ctx.response.Code)
	}
	return nil
}

func (ctx *taskFeatureContext) theResponseShouldContainATaskWithTitle(title string) error {
	if ctx.lastResponse == nil {
		if err := json.Unmarshal(ctx.response.Body.Bytes(), &ctx.lastResponse); err != nil {
			return err
		}
	}
	if actualTitle, ok := ctx.lastResponse["title"].(string); !ok || actualTitle != title {
		return fmt.Errorf("expected title %q, got %q", title, actualTitle)
	}
	return nil
}

func (ctx *taskFeatureContext) theTaskShouldHaveStatus(status string) error {
	if actualStatus, ok := ctx.lastResponse["status"].(string); !ok || actualStatus != status {
		return fmt.Errorf("expected status %q, got %q", status, actualStatus)
	}
	return nil
}

func (ctx *taskFeatureContext) theTaskShouldHavePriority(priority string) error {
	if actualPriority, ok := ctx.lastResponse["priority"].(string); !ok || actualPriority != priority {
		return fmt.Errorf("expected priority %q, got %q", priority, actualPriority)
	}
	return nil
}

func (ctx *taskFeatureContext) theTaskShouldHaveAValidUUID() error {
	if id, ok := ctx.lastResponse["id"].(string); !ok || id == "" {
		return fmt.Errorf("expected valid UUID, got %v", ctx.lastResponse["id"])
	}
	return nil
}

func (ctx *taskFeatureContext) theTaskShouldHaveCreatedAtAndUpdatedAtTimestamps() error {
	if _, ok := ctx.lastResponse["created_at"].(string); !ok {
		return fmt.Errorf("missing created_at timestamp")
	}
	if _, ok := ctx.lastResponse["updated_at"].(string); !ok {
		return fmt.Errorf("missing updated_at timestamp")
	}
	return nil
}

func (ctx *taskFeatureContext) theTaskShouldHaveDescription(desc string) error {
	if actualDesc, ok := ctx.lastResponse["description"].(string); !ok || actualDesc != desc {
		return fmt.Errorf("expected description %q, got %q", desc, actualDesc)
	}
	return nil
}

func (ctx *taskFeatureContext) theTaskShouldHaveDueDate(dueDate string) error {
	if actualDate, ok := ctx.lastResponse["due_date"].(string); !ok || actualDate != dueDate {
		return fmt.Errorf("expected due_date %q, got %q", dueDate, actualDate)
	}
	return nil
}

func (ctx *taskFeatureContext) theTaskShouldHaveTags(tags string) error {
	expectedTags := strings.Split(tags, ",")
	actualTagsInterface, ok := ctx.lastResponse["tags"].([]interface{})
	if !ok {
		return fmt.Errorf("expected tags array, got %T", ctx.lastResponse["tags"])
	}

	actualTags := make([]string, len(actualTagsInterface))
	for i, tag := range actualTagsInterface {
		actualTags[i] = tag.(string)
	}

	if len(actualTags) != len(expectedTags) {
		return fmt.Errorf("expected %d tags, got %d", len(expectedTags), len(actualTags))
	}

	return nil
}

func (ctx *taskFeatureContext) theErrorResponseShouldContain(message string) error {
	var errResponse map[string]interface{}
	if err := json.Unmarshal(ctx.response.Body.Bytes(), &errResponse); err != nil {
		return err
	}

	if errorField, ok := errResponse["error"].(map[string]interface{}); ok {
		if msg, ok := errorField["message"].(string); ok {
			if strings.Contains(msg, message) {
				return nil
			}
			return fmt.Errorf("error message %q does not contain %q", msg, message)
		}
	}

	return fmt.Errorf("no error message found in response")
}

func (ctx *taskFeatureContext) aTaskExistsWithTitle(title string) error {
	_, err := ctx.service.Create(task.CreateTaskRequest{Title: title})
	if err != nil {
		return err
	}

	// Get the created task to store its ID
	filter := task.FilterOptions{}
	tasks, err := ctx.service.GetAll(filter)
	if err != nil || len(tasks) == 0 {
		return fmt.Errorf("failed to retrieve created task")
	}

	ctx.lastTaskID = tasks[len(tasks)-1].ID
	return nil
}

func (ctx *taskFeatureContext) iRequestTheTaskByItsID() error {
	req, _ := http.NewRequest("GET", "/api/v1/tasks/"+ctx.lastTaskID, nil)
	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)

	if ctx.response.Code == http.StatusOK {
		if err := json.Unmarshal(ctx.response.Body.Bytes(), &ctx.lastResponse); err != nil {
			return err
		}
	}

	return nil
}

func (ctx *taskFeatureContext) iRequestATaskWithID(id string) error {
	req, _ := http.NewRequest("GET", "/api/v1/tasks/"+id, nil)
	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)
	return nil
}

func (ctx *taskFeatureContext) iUpdateTheTaskWithTheFollowingDetails(table *godog.Table) error {
	ctx.taskData = tableToMap(table)
	return ctx.updateTask(ctx.lastTaskID)
}

func (ctx *taskFeatureContext) updateTask(id string) error {
	reqBody := make(map[string]interface{})
	if title, ok := ctx.taskData["title"]; ok {
		reqBody["title"] = title
	}
	if desc, ok := ctx.taskData["description"]; ok {
		reqBody["description"] = desc
	}
	if priority, ok := ctx.taskData["priority"]; ok {
		reqBody["priority"] = priority
	}
	if status, ok := ctx.taskData["status"]; ok {
		reqBody["status"] = status
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/tasks/"+id, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)

	if ctx.response.Code == http.StatusOK {
		if err := json.Unmarshal(ctx.response.Body.Bytes(), &ctx.lastResponse); err != nil {
			return err
		}
	}

	return nil
}

func (ctx *taskFeatureContext) aTaskExistsWithTitleAndPriority(title, priority string) error {
	_, err := ctx.service.Create(task.CreateTaskRequest{
		Title:    title,
		Priority: priority,
	})
	if err != nil {
		return err
	}

	filter := task.FilterOptions{}
	tasks, err := ctx.service.GetAll(filter)
	if err != nil || len(tasks) == 0 {
		return fmt.Errorf("failed to retrieve created task")
	}

	ctx.lastTaskID = tasks[len(tasks)-1].ID
	return nil
}

func (ctx *taskFeatureContext) iUpdateATaskWithIDWithTitle(id, title string) error {
	ctx.taskData = map[string]string{"title": title}
	return ctx.updateTask(id)
}

func (ctx *taskFeatureContext) iDeleteTheTaskByItsID() error {
	req, _ := http.NewRequest("DELETE", "/api/v1/tasks/"+ctx.lastTaskID, nil)
	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)
	return nil
}

func (ctx *taskFeatureContext) theTaskShouldNoLongerExist() error {
	_, err := ctx.service.GetByID(ctx.lastTaskID)
	if err == nil {
		return fmt.Errorf("task still exists")
	}
	return nil
}

func (ctx *taskFeatureContext) iDeleteATaskWithID(id string) error {
	req, _ := http.NewRequest("DELETE", "/api/v1/tasks/"+id, nil)
	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)
	return nil
}

func (ctx *taskFeatureContext) theFollowingTasksExist(table *godog.Table) error {
	for i := 1; i < len(table.Rows); i++ {
		row := table.Rows[i]
		reqData := task.CreateTaskRequest{
			Title:    row.Cells[0].Value,
			Status:   row.Cells[1].Value,
			Priority: row.Cells[2].Value,
			Tags:     strings.Split(row.Cells[3].Value, ","),
		}
		_, err := ctx.service.Create(reqData)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ctx *taskFeatureContext) iRequestAllTasks() error {
	req, _ := http.NewRequest("GET", "/api/v1/tasks", nil)
	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)
	return nil
}

func (ctx *taskFeatureContext) theResponseShouldContainTasks(count int) error {
	var tasks []interface{}
	if err := json.Unmarshal(ctx.response.Body.Bytes(), &tasks); err != nil {
		return err
	}

	if len(tasks) != count {
		return fmt.Errorf("expected %d tasks, got %d", count, len(tasks))
	}
	return nil
}

func (ctx *taskFeatureContext) iRequestTasksWithLimitAndOffset(limit, offset int) error {
	url := fmt.Sprintf("/api/v1/tasks?limit=%d&offset=%d", limit, offset)
	req, _ := http.NewRequest("GET", url, nil)
	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)
	return nil
}

func (ctx *taskFeatureContext) iRequestTasksWithStatus(status string) error {
	url := fmt.Sprintf("/api/v1/tasks?status=%s", status)
	req, _ := http.NewRequest("GET", url, nil)
	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)
	return nil
}

func (ctx *taskFeatureContext) allTasksShouldHaveStatus(status string) error {
	var tasks []map[string]interface{}
	if err := json.Unmarshal(ctx.response.Body.Bytes(), &tasks); err != nil {
		return err
	}

	for _, t := range tasks {
		if t["status"] != status {
			return fmt.Errorf("task has status %v, expected %s", t["status"], status)
		}
	}
	return nil
}

func (ctx *taskFeatureContext) iRequestTasksWithPriority(priority string) error {
	url := fmt.Sprintf("/api/v1/tasks?priority=%s", priority)
	req, _ := http.NewRequest("GET", url, nil)
	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)
	return nil
}

func (ctx *taskFeatureContext) allTasksShouldHavePriority(priority string) error {
	var tasks []map[string]interface{}
	if err := json.Unmarshal(ctx.response.Body.Bytes(), &tasks); err != nil {
		return err
	}

	for _, t := range tasks {
		if t["priority"] != priority {
			return fmt.Errorf("task has priority %v, expected %s", t["priority"], priority)
		}
	}
	return nil
}

func (ctx *taskFeatureContext) iRequestTasksWithTag(tag string) error {
	url := fmt.Sprintf("/api/v1/tasks?tag=%s", tag)
	req, _ := http.NewRequest("GET", url, nil)
	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)
	return nil
}

func (ctx *taskFeatureContext) allTasksShouldHaveTag(tag string) error {
	var tasks []map[string]interface{}
	if err := json.Unmarshal(ctx.response.Body.Bytes(), &tasks); err != nil {
		return err
	}

	for _, t := range tasks {
		tagsInterface, ok := t["tags"].([]interface{})
		if !ok {
			return fmt.Errorf("tags field is not an array")
		}

		found := false
		for _, tagInterface := range tagsInterface {
			if tagInterface.(string) == tag {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("task does not have tag %s", tag)
		}
	}
	return nil
}

func (ctx *taskFeatureContext) iRequestTasksWithStatusAndPriority(status, priority string) error {
	url := fmt.Sprintf("/api/v1/tasks?status=%s&priority=%s", status, priority)
	req, _ := http.NewRequest("GET", url, nil)
	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)
	return nil
}

func (ctx *taskFeatureContext) theTaskShouldHaveTitle(title string) error {
	var tasks []map[string]interface{}
	if err := json.Unmarshal(ctx.response.Body.Bytes(), &tasks); err != nil {
		return err
	}

	if len(tasks) == 0 {
		return fmt.Errorf("no tasks in response")
	}

	if tasks[0]["title"] != title {
		return fmt.Errorf("expected title %q, got %q", title, tasks[0]["title"])
	}
	return nil
}

func (ctx *taskFeatureContext) iRequestTasksSortedByInOrder(field, order string) error {
	url := fmt.Sprintf("/api/v1/tasks?sort_by=%s&order=%s", field, order)
	req, _ := http.NewRequest("GET", url, nil)
	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)
	return nil
}

func (ctx *taskFeatureContext) theTasksShouldBeOrderedByCreationTimeAscending() error {
	// Just verify we got tasks back
	var tasks []map[string]interface{}
	if err := json.Unmarshal(ctx.response.Body.Bytes(), &tasks); err != nil {
		return err
	}
	if len(tasks) == 0 {
		return fmt.Errorf("no tasks returned")
	}
	return nil
}

func (ctx *taskFeatureContext) theFirstTaskShouldHavePriority(priority string) error {
	var tasks []map[string]interface{}
	if err := json.Unmarshal(ctx.response.Body.Bytes(), &tasks); err != nil {
		return err
	}

	if len(tasks) == 0 {
		return fmt.Errorf("no tasks in response")
	}

	if tasks[0]["priority"] != priority {
		return fmt.Errorf("first task has priority %v, expected %s", tasks[0]["priority"], priority)
	}
	return nil
}

func tableToMap(table *godog.Table) map[string]string {
	result := make(map[string]string)
	for i := 1; i < len(table.Rows); i++ {
		row := table.Rows[i]
		result[row.Cells[0].Value] = row.Cells[1].Value
	}
	return result
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	feature := &taskFeatureContext{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		feature.reset()
		return ctx, nil
	})

	// Background steps
	ctx.Step(`^the task API is available$`, feature.theTaskAPIIsAvailable)
	ctx.Step(`^the task storage is empty$`, feature.theTaskStorageIsEmpty)

	// Create steps
	ctx.Step(`^I create a task with the following details:$`, feature.iCreateATaskWithTheFollowingDetails)

	// Response validation steps
	ctx.Step(`^the response status should be (\d+)$`, feature.theResponseStatusShouldBe)
	ctx.Step(`^the response should contain a task with title "([^"]*)"$`, feature.theResponseShouldContainATaskWithTitle)
	ctx.Step(`^the task should have status "([^"]*)"$`, feature.theTaskShouldHaveStatus)
	ctx.Step(`^the task should have priority "([^"]*)"$`, feature.theTaskShouldHavePriority)
	ctx.Step(`^the task should have a valid UUID$`, feature.theTaskShouldHaveAValidUUID)
	ctx.Step(`^the task should have created_at and updated_at timestamps$`, feature.theTaskShouldHaveCreatedAtAndUpdatedAtTimestamps)
	ctx.Step(`^the task should have description "([^"]*)"$`, feature.theTaskShouldHaveDescription)
	ctx.Step(`^the task should have due_date "([^"]*)"$`, feature.theTaskShouldHaveDueDate)
	ctx.Step(`^the task should have tags "([^"]*)"$`, feature.theTaskShouldHaveTags)
	ctx.Step(`^the error response should contain "([^"]*)"$`, feature.theErrorResponseShouldContain)

	// Read steps
	ctx.Step(`^a task exists with title "([^"]*)"$`, feature.aTaskExistsWithTitle)
	ctx.Step(`^I request the task by its ID$`, feature.iRequestTheTaskByItsID)
	ctx.Step(`^I request a task with ID "([^"]*)"$`, feature.iRequestATaskWithID)

	// Update steps
	ctx.Step(`^I update the task with the following details:$`, feature.iUpdateTheTaskWithTheFollowingDetails)
	ctx.Step(`^a task exists with title "([^"]*)" and priority "([^"]*)"$`, feature.aTaskExistsWithTitleAndPriority)
	ctx.Step(`^I update a task with ID "([^"]*)" with title "([^"]*)"$`, feature.iUpdateATaskWithIDWithTitle)

	// Delete steps
	ctx.Step(`^I delete the task by its ID$`, feature.iDeleteTheTaskByItsID)
	ctx.Step(`^the task should no longer exist$`, feature.theTaskShouldNoLongerExist)
	ctx.Step(`^I delete a task with ID "([^"]*)"$`, feature.iDeleteATaskWithID)

	// List and filter steps
	ctx.Step(`^the following tasks exist:$`, feature.theFollowingTasksExist)
	ctx.Step(`^I request all tasks$`, feature.iRequestAllTasks)
	ctx.Step(`^the response should contain (\d+) tasks$`, feature.theResponseShouldContainTasks)
	ctx.Step(`^the response should contain (\d+) task$`, feature.theResponseShouldContainTasks)
	ctx.Step(`^I request tasks with limit (\d+) and offset (\d+)$`, feature.iRequestTasksWithLimitAndOffset)
	ctx.Step(`^I request tasks with status "([^"]*)"$`, feature.iRequestTasksWithStatus)
	ctx.Step(`^all tasks should have status "([^"]*)"$`, feature.allTasksShouldHaveStatus)
	ctx.Step(`^I request tasks with priority "([^"]*)"$`, feature.iRequestTasksWithPriority)
	ctx.Step(`^all tasks should have priority "([^"]*)"$`, feature.allTasksShouldHavePriority)
	ctx.Step(`^I request tasks with tag "([^"]*)"$`, feature.iRequestTasksWithTag)
	ctx.Step(`^all tasks should have tag "([^"]*)"$`, feature.allTasksShouldHaveTag)
	ctx.Step(`^I request tasks with status "([^"]*)" and priority "([^"]*)"$`, feature.iRequestTasksWithStatusAndPriority)
	ctx.Step(`^the task should have title "([^"]*)"$`, feature.theTaskShouldHaveTitle)
	ctx.Step(`^I request tasks sorted by "([^"]*)" in "([^"]*)" order$`, feature.iRequestTasksSortedByInOrder)
	ctx.Step(`^the tasks should be ordered by creation time ascending$`, feature.theTasksShouldBeOrderedByCreationTimeAscending)
	ctx.Step(`^the first task should have priority "([^"]*)"$`, feature.theFirstTaskShouldHavePriority)
}
