package cache

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path"
	"time"
)

const (
	zipPath = "/tldr.zip"
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
		err = cache.CreateCacheFolder()
		if err != nil {
			return nil, fmt.Errorf("ERROR: creating cache directory: %s", err)
		}
		err = cache.LoadFromRemote()
		if err != nil {
			return nil, fmt.Errorf("ERROR: loading data from remote: %s", err)
		}
	} else if err != nil || info.ModTime().Before(time.Now().Add(-ttl)) {
		// err = cache.Reload()
		// if err != nil {
		// 	return nil, fmt.Errorf("ERROR: reloading cache: %s", err)
		// }
	}

	return cache, nil
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
		return fmt.Errorf("ERROR: downloading tldr.zip: %s", err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(dir, resp.Body)
	if err != nil {
		return fmt.Errorf("ERROR: saving tldr.zip to cache: %s", err)
	}

	//TODO: Unzip tldr.zip
	zip.Extract(cache.location+zipPath, cache.location)
	return nil
}

func GetCacheDir() (string, err) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("ERROR: getting current user: %s", err)
	}
	if usr.HomeDir == "" {
		return "", fmt.Errorf("ERROR: loading current user's home directory")
	}
	return path.Join(usr.HomeDir, ".tldr"), nil
}
