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
		{"Empty string", ""},
		{"Non-empty string", "hello"},
		{"String with spaces", "hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptr := StringPtr(tt.input)

			assert.NotNil(t, ptr)
			assert.Equal(t, tt.input, *ptr)
		})
	}
}

func TestIntPtr(t *testing.T) {
	tests := []struct {
		name  string
		input int
	}{
		{"Zero", 0},
		{"Positive", 42},
		{"Negative", -10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptr := IntPtr(tt.input)

			assert.NotNil(t, ptr)
			assert.Equal(t, tt.input, *ptr)
		})
	}
}
