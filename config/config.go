package config

import (
	"runtime"
)

// OSName will return the corrected name of the current platform
func OSName() (n string) {
	n = runtime.GOOS
	if n == "darwin" {
		n = "osx"
	}
	return
}
