package initCache

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/eiladin/tldr/internal/pipe"
	"github.com/eiladin/tldr/internal/zip"
	"github.com/eiladin/tldr/pkg/context"
)

var (
	errCreatingCacheFolder = errors.New("unable to create cache folder")
	errGettingCacheFolder  = errors.New("unable to get cache folder")
	errCreatingZip         = errors.New("unable to create zip")
	errDownloadingFile     = errors.New("unable to download file")
	errSavingZipToCache    = errors.New("unable to save zip to cache")
	errRemovingZip         = errors.New("unable to remove zip")
)

const (
	zipPath = "/tldr.zip"
)

type Pipe struct{}

func (Pipe) String() string {
	return "initializing cache"
}

func (Pipe) Run(ctx *context.Context) error {
	_, err := os.Stat(ctx.Cache.Location)
	if os.IsNotExist(err) {
		if err = createAndLoad(ctx); err != nil {
			return err
		}
	} else if err != nil {
		return fmt.Errorf("cache: %s: %s", err, errGettingCacheFolder.Error())
	} else {
		return pipe.Skip("cache is up-to-date")
	}
	return nil
}

func createAndLoad(ctx *context.Context) error {
	if err := createCacheFolder(ctx); err != nil {
		return err
	}
	if err := loadFromRemote(ctx); err != nil {
		return err
	}
	return nil
}

func createCacheFolder(ctx *context.Context) error {
	err := os.MkdirAll(ctx.Cache.Location, 0755)
	if err != nil {
		return fmt.Errorf("cache: %s: %s", err, errCreatingCacheFolder.Error())
	}
	return nil
}

func loadFromRemote(ctx *context.Context) error {
	dir, err := os.Create(ctx.Cache.Location + zipPath)
	if err != nil {
		return fmt.Errorf("cache: %s: %s", err, errCreatingZip.Error())
	}

	resp, err := http.Get(ctx.Cache.Remote)
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("cache: %q: %s", ctx.Cache.Remote, errDownloadingFile.Error())
	}

	if _, err = io.Copy(dir, resp.Body); err != nil {
		return fmt.Errorf("cache: %s: %s", err.Error(), errSavingZipToCache.Error())
	}

	if _, err = zip.Extract(ctx.Cache.Location+zipPath, ctx.Cache.Location); err != nil {
		return err
	}
	dir.Close()

	if err = os.Remove(ctx.Cache.Location + zipPath); err != nil {
		return fmt.Errorf("cache: %s, %s", err, errRemovingZip.Error())
	}
	return nil
}
