package features

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
	"github.com/f00b455/blank-go/internal/handlers"
	"github.com/f00b455/blank-go/pkg/dax"
	"github.com/gin-gonic/gin"
)

type daxContext struct {
	repo          dax.Repository
	service       *dax.Service
	handler       *handlers.DAXHandler
	router        *gin.Engine
	response      *httptest.ResponseRecorder
	lastResponse  map[string]interface{}
	errorResponse map[string]interface{}
}

func (ctx *daxContext) reset() {
	gin.SetMode(gin.TestMode)
	ctx.repo = dax.NewInMemoryRepository()
	ctx.service = dax.NewService(ctx.repo)
	ctx.handler = handlers.NewDAXHandler(ctx.service)
	ctx.setupRouter()

	ctx.response = nil
	ctx.lastResponse = nil
	ctx.errorResponse = nil
}

func (ctx *daxContext) setupRouter() {
	ctx.router = gin.New()
	api := ctx.router.Group("/api/v1/dax")
	{
		api.POST("/import", ctx.handler.ImportCSV)
		api.GET("", ctx.handler.GetByFilters)
		api.GET("/metrics", ctx.handler.GetMetrics)
	}
}

func (ctx *daxContext) cleanDatabase() {
	ctx.repo.DeleteAll()
}

func (ctx *daxContext) theDAXAPIIsAvailable() error {
	return nil
}

func (ctx *daxContext) thePostgreSQLDatabaseIsClean() error {
	ctx.cleanDatabase()
	return nil
}

func (ctx *daxContext) iUploadACSVFileWithTheFollowingContent(csvContent *godog.DocString) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "test.csv")
	if err != nil {
		return err
	}

	if _, err := io.WriteString(part, csvContent.Content); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "/api/v1/dax/import", body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)

	if ctx.response.Code == http.StatusOK {
		_ = json.Unmarshal(ctx.response.Body.Bytes(), &ctx.lastResponse)
	} else {
		_ = json.Unmarshal(ctx.response.Body.Bytes(), &ctx.errorResponse)
	}

	return nil
}

func (ctx *daxContext) theResponseStatusShouldBe(expectedStatus int) error {
	if ctx.response.Code != expectedStatus {
		return fmt.Errorf("expected status %d, got %d", expectedStatus, ctx.response.Code)
	}
	return nil
}

func (ctx *daxContext) theResponseShouldIndicateRecordsImported(expectedCount int) error {
	recordsImported, ok := ctx.lastResponse["records_imported"].(float64)
	if !ok {
		return fmt.Errorf("records_imported not found in response")
	}

	if int(recordsImported) != expectedCount {
		return fmt.Errorf("expected %d records imported, got %d", expectedCount, int(recordsImported))
	}

	return nil
}

func (ctx *daxContext) theDatabaseShouldContainDAXRecords(expectedCount int) error {
	count, err := ctx.repo.Count()
	if err != nil {
		return err
	}

	if count != expectedCount {
		return fmt.Errorf("expected %d records in database, got %d", expectedCount, count)
	}

	return nil
}

func (ctx *daxContext) theFollowingDAXRecordExists(table *godog.Table) error {
	if len(table.Rows) < 2 {
		return fmt.Errorf("table must have header and at least one data row")
	}

	headers := table.Rows[0].Cells
	dataRow := table.Rows[1].Cells

	record := dax.DAXRecord{}

	for i, header := range headers {
		value := dataRow[i].Value
		switch header.Value {
		case "company":
			record.Company = value
		case "ticker":
			record.Ticker = value
		case "report_type":
			record.ReportType = value
		case "metric":
			record.Metric = value
		case "year":
			year, _ := strconv.Atoi(value)
			record.Year = year
		case "value":
			val, _ := strconv.ParseFloat(value, 64)
			record.Value = &val
		case "currency":
			record.Currency = value
		}
	}

	return ctx.repo.Create(&record)
}

func (ctx *daxContext) theEBITDAValueForSIEShouldBe(year int, expectedValue float64) error {
	records, _, err := ctx.repo.FindByFilters("SIE", &year, 1, 100)
	if err != nil {
		return err
	}

	for _, record := range records {
		if record.Metric == "EBITDA" {
			if record.Value == nil || *record.Value != expectedValue {
				return fmt.Errorf("expected EBITDA value %f, got %f", expectedValue, *record.Value)
			}
			return nil
		}
	}

	return fmt.Errorf("EBITDA record not found for SIE in %d", year)
}

func (ctx *daxContext) theErrorResponseShouldContain(expectedError string) error {
	errorMsg, ok := ctx.errorResponse["error"].(string)
	if !ok {
		return fmt.Errorf("error message not found in response")
	}

	if !strings.Contains(errorMsg, expectedError) {
		return fmt.Errorf("expected error to contain '%s', got '%s'", expectedError, errorMsg)
	}

	return nil
}

func (ctx *daxContext) iUploadACSVFileWithRecords(recordCount int) error {
	csvContent := "company,ticker,report_type,metric,year,value,currency\n"
	for i := 0; i < recordCount; i++ {
		csvContent += fmt.Sprintf("Company%d,TICK%d,income,EBITDA,2025,%d.0,EUR\n", i, i%10, i*1000000)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "test.csv")
	if err != nil {
		return err
	}

	if _, err := io.WriteString(part, csvContent); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "/api/v1/dax/import", body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)

	if ctx.response.Code == http.StatusOK {
		_ = json.Unmarshal(ctx.response.Body.Bytes(), &ctx.lastResponse)
	}

	return nil
}

func (ctx *daxContext) theFollowingDAXRecordsExist(table *godog.Table) error {
	if len(table.Rows) < 2 {
		return fmt.Errorf("table must have header and at least one data row")
	}

	headers := table.Rows[0].Cells

	for i := 1; i < len(table.Rows); i++ {
		dataRow := table.Rows[i].Cells
		record := dax.DAXRecord{}

		for j, header := range headers {
			value := dataRow[j].Value
			switch header.Value {
			case "company":
				record.Company = value
			case "ticker":
				record.Ticker = value
			case "report_type":
				record.ReportType = value
			case "metric":
				record.Metric = value
			case "year":
				year, _ := strconv.Atoi(value)
				record.Year = year
			case "value":
				val, _ := strconv.ParseFloat(value, 64)
				record.Value = &val
			case "currency":
				record.Currency = value
			}
		}

		if err := ctx.repo.Create(&record); err != nil {
			return err
		}
	}

	return nil
}

func (ctx *daxContext) iRequestAllDAXRecordsWithPageAndLimit(page, limit int) error {
	url := fmt.Sprintf("/api/v1/dax?page=%d&limit=%d", page, limit)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)

	if ctx.response.Code == http.StatusOK {
		_ = json.Unmarshal(ctx.response.Body.Bytes(), &ctx.lastResponse)
	}

	return nil
}

func (ctx *daxContext) theResponseShouldContainRecords(expectedCount int) error {
	data, ok := ctx.lastResponse["data"].([]interface{})
	if !ok {
		return fmt.Errorf("data not found in response")
	}

	if len(data) != expectedCount {
		return fmt.Errorf("expected %d records, got %d", expectedCount, len(data))
	}

	return nil
}

func (ctx *daxContext) theResponseShouldIncludePaginationMetadata() error {
	_, ok := ctx.lastResponse["pagination"]
	if !ok {
		return fmt.Errorf("pagination metadata not found in response")
	}
	return nil
}

func (ctx *daxContext) iRequestDAXRecordsForTicker(ticker string) error {
	url := fmt.Sprintf("/api/v1/dax?ticker=%s", ticker)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)

	if ctx.response.Code == http.StatusOK {
		_ = json.Unmarshal(ctx.response.Body.Bytes(), &ctx.lastResponse)
	}

	return nil
}

func (ctx *daxContext) allRecordsShouldHaveTicker(expectedTicker string) error {
	data, ok := ctx.lastResponse["data"].([]interface{})
	if !ok {
		return fmt.Errorf("data not found in response")
	}

	for _, item := range data {
		record := item.(map[string]interface{})
		ticker, _ := record["ticker"].(string)
		if ticker != expectedTicker {
			return fmt.Errorf("expected ticker %s, got %s", expectedTicker, ticker)
		}
	}

	return nil
}

func (ctx *daxContext) iRequestDAXRecordsForTickerAndYear(ticker string, year int) error {
	url := fmt.Sprintf("/api/v1/dax?ticker=%s&year=%d", ticker, year)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)

	if ctx.response.Code == http.StatusOK {
		_ = json.Unmarshal(ctx.response.Body.Bytes(), &ctx.lastResponse)
	}

	return nil
}

func (ctx *daxContext) allRecordsShouldHaveTickerAndYear(expectedTicker string, expectedYear int) error {
	data, ok := ctx.lastResponse["data"].([]interface{})
	if !ok {
		return fmt.Errorf("data not found in response")
	}

	for _, item := range data {
		record := item.(map[string]interface{})
		ticker, _ := record["ticker"].(string)
		year, _ := record["year"].(float64)

		if ticker != expectedTicker || int(year) != expectedYear {
			return fmt.Errorf("expected ticker %s and year %d, got %s and %d",
				expectedTicker, expectedYear, ticker, int(year))
		}
	}

	return nil
}

func (ctx *daxContext) iRequestDAXRecordsForYear(year int) error {
	url := fmt.Sprintf("/api/v1/dax?year=%d", year)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)

	if ctx.response.Code == http.StatusOK {
		_ = json.Unmarshal(ctx.response.Body.Bytes(), &ctx.lastResponse)
	}

	return nil
}

func (ctx *daxContext) allRecordsShouldHaveYear(expectedYear int) error {
	data, ok := ctx.lastResponse["data"].([]interface{})
	if !ok {
		return fmt.Errorf("data not found in response")
	}

	for _, item := range data {
		record := item.(map[string]interface{})
		year, _ := record["year"].(float64)

		if int(year) != expectedYear {
			return fmt.Errorf("expected year %d, got %d", expectedYear, int(year))
		}
	}

	return nil
}

func (ctx *daxContext) iRequestAvailableMetricsForTicker(ticker string) error {
	url := fmt.Sprintf("/api/v1/dax/metrics?ticker=%s", ticker)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)

	if ctx.response.Code == http.StatusOK {
		_ = json.Unmarshal(ctx.response.Body.Bytes(), &ctx.lastResponse)
	}

	return nil
}

func (ctx *daxContext) theResponseShouldContainMetrics(expectedMetrics string) error {
	metrics, ok := ctx.lastResponse["metrics"].([]interface{})
	if !ok {
		return fmt.Errorf("metrics not found in response")
	}

	expectedList := strings.Split(expectedMetrics, ",")
	if len(metrics) != len(expectedList) {
		return fmt.Errorf("expected %d metrics, got %d", len(expectedList), len(metrics))
	}

	return nil
}

func (ctx *daxContext) theResponseShouldContainMetrics0(expectedCount int) error {
	metrics, ok := ctx.lastResponse["metrics"].([]interface{})
	if !ok {
		return fmt.Errorf("metrics not found in response")
	}

	if len(metrics) != expectedCount {
		return fmt.Errorf("expected %d metrics, got %d", expectedCount, len(metrics))
	}

	return nil
}

func InitializeDAXScenario(sc *godog.ScenarioContext) {
	ctx := &daxContext{}

	sc.Before(func(c context.Context, sc *godog.Scenario) (context.Context, error) {
		ctx.reset()
		return c, nil
	})

	sc.Step(`^the DAX API is available$`, ctx.theDAXAPIIsAvailable)
	sc.Step(`^the PostgreSQL database is clean$`, ctx.thePostgreSQLDatabaseIsClean)
	sc.Step(`^I upload a CSV file with the following content:$`, ctx.iUploadACSVFileWithTheFollowingContent)
	sc.Step(`^the response status should be (\d+)$`, ctx.theResponseStatusShouldBe)
	sc.Step(`^the response should indicate (\d+) records imported$`, ctx.theResponseShouldIndicateRecordsImported)
	sc.Step(`^the database should contain (\d+) DAX records$`, ctx.theDatabaseShouldContainDAXRecords)
	sc.Step(`^the following DAX record exists:$`, ctx.theFollowingDAXRecordExists)
	sc.Step(`^the EBITDA value for SIE (\d+) should be ([0-9.]+)$`, ctx.theEBITDAValueForSIEShouldBe)
	sc.Step(`^the error response should contain "([^"]*)"$`, ctx.theErrorResponseShouldContain)
	sc.Step(`^I upload a CSV file with (\d+) records$`, ctx.iUploadACSVFileWithRecords)
	sc.Step(`^the following DAX records exist:$`, ctx.theFollowingDAXRecordsExist)
	sc.Step(`^I request all DAX records with page (\d+) and limit (\d+)$`, ctx.iRequestAllDAXRecordsWithPageAndLimit)
	sc.Step(`^the response should contain (\d+) records$`, ctx.theResponseShouldContainRecords)
	sc.Step(`^the response should include pagination metadata$`, ctx.theResponseShouldIncludePaginationMetadata)
	sc.Step(`^I request DAX records for ticker "([^"]*)"$`, ctx.iRequestDAXRecordsForTicker)
	sc.Step(`^all records should have ticker "([^"]*)"$`, ctx.allRecordsShouldHaveTicker)
	sc.Step(`^I request DAX records for ticker "([^"]*)" and year (\d+)$`, ctx.iRequestDAXRecordsForTickerAndYear)
	sc.Step(`^all records should have ticker "([^"]*)" and year (\d+)$`, ctx.allRecordsShouldHaveTickerAndYear)
	sc.Step(`^I request DAX records for year (\d+)$`, ctx.iRequestDAXRecordsForYear)
	sc.Step(`^all records should have year (\d+)$`, ctx.allRecordsShouldHaveYear)
	sc.Step(`^I request available metrics for ticker "([^"]*)"$`, ctx.iRequestAvailableMetricsForTicker)
	sc.Step(`^the response should contain metrics "([^"]*)"$`, ctx.theResponseShouldContainMetrics)
	sc.Step(`^the response should contain (\d+) metrics$`, ctx.theResponseShouldContainMetrics0)
}
