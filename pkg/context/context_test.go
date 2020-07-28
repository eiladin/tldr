package context

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetHomeDir(t *testing.T) {
	os.Setenv("HOME", "./test-home")
	dir, err := getCacheDir()
	assert.NoError(t, err)
	assert.Equal(t, "test-home/.tldr", dir)
}

func TestNew(t *testing.T) {
	ctx := New()
	assert.NotEmpty(t, ctx.Cache.Location, "Cache Location should be set")
	assert.Equal(t, time.Hour*24*7, ctx.Cache.TTL, "TTL should be 7 days")
	assert.Equal(t, os.Stdout, ctx.Writer, "Writer should be `os.StdOut`")
	assert.Equal(t, "pages", ctx.PagesDirectory, "PagesDirectory should be `pages`")
	assert.Equal(t, ".md", ctx.PageSuffix, "PageSuffix should be `.md`")
	assert.Equal(t, OperationRenderPage, ctx.Operation, "Operation should be RenderPage")
}

func TestRenderPlatform(t *testing.T) {
	cases := []struct {
		platform string
		found    string
		expected string
	}{
		{"linux", "linux", "linux"},
		{"linux", "common", "linux (common)"},
	}
	for _, c := range cases {
		ctx := New()
		ctx.FoundPlatform = c.found
		ctx.Platform = c.platform
		assert.Equal(t, c.expected, ctx.RenderPlatform())
	}
}
