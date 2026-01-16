package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringPtr(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "returns pointer to empty string",
			input: "",
		},
		{
			name:  "returns pointer to non-empty string",
			input: "test",
		},
		{
			name:  "returns pointer to string with spaces",
			input: "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringPtr(tt.input)
			assert.NotNil(t, result)
			assert.Equal(t, tt.input, *result)
		})
	}
}

func TestIntPtr(t *testing.T) {
	tests := []struct {
		name  string
		input int
	}{
		{
			name:  "returns pointer to zero",
			input: 0,
		},
		{
			name:  "returns pointer to positive number",
			input: 42,
		},
		{
			name:  "returns pointer to negative number",
			input: -100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IntPtr(tt.input)
			assert.NotNil(t, result)
			assert.Equal(t, tt.input, *result)
		})
	}
}
