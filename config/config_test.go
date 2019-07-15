package config

import (
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurrentPlatform(t *testing.T) {
	expected := strings.ToLower(runtime.GOOS)
	if expected == "darwin" {
		expected = "osx"
	}
	actual := CurrentPlatform()
	assert.Equal(t, expected, actual, "Expected %s, Actual: %s", expected, actual)
}
