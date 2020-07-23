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

type Pipe struct{}

func (Pipe) String() string {
	return "initializing cache"
}

func (Pipe) Run(ctx *context.Context) error {
	if err := createAndLoad(ctx); err != nil {
		return err
	}

	available, err := ioutil.ReadDir(path.Join(ctx.Cache.Location, ctx.PagesDirectory))
	if err != nil {
		return fmt.Errorf("cache: %s: %s", err, errReadingPagesDir.Error())
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
			return fmt.Errorf("cache: %s: %s", err, errCreatingCacheFolder.Error())
		}
		if err := loadFromRemote(ctx); err != nil {
			return err
		}
	} else if err != nil {
		return fmt.Errorf("cache: %s: %s", err, errGettingCacheFolder.Error())
	}
	return nil
}

func loadFromRemote(ctx *context.Context) error {
	dir, err := os.Create(path.Join(ctx.Cache.Location + zipPath))
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
