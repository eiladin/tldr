package config

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOSName(t *testing.T) {
	expected := runtime.GOOS
	if expected == "darwin" {
		expected = "osx"
	}
	actual := OSName()
	assert.Equal(t, expected, actual, "Expected %s, Actual: %s", expected, actual)
}
