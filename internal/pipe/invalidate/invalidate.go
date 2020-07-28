package invalidate

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/eiladin/tldr/internal/pipe"
	"github.com/eiladin/tldr/pkg/context"
)

var errRemovingCacheFolder = errors.New("unable to remove cache folder")

// Pipe for invalidating the cache
type Pipe struct{}

func (Pipe) String() string {
	return "clearing cache"
}

// Run the pipe
func (Pipe) Run(ctx *context.Context) error {
	cacheExpired := false
	info, err := os.Stat(ctx.Cache.Location)
	if !os.IsNotExist(err) {
		before := info.ModTime().Before(time.Now().Add(-ctx.Cache.TTL))
		children, _ := ioutil.ReadDir(ctx.Cache.Location)
		if before || len(children) == 0 {
			cacheExpired = true
		}
	}

	if ctx.PurgeCache || cacheExpired {
		err = os.RemoveAll(ctx.Cache.Location)
		if err != nil {
			return fmt.Errorf("cache: %s: %s", err, errRemovingCacheFolder.Error())
		}
	} else {
		return pipe.Skip("cache is up-to-date")
	}
	return nil
}
