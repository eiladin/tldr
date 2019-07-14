package cmd

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/eiladin/tldr/cache"
	"github.com/stretchr/testify/assert"
)

func TestGetAvailablePlatforms(t *testing.T) {
	settings := cache.Cache{
		Remote:   "http://tldr-pages.github.com/assets/tldr.zip",
		TTL:      time.Minute,
		Location: "./tldr-platforms-cmd-test",
	}

	var b bytes.Buffer
	listAvailablePlatforms(&b, settings)
	out := b.String()
	assert.Contains(t, out, "common")
	assert.Contains(t, out, "linux")
	assert.Contains(t, out, "osx")
	assert.Contains(t, out, "sunos")
	assert.Contains(t, out, "windows")

	os.RemoveAll(settings.Location)
}
