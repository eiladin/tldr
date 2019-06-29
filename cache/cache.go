package cache

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/user"
	"path"
	"strings"
	"time"

	"github.com/eiladin/tldr/zip"
)

const (
	zipPath        = "/tldr.zip"
	pagesDirectory = "pages"
	pageSuffix     = ".md"
)

// Cache stuct
type Cache struct {
	location string
	remote   string
	ttl      time.Duration
}

// Create a new Cache and populate it
func Create(remote string, ttl time.Duration) (*Cache, error) {
	dir, err := getCacheDir()
	if err != nil {
		return nil, fmt.Errorf("ERROR: getting cache directory: %s", err)
	}

	cache := &Cache{location: dir, remote: remote, ttl: ttl}

	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err = cache.createAndLoad(); err != nil {
			return nil, fmt.Errorf("ERROR: creating cache: %s", err)
		}
	} else if err != nil || info.ModTime().Before(time.Now().Add(-ttl)) {
		if err = cache.Refresh(); err != nil {
			return nil, fmt.Errorf("ERROR: refreshing cache: %s", err)
		}
	}

	return cache, nil
}

// Refresh the cache with the latest info
func (cache *Cache) Refresh() error {
	if err := os.RemoveAll(cache.location); err != nil {
		return fmt.Errorf("ERROR: removing cache directory: %s", err)
	}
	if err := cache.createAndLoad(); err != nil {
		return fmt.Errorf("ERROR: creating cache directory: %s", err)
	}
	return nil
}

// Fetch a specific page from cache
func (cache *Cache) Fetch(platform, page string) (io.ReadCloser, error) {
	platformPath := path.Join(cache.location, pagesDirectory, platform, page+pageSuffix)
	commonPath := path.Join(cache.location, pagesDirectory, "common", page+pageSuffix)

	paths := []string{platformPath, commonPath}
	for _, p := range paths {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			continue
		} else {
			return os.Open(p)
		}
	}

	return nil, errors.New("This page (" + page + ") does not exist yet!\n" +
		"Submit new pages here: https://github.com/tldr-pages/tldr")
}

// FetchRandom returns a random page from cache
func (cache *Cache) FetchRandom(platform string) (io.ReadCloser, error) {
	commonPath := path.Join(cache.location, pagesDirectory, "common")
	platformPath := path.Join(cache.location, pagesDirectory, platform)
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
	return os.Open(page)
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
	return os.MkdirAll(cache.location, 0755)
}

func (cache *Cache) loadFromRemote() error {
	dir, err := os.Create(cache.location + zipPath)
	if err != nil {
		return fmt.Errorf("ERROR: creating cache folder: %s", err)
	}
	defer dir.Close()

	resp, err := http.Get(cache.remote)
	if err != nil {
		return fmt.Errorf("ERROR: downloading zip: %s", err)
	}
	defer resp.Body.Close()

	if _, err = io.Copy(dir, resp.Body); err != nil {
		return fmt.Errorf("ERROR: saving zip to cache: %s", err)
	}

	if _, err = zip.Extract(cache.location+zipPath, cache.location); err != nil {
		return fmt.Errorf("ERROR: extracting zip: %s", err)
	}

	if err = os.Remove(cache.location + zipPath); err != nil {
		return fmt.Errorf("ERROR: removing zip file: %s", err)
	}
	return nil
}

func getCacheDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("ERROR: getting current user: %s", err)
	}
	if usr.HomeDir == "" {
		return "", fmt.Errorf("ERROR: loading current user's home directory")
	}
	return path.Join(usr.HomeDir, ".tldr"), nil
}
