package cmd

import (
	"bytes"
	"testing"
	"time"

	"github.com/eiladin/tldr/cache"
	"github.com/stretchr/testify/assert"
)

func TestPlatformListPages(t *testing.T) {

	settings := cache.Cache{
		Remote:   "http://tldr-pages.github.com/assets/tldr.zip",
		TTL:      time.Minute,
		Location: "./tldr-list-cmd-test",
	}

	tests := []struct {
		platform     string
		expectations []string
	}{
		{"common", []string{"curl", "ls"}},
		{"linux", []string{"dmesg", "alpine"}},
		{"osx", []string{"dmesg", "brew"}},
		{"sunos", []string{"dmesg", "stty"}},
		{"windows", []string{"rmdir", "mkdir"}},
		{"all", []string{"where", "choco", "brew"}},
	}

	for _, test := range tests {
		var b bytes.Buffer
		listPlatformPages(&b, settings, test.platform)
		out := b.String()
		for _, expectation := range test.expectations {
			assert.Contains(t, out, expectation)
		}
	}

	settings.Purge() // nolint: errcheck
}
