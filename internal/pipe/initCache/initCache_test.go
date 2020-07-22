package initCache

import (
	"os"
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

func TestInitCache(t *testing.T) {
	defer cleanTest()
	ctx := context.New()
	ctx.Cache.Location = testCacheDir
	ctx.Cache.TTL = time.Minute
	err := Pipe{}.Run(ctx)
	assert.NoError(t, err)
	assert.DirExists(t, "./test-cache/pages")
}
