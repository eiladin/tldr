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

type Cache struct {
	location string
	remote   string
	ttl      time.Duration
}

func Create(remote string, ttl time.Duration) (*Cache, error) {
	dir, err := GetCacheDir()
	if err != nil {
		return nil, fmt.Errorf("ERROR: getting cache directory: %s", err)
	}

	cache := &Cache{location: dir, remote: remote, ttl: ttl}

	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = cache.CreateAndLoad()
		if err != nil {
			return nil, fmt.Errorf("ERROR: creating cache: %s", err)
		}
	} else if err != nil || info.ModTime().Before(time.Now().Add(-ttl)) {
		err = cache.Refresh()
		if err != nil {
			return nil, fmt.Errorf("ERROR: refreshing cache: %s", err)
		}
	}

	return cache, nil
}

func (cache *Cache) CreateAndLoad() error {
	err := cache.CreateCacheFolder()
	if err != nil {
		return fmt.Errorf("ERROR: creating cache directory: %s", err)
	}
	err = cache.LoadFromRemote()
	if err != nil {
		return fmt.Errorf("ERROR: loading data from remote: %s", err)
	}
	return nil
}

func (cache *Cache) Refresh() error {
	err := os.RemoveAll(cache.location)
	if err != nil {
		return fmt.Errorf("ERROR: removing cache directory: %s", err)
	}
	err = cache.CreateAndLoad()
	if err != nil {
		return fmt.Errorf("ERROR: creating cache directory: %s", err)
	}
	return nil
}

func (cache *Cache) Fetch(platform, page string) (io.ReadCloser, error) {
	filePath := path.Join(cache.location, pagesDirectory, platform, page+pageSuffix)
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		filePath = path.Join(cache.location, pagesDirectory, "common", page+pageSuffix)
		_, err = os.Stat(filePath)
		if os.IsNotExist(err) {
			return nil, errors.New("This page (" + page + ") doesn't exist yet!\n" +
				"Submit new pages here: https://github.com/tldr-pages/tldr")
		}
	}

	return os.Open(filePath)
}

func (cache *Cache) FetchRandom(platform string) (io.ReadCloser, error) {
	cmn := path.Join(cache.location, pagesDirectory, "common")
	plt := path.Join(cache.location, pagesDirectory, platform)
	paths := []string{cmn, plt}
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

func (cache *Cache) CreateCacheFolder() error {
	return os.MkdirAll(cache.location, 0755)
}

func (cache *Cache) LoadFromRemote() error {
	dir, err := os.Create(cache.location + zipPath)
	if err != nil {
		return fmt.Errorf("ERROR: creating cache: %s", err)
	}
	defer dir.Close()

	resp, err := http.Get(cache.remote)
	if err != nil {
		return fmt.Errorf("ERROR: downloading zip: %s", err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(dir, resp.Body)
	if err != nil {
		return fmt.Errorf("ERROR: saving zip to cache: %s", err)
	}

	_, err = zip.Extract(cache.location+zipPath, cache.location)
	if err != nil {
		return fmt.Errorf("ERROR: extracting zip: %s", err)
	}

	err = os.Remove(cache.location + zipPath)
	if err != nil {
		return fmt.Errorf("ERROR: removing zip file: %s", err)
	}
	return nil
}

func GetCacheDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("ERROR: getting current user: %s", err)
	}
	if usr.HomeDir == "" {
		return "", fmt.Errorf("ERROR: loading current user's home directory")
	}
	return path.Join(usr.HomeDir, ".tldr"), nil
}
