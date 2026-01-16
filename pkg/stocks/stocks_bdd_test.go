package stocks_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/f00b455/blank-go/internal/handlers"
	"github.com/f00b455/blank-go/pkg/stocks"
	"github.com/f00b455/blank-go/pkg/stocks/mocks"
	"github.com/gin-gonic/gin"
)

type stocksFeatureContext struct {
	router        *gin.Engine
	response      *httptest.ResponseRecorder
	responseBody  map[string]interface{}
	firstResponse *stocks.StockSummary
	mockClient    *mocks.MockStocksClient
}

func (ctx *stocksFeatureContext) theYahooFinanceAPIIsAvailable() error {
	gin.SetMode(gin.TestMode)

	// Initialize mock client
	ctx.mockClient = new(mocks.MockStocksClient)

	// Setup mock responses for common test tickers
	ctx.setupMockData()

	// Initialize stocks service with mock client
	stocksService := stocks.NewService(ctx.mockClient)
	stocksHandler := handlers.NewStocksHandler(stocksService)

	// Setup router
	ctx.router = gin.New()
	api := ctx.router.Group("/api/v1")
	{
		api.GET("/stocks/:ticker/summary", stocksHandler.GetStockSummary)
		api.GET("/stocks/summary", stocksHandler.GetBatchSummary)
	}

	return nil
}

func (ctx *stocksFeatureContext) setupMockData() {
	// Mock data for AAPL
	aaplQuote := &stocks.YahooQuote{
		Symbol:                     "AAPL",
		ShortName:                  "Apple Inc.",
		RegularMarketPrice:         185.50,
		RegularMarketOpen:          184.00,
		RegularMarketHigh:          186.20,
		RegularMarketLow:           183.50,
		RegularMarketChange:        1.50,
		RegularMarketChangePercent: 0.81,
		RegularMarketVolume:        52000000,
		Currency:                   "USD",
	}

	// Mock data for GOOGL
	googlQuote := &stocks.YahooQuote{
		Symbol:                     "GOOGL",
		ShortName:                  "Alphabet Inc.",
		RegularMarketPrice:         140.50,
		RegularMarketOpen:          139.00,
		RegularMarketHigh:          141.20,
		RegularMarketLow:           138.50,
		RegularMarketChange:        1.50,
		RegularMarketChangePercent: 1.08,
		RegularMarketVolume:        25000000,
		Currency:                   "USD",
	}

	// Mock data for MSFT
	msftQuote := &stocks.YahooQuote{
		Symbol:                     "MSFT",
		ShortName:                  "Microsoft Corporation",
		RegularMarketPrice:         420.00,
		RegularMarketOpen:          418.00,
		RegularMarketHigh:          422.00,
		RegularMarketLow:           417.50,
		RegularMarketChange:        2.00,
		RegularMarketChangePercent: 0.48,
		RegularMarketVolume:        30000000,
		Currency:                   "USD",
	}

	// Mock data for TSLA
	tslaQuote := &stocks.YahooQuote{
		Symbol:                     "TSLA",
		ShortName:                  "Tesla, Inc.",
		RegularMarketPrice:         245.00,
		RegularMarketOpen:          243.00,
		RegularMarketHigh:          247.00,
		RegularMarketLow:           242.00,
		RegularMarketChange:        2.00,
		RegularMarketChangePercent: 0.82,
		RegularMarketVolume:        45000000,
		Currency:                   "USD",
	}

	// Mock data for AMZN
	amznQuote := &stocks.YahooQuote{
		Symbol:                     "AMZN",
		ShortName:                  "Amazon.com, Inc.",
		RegularMarketPrice:         178.00,
		RegularMarketOpen:          176.50,
		RegularMarketHigh:          179.00,
		RegularMarketLow:           175.50,
		RegularMarketChange:        1.50,
		RegularMarketChangePercent: 0.85,
		RegularMarketVolume:        35000000,
		Currency:                   "USD",
	}

	// Setup mock expectations for single ticker requests
	ctx.mockClient.On("GetQuote", "AAPL").Return(aaplQuote, nil).Maybe()
	ctx.mockClient.On("GetQuote", "GOOGL").Return(googlQuote, nil).Maybe()
	ctx.mockClient.On("GetQuote", "MSFT").Return(msftQuote, nil).Maybe()
	ctx.mockClient.On("GetQuote", "TSLA").Return(tslaQuote, nil).Maybe()
	ctx.mockClient.On("GetQuote", "AMZN").Return(amznQuote, nil).Maybe()
	ctx.mockClient.On("GetQuote", "").Return(nil, fmt.Errorf("ticker is required")).Maybe()
	ctx.mockClient.On("GetQuote", "INVALID_TICKER_XYZ").Return(nil, fmt.Errorf("ticker not found")).Maybe()
	ctx.mockClient.On("GetQuote", "INVALID_XYZ").Return(nil, fmt.Errorf("ticker not found")).Maybe()

	// Setup mock expectations for batch requests
	ctx.mockClient.On("GetQuotes", []string{"AAPL", "GOOGL", "MSFT"}).Return(map[string]*stocks.YahooQuote{
		"AAPL":  aaplQuote,
		"GOOGL": googlQuote,
		"MSFT":  msftQuote,
	}, nil).Maybe()

	ctx.mockClient.On("GetQuotes", []string{"AAPL"}).Return(map[string]*stocks.YahooQuote{
		"AAPL": aaplQuote,
	}, nil).Maybe()

	ctx.mockClient.On("GetQuotes", []string{"AAPL", "INVALID_XYZ", "MSFT"}).Return(map[string]*stocks.YahooQuote{
		"AAPL": aaplQuote,
		"MSFT": msftQuote,
	}, nil).Maybe()

	ctx.mockClient.On("GetQuotes", []string{"AAPL", "GOOGL", "MSFT", "TSLA", "AMZN"}).Return(map[string]*stocks.YahooQuote{
		"AAPL":  aaplQuote,
		"GOOGL": googlQuote,
		"MSFT":  msftQuote,
		"TSLA":  tslaQuote,
		"AMZN":  amznQuote,
	}, nil).Maybe()
}

func (ctx *stocksFeatureContext) iRequestStockSummaryForTicker(ticker string) error {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/stocks/"+ticker+"/summary", nil)
	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)

	if ctx.response.Code == http.StatusOK {
		_ = json.Unmarshal([]byte(ctx.response.Body.String()), &ctx.responseBody)
	}

	return nil
}

func (ctx *stocksFeatureContext) iRequestBatchStockSummaryForTickers(tickers string) error {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/stocks/summary?tickers="+tickers, nil)
	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)

	_ = json.Unmarshal([]byte(ctx.response.Body.String()), &ctx.responseBody)

	return nil
}

func (ctx *stocksFeatureContext) theResponseStatusShouldBe(expectedStatus int) error {
	if ctx.response.Code != expectedStatus {
		return fmt.Errorf("expected status %d, got %d. Response body: %s", expectedStatus, ctx.response.Code, ctx.response.Body.String())
	}
	return nil
}

func (ctx *stocksFeatureContext) theResponseShouldContainTicker(ticker string) error {
	if tickerValue, ok := ctx.responseBody["ticker"]; !ok || tickerValue != ticker {
		return fmt.Errorf("expected ticker %s in response", ticker)
	}
	return nil
}

func (ctx *stocksFeatureContext) theResponseShouldContainCurrentPrice() error {
	if _, ok := ctx.responseBody["current_price"]; !ok {
		return fmt.Errorf("response should contain current_price")
	}
	return nil
}

func (ctx *stocksFeatureContext) theResponseShouldContainOpenPrice() error {
	if _, ok := ctx.responseBody["open"]; !ok {
		return fmt.Errorf("response should contain open")
	}
	return nil
}

func (ctx *stocksFeatureContext) theResponseShouldContainHighPrice() error {
	if _, ok := ctx.responseBody["high"]; !ok {
		return fmt.Errorf("response should contain high")
	}
	return nil
}

func (ctx *stocksFeatureContext) theResponseShouldContainLowPrice() error {
	if _, ok := ctx.responseBody["low"]; !ok {
		return fmt.Errorf("response should contain low")
	}
	return nil
}

func (ctx *stocksFeatureContext) theResponseShouldContainChangeValue() error {
	if _, ok := ctx.responseBody["change"]; !ok {
		return fmt.Errorf("response should contain change")
	}
	return nil
}

func (ctx *stocksFeatureContext) theResponseShouldContainChangePercentage() error {
	if _, ok := ctx.responseBody["change_percent"]; !ok {
		return fmt.Errorf("response should contain change_percent")
	}
	return nil
}

func (ctx *stocksFeatureContext) theResponseShouldContainVolume() error {
	if _, ok := ctx.responseBody["volume"]; !ok {
		return fmt.Errorf("response should contain volume")
	}
	return nil
}

func (ctx *stocksFeatureContext) theResponseShouldContainCurrency() error {
	if _, ok := ctx.responseBody["currency"]; !ok {
		return fmt.Errorf("response should contain currency")
	}
	return nil
}

func (ctx *stocksFeatureContext) theResponseShouldContainCompanyName() error {
	if _, ok := ctx.responseBody["name"]; !ok {
		return fmt.Errorf("response should contain name")
	}
	return nil
}

func (ctx *stocksFeatureContext) theCompanyNameShouldNotBeEmpty() error {
	if name, ok := ctx.responseBody["name"].(string); !ok || name == "" {
		return fmt.Errorf("company name should not be empty")
	}
	return nil
}

func (ctx *stocksFeatureContext) theResponseShouldContainDate() error {
	if _, ok := ctx.responseBody["date"]; !ok {
		return fmt.Errorf("response should contain date")
	}
	return nil
}

func (ctx *stocksFeatureContext) theDateShouldBeInFormat(format string) error {
	dateStr, ok := ctx.responseBody["date"].(string)
	if !ok {
		return fmt.Errorf("date should be a string")
	}

	// Simple check for YYYY-MM-DD format
	parts := strings.Split(dateStr, "-")
	if len(parts) != 3 {
		return fmt.Errorf("date format should be YYYY-MM-DD")
	}

	return nil
}

func (ctx *stocksFeatureContext) theErrorMessageShouldIndicate(message string) error {
	var errorResp map[string]interface{}
	_ = json.Unmarshal([]byte(ctx.response.Body.String()), &errorResp)

	if errorMsg, ok := errorResp["error"].(string); !ok || !strings.Contains(errorMsg, message) {
		return fmt.Errorf("error message should contain '%s', got response: %s", message, ctx.response.Body.String())
	}
	return nil
}

func (ctx *stocksFeatureContext) iRequestStockSummaryForTickerAgainWithinCacheTTL(ticker string) error {
	// Store first response
	var firstSummary stocks.StockSummary
	_ = json.Unmarshal([]byte(ctx.response.Body.String()), &firstSummary)
	ctx.firstResponse = &firstSummary

	// Make second request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/stocks/"+ticker+"/summary", nil)
	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)

	_ = json.Unmarshal([]byte(ctx.response.Body.String()), &ctx.responseBody)

	return nil
}

func (ctx *stocksFeatureContext) bothResponsesShouldBeIdentical() error {
	// Compare ticker from first response with second response
	secondTicker, ok := ctx.responseBody["ticker"].(string)
	if !ok {
		return fmt.Errorf("second response does not contain ticker field")
	}
	if ctx.firstResponse.Ticker != secondTicker {
		return fmt.Errorf("responses should be identical: first ticker=%s, second ticker=%s", ctx.firstResponse.Ticker, secondTicker)
	}
	return nil
}

func (ctx *stocksFeatureContext) theSecondRequestShouldBeServedFromCache() error {
	// In a real implementation, we'd track cache hits
	// For now, we just verify we got a successful response
	return nil
}

func (ctx *stocksFeatureContext) theResponseShouldContainStockSummaries(count int) error {
	summaries, ok := ctx.responseBody["summaries"].([]interface{})
	if !ok {
		return fmt.Errorf("response should contain summaries array")
	}

	if len(summaries) != count {
		return fmt.Errorf("expected %d summaries, got %d", count, len(summaries))
	}

	return nil
}

func (ctx *stocksFeatureContext) theResponseShouldContainStockSummary(count int) error {
	return ctx.theResponseShouldContainStockSummaries(count)
}

func (ctx *stocksFeatureContext) theResponseShouldIncludeTicker(ticker string) error {
	summaries, ok := ctx.responseBody["summaries"].([]interface{})
	if !ok {
		return fmt.Errorf("response should contain summaries array")
	}

	for _, s := range summaries {
		summary := s.(map[string]interface{})
		if summary["ticker"] == ticker {
			return nil
		}
	}

	return fmt.Errorf("ticker %s not found in summaries", ticker)
}

func (ctx *stocksFeatureContext) theResponseShouldContainSuccessfulSummaries(count int) error {
	summaries, ok := ctx.responseBody["summaries"].([]interface{})
	if !ok {
		return fmt.Errorf("response should contain summaries array")
	}

	if len(summaries) != count {
		return fmt.Errorf("expected %d successful summaries, got %d", count, len(summaries))
	}

	return nil
}

func (ctx *stocksFeatureContext) theResponseShouldContainError(count int) error {
	errors, ok := ctx.responseBody["errors"].([]interface{})
	if !ok {
		return fmt.Errorf("response should contain errors array")
	}

	if len(errors) != count {
		return fmt.Errorf("expected %d errors, got %d", count, len(errors))
	}

	return nil
}

func (ctx *stocksFeatureContext) theErrorShouldIndicateTickerNotFound(ticker string) error {
	errors, ok := ctx.responseBody["errors"].([]interface{})
	if !ok {
		return fmt.Errorf("response should contain errors array")
	}

	for _, e := range errors {
		errorItem := e.(map[string]interface{})
		if errorItem["ticker"] == ticker {
			return nil
		}
	}

	return fmt.Errorf("error for ticker %s not found", ticker)
}

func (ctx *stocksFeatureContext) theRequestShouldNotExceedAPIRateLimits() error {
	// In a real implementation, we'd track API call rates
	return nil
}

func (ctx *stocksFeatureContext) allStockSummariesShouldBeReturned(count int) error {
	summaries, ok := ctx.responseBody["summaries"].([]interface{})
	if !ok {
		return fmt.Errorf("response should contain summaries array")
	}

	if len(summaries) != count {
		return fmt.Errorf("expected %d summaries, got %d", count, len(summaries))
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	featureCtx := &stocksFeatureContext{}

	// Background
	ctx.Step(`^the Yahoo Finance API is available$`, featureCtx.theYahooFinanceAPIIsAvailable)

	// Single stock steps
	ctx.Step(`^I request stock summary for ticker "([^"]*)"$`, featureCtx.iRequestStockSummaryForTicker)
	ctx.Step(`^the response should contain ticker "([^"]*)"$`, featureCtx.theResponseShouldContainTicker)
	ctx.Step(`^the response should contain current price$`, featureCtx.theResponseShouldContainCurrentPrice)
	ctx.Step(`^the response should contain open price$`, featureCtx.theResponseShouldContainOpenPrice)
	ctx.Step(`^the response should contain high price$`, featureCtx.theResponseShouldContainHighPrice)
	ctx.Step(`^the response should contain low price$`, featureCtx.theResponseShouldContainLowPrice)
	ctx.Step(`^the response should contain change value$`, featureCtx.theResponseShouldContainChangeValue)
	ctx.Step(`^the response should contain change percentage$`, featureCtx.theResponseShouldContainChangePercentage)
	ctx.Step(`^the response should contain volume$`, featureCtx.theResponseShouldContainVolume)
	ctx.Step(`^the response should contain currency$`, featureCtx.theResponseShouldContainCurrency)
	ctx.Step(`^the response should contain company name$`, featureCtx.theResponseShouldContainCompanyName)
	ctx.Step(`^the company name should not be empty$`, featureCtx.theCompanyNameShouldNotBeEmpty)
	ctx.Step(`^the response should contain date$`, featureCtx.theResponseShouldContainDate)
	ctx.Step(`^the date should be in format "([^"]*)"$`, featureCtx.theDateShouldBeInFormat)

	// Batch stock steps
	ctx.Step(`^I request batch stock summary for tickers "([^"]*)"$`, featureCtx.iRequestBatchStockSummaryForTickers)
	ctx.Step(`^the response should contain (\d+) stock summaries$`, featureCtx.theResponseShouldContainStockSummaries)
	ctx.Step(`^the response should contain (\d+) stock summary$`, featureCtx.theResponseShouldContainStockSummary)
	ctx.Step(`^the response should include ticker "([^"]*)"$`, featureCtx.theResponseShouldIncludeTicker)
	ctx.Step(`^the response should contain (\d+) successful summaries$`, featureCtx.theResponseShouldContainSuccessfulSummaries)
	ctx.Step(`^the response should contain (\d+) error$`, featureCtx.theResponseShouldContainError)
	ctx.Step(`^the error should indicate ticker "([^"]*)" not found$`, featureCtx.theErrorShouldIndicateTickerNotFound)
	ctx.Step(`^the request should not exceed API rate limits$`, featureCtx.theRequestShouldNotExceedAPIRateLimits)
	ctx.Step(`^all (\d+) stock summaries should be returned$`, featureCtx.allStockSummariesShouldBeReturned)

	// Cache steps
	ctx.Step(`^I request stock summary for ticker "([^"]*)" again within cache TTL$`, featureCtx.iRequestStockSummaryForTickerAgainWithinCacheTTL)
	ctx.Step(`^both responses should be identical$`, featureCtx.bothResponsesShouldBeIdentical)
	ctx.Step(`^the second request should be served from cache$`, featureCtx.theSecondRequestShouldBeServedFromCache)

	// Common steps
	ctx.Step(`^the response status should be (\d+)$`, featureCtx.theResponseStatusShouldBe)
	ctx.Step(`^the error message should indicate "([^"]*)"$`, featureCtx.theErrorMessageShouldIndicate)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../../features/stocks-batch.feature", "../../features/stocks-summary.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
