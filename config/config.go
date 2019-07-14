package config

import (
	"runtime"
	"time"
)

// OSName will return the corrected name of the current platform
func OSName() (n string) {
	n = runtime.GOOS
	if n == "darwin" {
		n = "osx"
	}
	return
}

const (
	// DefaultTTL is 1 week
	DefaultTTL = time.Hour * 24 * 7
	// DefaultRemoteURL pulls from tldr-pages.github.com
	DefaultRemoteURL = "http://tldr-pages.github.com/assets/tldr.zip"
)
