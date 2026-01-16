package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "returns correct version",
			expected: "0.1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Version()
			assert.Equal(t, tt.expected, result)
			assert.NotEmpty(t, result)
		})
	}
}
