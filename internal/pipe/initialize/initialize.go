package initialize

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

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
	errReadingPagesDir     = errors.New("unable to read pages folder")
)

const zipPath = "tldr.zip"

// Pipe for initializing cache
type Pipe struct{}

func (Pipe) String() string {
	return "initializing cache"
}

// Run the pipe
func (Pipe) Run(ctx *context.Context) error {
	if err := createAndLoad(ctx); err != nil {
		return err
	}

	available, err := ioutil.ReadDir(path.Join(ctx.Cache.Location, ctx.PagesDirectory))
	if err != nil {
		return fmt.Errorf("cache: %s: %w", err, errReadingPagesDir)
	}

	for _, f := range available {
		if f.IsDir() {
			ctx.AvailablePlatforms = append(ctx.AvailablePlatforms, f.Name())
		}
	}
	return nil
}

func createAndLoad(ctx *context.Context) error {
	_, err := os.Stat(ctx.Cache.Location)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(ctx.Cache.Location, 0755); err != nil {
			return fmt.Errorf("cache: %s: %w", err, errCreatingCacheFolder)
		}
		if err := loadFromRemote(ctx); err != nil {
			return err
		}
	} else if err != nil {
		return fmt.Errorf("cache: %s: %w", err, errGettingCacheFolder)
	}
	return nil
}

func loadFromRemote(ctx *context.Context) error {
	dir, err := os.Create(path.Join(ctx.Cache.Location, zipPath))
	if err != nil {
		return fmt.Errorf("cache: %s: %w", err, errCreatingZip)
	}

	resp, err := http.Get(ctx.Cache.Remote)
	if err != nil {
		return fmt.Errorf("cache: %q: %w", ctx.Cache.Remote, errDownloadingFile)
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("cache: %q: %w", ctx.Cache.Remote, errDownloadingFile)
	}

	if _, err = io.Copy(dir, resp.Body); err != nil {
		return fmt.Errorf("cache: %s: %w", err, errSavingZipToCache)
	}

	if _, err = zip.Extract(path.Join(ctx.Cache.Location, zipPath), ctx.Cache.Location); err != nil {
		return err
	}
	dir.Close()

	if err = os.Remove(path.Join(ctx.Cache.Location, zipPath)); err != nil {
		return fmt.Errorf("cache: %s, %w", err, errRemovingZip)
	}
	return nil
}
