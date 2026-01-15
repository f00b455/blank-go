package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	version := Version()

	assert.NotEmpty(t, version)
	assert.Equal(t, "0.1.0", version)
}
