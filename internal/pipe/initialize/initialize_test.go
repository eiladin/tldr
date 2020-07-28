package initialize

import (
	"errors"
	"os"
	"path"
	"testing"
	"time"

	"github.com/eiladin/tldr/pkg/context"
	"github.com/stretchr/testify/assert"
)

var testCacheDir = "./test-cache"

func cleanTest() {
	os.RemoveAll(testCacheDir)
}

type test struct {
	platform     string
	expectations []string
}

func TestString(t *testing.T) {
	p := Pipe{}
	assert.NotEmpty(t, p.String())
}

func TestInitCache(t *testing.T) {
	defer cleanTest()
	ctx := context.New()
	ctx.Cache.Location = testCacheDir
	ctx.Cache.TTL = time.Minute
	err := Pipe{}.Run(ctx)
	assert.NoError(t, err)
	assert.DirExists(t, "./test-cache/pages")
	assert.Contains(t, ctx.AvailablePlatforms, "linux")
	assert.Contains(t, ctx.AvailablePlatforms, "osx")
	assert.Contains(t, ctx.AvailablePlatforms, "sunos")
	assert.Contains(t, ctx.AvailablePlatforms, "windows")
}

func TestReadError(t *testing.T) {
	ctx := context.New()
	ctx.Cache.Location = "./read-error-test"
	os.MkdirAll(path.Join(ctx.Cache.Location, "pages"), 0)
	defer os.RemoveAll(ctx.Cache.Location)
	err := Pipe{}.Run(ctx)
	assert.True(t, errors.Is(err, errReadingPagesDir))
}

func TestDownloadError(t *testing.T) {
	ctx := context.New()
	ctx.Cache.Location = "./download-error-test"
	defer os.RemoveAll(ctx.Cache.Location)
	ctx.Cache.Remote = "http://404.not-found.url/tldr.zip"
	err := Pipe{}.Run(ctx)
	assert.True(t, errors.Is(err, errDownloadingFile))
}
