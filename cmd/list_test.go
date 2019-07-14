package cmd

import (
	"bytes"
	"io"
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
		Location: "./tldr-list-test",
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
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		listPlatformPages(settings, test.platform)

		outC := make(chan string)
		go func() {
			var buf bytes.Buffer
			io.Copy(&buf, r)
			outC <- buf.String()
		}()
		w.Close()
		out := <-outC
		os.Stdout = old
		for _, expectation := range test.expectations {
			assert.Contains(t, out, expectation)
		}
	}

	os.RemoveAll("./tldr-test")
}
