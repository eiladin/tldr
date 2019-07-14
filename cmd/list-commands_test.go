package cmd

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/eiladin/tldr/cache"
	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
)

func TestListPages(t *testing.T) {
	color.NoColor = false

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
	}

	for _, test := range tests {
		var b bytes.Buffer
		listPlatformPages(&b, settings, test.platform)
		out := b.String()
		for _, expectation := range test.expectations {
			assert.Contains(t, out, expectation)
		}
	}

	os.RemoveAll(settings.Location)
}
