package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/eiladin/tldr/cache"
	"github.com/stretchr/testify/assert"
)

func TestListAvailablePlatforms(t *testing.T) {
	settings := cache.Cache{
		Remote:   "http://tldr-pages.github.com/assets/tldr.zip",
		TTL:      time.Minute,
		Location: "./tldr-platforms-cmd-test",
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	listAvailablePlatforms(settings)

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()
	w.Close()
	out := <-outC
	os.Stdout = old
	assert.Contains(t, out, "common")
	assert.Contains(t, out, "linux")
	assert.Contains(t, out, "osx")
	assert.Contains(t, out, "sunos")
	assert.Contains(t, out, "windows")

	os.RemoveAll(settings.Location)
}
