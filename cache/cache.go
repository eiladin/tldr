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
	"strings"
	"time"

	"github.com/eiladin/tldr/zip"
	"github.com/mitchellh/go-homedir"
)

const (
	zipPath        = "/tldr.zip"
	pagesDirectory = "pages"
	pageSuffix     = ".md"
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
func Create(remote string, ttl time.Duration, folder string) (*Cache, error) {
	dir, err := getCacheDir(folder)
	if err != nil {
		return nil, fmt.Errorf("ERROR: getting cache directory: %s", err)
	}

	cache := &Cache{Location: dir, Remote: remote, TTL: ttl}

	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err = cache.createAndLoad(); err != nil {
			return nil, fmt.Errorf("ERROR: creating cache: %s", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("ERROR: creating cache: %s", err)
	} else {
		cacheExpired := info.ModTime().Before(time.Now().Add(-ttl))
		children, _ := ioutil.ReadDir(cache.Location)
		if cacheExpired || len(children) == 0 {
			if err = cache.Refresh(); err != nil {
				return nil, fmt.Errorf("ERROR: refreshing cache: %s", err)
			}
		}
	}

	return cache, nil
}

// Refresh the cache with the latest info
func (cache *Cache) Refresh() error {
	if err := os.RemoveAll(cache.Location); err != nil {
		return fmt.Errorf("ERROR: removing cache directory: %s", err)
	}
	if err := cache.createAndLoad(); err != nil {
		return fmt.Errorf("ERROR: creating cache directory: %s", err)
	}
	return nil
}

// FetchPage returns a specific page from cache
func (cache *Cache) FetchPage(platform, page string) (io.ReadCloser, string, error) {
	pform := platform
	platformPath := path.Join(cache.Location, pagesDirectory, platform, page+pageSuffix)
	commonPath := path.Join(cache.Location, pagesDirectory, "common", page+pageSuffix)

	paths := []string{platformPath, commonPath}
	for _, p := range paths {
		if p == commonPath {
			pform = "common"
		}
		if _, err := os.Stat(p); os.IsNotExist(err) {
			continue
		} else {
			closer, err := os.Open(p)
			return closer, pform, err
		}
	}

	return nil, "", errors.New("This page (" + page + ") does not exist yet!\n" +
		"Submit new pages here: https://github.com/tldr-pages/tldr")
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

// ListPages returns all pages for a given platform
func (cache *Cache) ListPages(platform string) ([]string, error) {
	dir := path.Join(cache.Location, pagesDirectory, platform)
	pages := []os.FileInfo{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".md") {
			pages = append(pages, f)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("ERROR: getting pages list")
	}

	names := make([]string, len(pages))
	for i, page := range pages {
		name := page.Name()
		names[i] = name[:len(name)-3]
	}
	return names, nil
}

// AvailablePlatforms returns all platforms available in the cache
func (cache *Cache) AvailablePlatforms() ([]string, error) {
	var platforms []string
	available, err := ioutil.ReadDir(path.Join(cache.Location, pagesDirectory))
	if err != nil {
		return nil, fmt.Errorf("ERROR: reading pages folder: %s", err)
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
		return fmt.Errorf("ERROR: getting cache folder: %s", err)
	}
	return os.RemoveAll(dir)
}

func (cache *Cache) createAndLoad() error {
	if err := cache.createCacheFolder(); err != nil {
		return fmt.Errorf("ERROR: creating cache directory: %s", err)
	}
	if err := cache.loadFromRemote(); err != nil {
		return fmt.Errorf("ERROR: loading data from remote: %s", err)
	}
	return nil
}

func (cache *Cache) createCacheFolder() error {
	return os.MkdirAll(cache.Location, 0755)
}

func (cache *Cache) loadFromRemote() error {
	dir, err := os.Create(cache.Location + zipPath)
	if err != nil {
		return fmt.Errorf("ERROR: creating cache folder: %s", err)
	}
	defer dir.Close()

	resp, err := http.Get(cache.Remote)
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ERROR: downloading zip: %s", err)
	}

	if _, err = io.Copy(dir, resp.Body); err != nil {
		return fmt.Errorf("ERROR: saving zip to cache: %s", err)
	}

	if _, err = zip.Extract(cache.Location+zipPath, cache.Location); err != nil {
		return fmt.Errorf("ERROR: extracting zip: %s", err)
	}

	if err = os.Remove(cache.Location + zipPath); err != nil {
		return fmt.Errorf("ERROR: removing zip file: %s", err)
	}
	return nil
}

func getCacheDir(folder string) (string, error) {
	if folder == "" {
		home, err := homedir.Dir()
		if err != nil {
			return "", fmt.Errorf("ERROR: getting current user's home directory: %s", err)
		}
		return path.Join(home, ".tldr"), nil
	}
	return folder, nil
}
