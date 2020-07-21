package config

import (
	"runtime"
	"strings"
)

// CurrentPlatform will return the corrected name of the current platform
func CurrentPlatform() string {
	n := runtime.GOOS
	if strings.ToLower(n) == "darwin" {
		n = "osx"
	}
	return strings.ToLower(n)
}
