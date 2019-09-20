package cache

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/eiladin/tldr/zip"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/xerrors"
)

const (
	zipPath        = "/tldr.zip"
	pagesDirectory = "pages"
	pageSuffix     = ".md"
)

var (
	ErrCreatingCacheFolder = errors.New("unable to create cache folder")
	ErrCreatingZip         = errors.New("unable to create zip")
	ErrDownloadingFile     = errors.New("unable to download file")
	ErrGettingCacheFolder  = errors.New("unable to get cache folder")
	ErrGettingHomeDir      = errors.New("unable to get user's home")
	ErrListingPages        = errors.New("unable to list pages")
	ErrReadingPagesDir     = errors.New("unable to read pages folder")
	ErrRemovingCacheFolder = errors.New("unable to remove cache folder")
	ErrRemovingZip         = errors.New("unable to remove zip")
	ErrSavingZipToCache    = errors.New("unable to save zip to cache")
)

// Cache stuct
type Cache struct {
	Location string
	Remote   string
	TTL      time.Duration
}

//DefaultSettings for the cache
var DefaultSettings = Cache{
	TTL:    time.Hour * 24 * 7,
	Remote: "http://tldr-pages.github.com/assets/tldr.zip",
}

// Create a new Cache and populate it
func Create(w io.Writer, remote string, ttl time.Duration, folder string) (*Cache, error) {
	dir, err := getCacheDir(folder)
	if err != nil {
		return nil, err
	}

	cache := &Cache{Location: dir, Remote: remote, TTL: ttl}

	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err = cache.createAndLoad(); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, xerrors.Errorf("cache: %s: %w", err, ErrGettingCacheFolder)
	} else {
		cacheExpired := info.ModTime().Before(time.Now().Add(-ttl))
		children, _ := ioutil.ReadDir(cache.Location)
		if cacheExpired || len(children) == 0 {
			fmt.Fprint(w, "Cache Expired - Refreshing ... ")
			if err = cache.Refresh(); err != nil {
				return nil, err
			}
			fmt.Fprintln(w, "Done")
		}
	}

	return cache, nil
}

// Refresh the cache with the latest info
func (cache *Cache) Refresh() error {
	if err := os.RemoveAll(cache.Location); err != nil {
		return xerrors.Errorf("cache: %s: %w", err, ErrRemovingCacheFolder)
	}
	return cache.createAndLoad()
}

type platformPath struct {
	Platform string
	Path     string
}

func (cache *Cache) makePath(platform string, page string) string {
	return path.Join(cache.Location, pagesDirectory, platform, page+pageSuffix)
}

func (cache *Cache) getPage(platforms []string, page string) (string, string, error) {
	ps := make([]platformPath, 0)
	for _, p := range platforms {
		ps = append(ps, platformPath{Platform: p, Path: cache.makePath(p, page)})
	}

	for _, pp := range ps {
		if _, err := os.Stat(pp.Path); os.IsNotExist(err) {
			continue
		} else {
			return pp.Platform, pp.Path, nil
		}
	}

	return "", "", xerrors.New("This page (" + page + ") does not exist yet!\n" +
		"Submit new pages here: https://github.com/tldr-pages/tldr/issues/new?title=page%20request:%20" + page)
}

// FetchPage returns a specific page from cache
func (cache *Cache) FetchPage(platform, page string) (io.ReadCloser, string, error) {
	platforms := []string{platform, "common"}
	p, f, err := cache.getPage(platforms, page)
	if err != nil {
		platforms, err = cache.AvailablePlatforms()
		if err != nil {
			return nil, "", err
		}
		p, f, err = cache.getPage(platforms, page)
		if err != nil {
			return nil, "", err
		}
	}

	closer, err := os.Open(f)
	return closer, p, err
}

// FetchRandomPage returns a random page from cache
func (cache *Cache) FetchRandomPage(platform string) (io.ReadCloser, string, error) {
	commonPath := path.Join(cache.Location, pagesDirectory, "common")
	platformPath := path.Join(cache.Location, pagesDirectory, platform)
	paths := []string{commonPath, platformPath}
	srcs := make([]string, 0)
	for _, p := range paths {
		files, err := ioutil.ReadDir(p)
		if err != nil {
			break
		}
		for _, f := range files {
			if strings.HasSuffix(f.Name(), pageSuffix) {
				srcs = append(srcs, path.Join(p, f.Name()))
			}
		}
	}
	rand.Seed(time.Now().UTC().UnixNano())
	page := srcs[rand.Intn(len(srcs))]
	if strings.Contains(page, "common") {
		platform = "common"
	}
	reader, err := os.Open(page)
	return reader, platform, err
}

func (cache *Cache) getPages(platform string) ([]string, error) {
	dir := path.Join(cache.Location, pagesDirectory, platform)
	pages := []os.FileInfo{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".md") {
			pages = append(pages, f)
		}
		return nil
	})
	if err != nil {
		return nil, xerrors.Errorf("cache: %s: %w", err, ErrListingPages)
	}

	names := make([]string, len(pages))
	for i, page := range pages {
		name := page.Name()
		names[i] = name[:len(name)-3]
	}
	return names, nil

}

func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	r := []string{}
	for _, s := range stringSlice {
		if _, v := keys[s]; !v {
			keys[s] = true
			r = append(r, s)
		}
	}
	return r
}

// ListPages returns all pages for a given platform
func (cache *Cache) ListPages(platform string) ([]string, error) {
	if strings.ToLower(platform) == "all" {
		pgs := make([]string, 0)
		pfs, err := cache.AvailablePlatforms()
		if err != nil {
			return nil, err
		}
		for _, p := range pfs {
			ppg, err := cache.getPages(p)
			if err != nil {
				return nil, err
			}
			for _, pg := range ppg {
				pgs = append(pgs, pg)
			}
		}
		pgs = unique(pgs)
		sort.Strings(pgs)
		return pgs, nil
	}
	return cache.getPages(platform)
}

// AvailablePlatforms returns all platforms available in the cache
func (cache *Cache) AvailablePlatforms() ([]string, error) {
	var platforms []string
	available, err := ioutil.ReadDir(path.Join(cache.Location, pagesDirectory))
	if err != nil {
		return nil, xerrors.Errorf("cache: %s: %w", err, ErrReadingPagesDir)
	}

	for _, f := range available {
		platform := f.Name()
		if f.IsDir() {
			platforms = append(platforms, platform)
		}
	}
	return platforms, nil
}

// IsPlatformValid ensures the provided platform is found in the cache
func (cache *Cache) IsPlatformValid(platform string) (bool, []string) {
	platforms, _ := cache.AvailablePlatforms()
	for _, p := range platforms {
		if p == platform {
			return true, platforms
		}
	}
	return false, platforms
}

//Purge deletes the cache
func (cache *Cache) Purge() error {
	dir, err := getCacheDir(cache.Location)
	if err != nil {
		return err
	}
	err = os.RemoveAll(dir)
	if err != nil {
		return xerrors.Errorf("cache: %s: %w", err, ErrRemovingCacheFolder)
	}
	return nil
}

func (cache *Cache) createAndLoad() error {
	if err := cache.createCacheFolder(); err != nil {
		return err
	}
	if err := cache.loadFromRemote(); err != nil {
		return err
	}
	return nil
}

func (cache *Cache) createCacheFolder() error {
	err := os.MkdirAll(cache.Location, 0755)
	if err != nil {
		return xerrors.Errorf("cache: %s: %w", err, ErrCreatingCacheFolder)
	}
	return nil
}

func (cache *Cache) loadFromRemote() error {
	dir, err := os.Create(cache.Location + zipPath)
	if err != nil {
		return xerrors.Errorf("cache: %s: %w", err, ErrCreatingZip)
	}

	resp, err := http.Get(cache.Remote)
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return xerrors.Errorf("cache: %q: %w", cache.Remote, ErrDownloadingFile)
	}

	if _, err = io.Copy(dir, resp.Body); err != nil {
		return xerrors.Errorf("cache: %s: %w", err, ErrSavingZipToCache)
	}

	if _, err = zip.Extract(cache.Location+zipPath, cache.Location); err != nil {
		return err
	}
	dir.Close()

	if err = os.Remove(cache.Location + zipPath); err != nil {
		return xerrors.Errorf("cache: %s, %w", err, ErrRemovingZip)
	}
	return nil
}

func getCacheDir(folder string) (string, error) {
	if folder == "" {
		home, err := homedir.Dir()
		if err != nil {
			return "", xerrors.Errorf("cache: %s: %w", err, ErrGettingHomeDir)
		}
		return path.Join(home, ".tldr"), nil
	}
	return folder, nil
}
