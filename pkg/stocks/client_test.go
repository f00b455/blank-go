package stocks_test

import (
	"testing"

	"github.com/f00b455/blank-go/pkg/stocks"
	"github.com/stretchr/testify/assert"
)

func TestClient_GetQuote_EmptyTicker(t *testing.T) {
	client := stocks.NewClient()

	quote, err := client.GetQuote("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ticker is required")
	assert.Nil(t, quote)
}

func TestClient_GetQuotes_EmptyTickers(t *testing.T) {
	client := stocks.NewClient()

	quotes, err := client.GetQuotes([]string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one ticker is required")
	assert.Nil(t, quotes)
}

func TestNewClient(t *testing.T) {
	client := stocks.NewClient()
	assert.NotNil(t, client)
}
