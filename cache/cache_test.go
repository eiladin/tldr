package cache

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	remoteURL = "http://tldr-pages.github.com/assets/tldr.zip"
	ttl       = time.Minute
	location  = "./cache-test"
)

var cache *Cache

func CreateCache() *Cache {
	cache, _ = Create(remoteURL, ttl, location)
	return cache
}

func DestroyCache() {
	cache.Purge()
}

func TestMain(m *testing.M) {
	CreateCache()
	code := m.Run()
	DestroyCache()
	os.Exit(code)
}

func TestFetchPage(t *testing.T) {
	tests := []struct {
		platform    string
		page        string
		outPlatform string
	}{
		{"linux", "cat", "common"},
		{"windows", "choco", "windows"},
		{"osx", "airport", "osx"},
		{"linux", "pacman", "linux"},
		{"sunos", "dmesg", "sunos"},
	}
	for _, test := range tests {
		readCloser, pform, _ := cache.FetchPage(test.platform, test.page)
		assert.Equal(t, test.outPlatform, pform, fmt.Sprintf("Platform should match: %s", test.outPlatform))
		readCloser.Close()
	}

	_, _, err := cache.FetchPage("linux", "qwaszx")
	assert.Error(t, err, "Should result in a not exist error")
}

func TestLoadFromRemote(t *testing.T) {
	tests := []struct {
		remote    string
		location  string
		fileMode  os.FileMode
		shouldErr bool
	}{
		{remoteURL, "./tldr-load", 0755, false},
		{"https://github.com/eiladin/not-found.zip", "tldr-not-found", 0755, true},
		{remoteURL, "./tldr-perm-error", 0100, true},
	}

	for _, test := range tests {
		os.Mkdir(test.location, test.fileMode)
		cache := Cache{TTL: ttl, Location: test.location, Remote: test.remote}
		err := cache.loadFromRemote()
		if test.shouldErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		cache.Purge()
	}
}

func TestCreateAndLoad(t *testing.T) {
	location := "./tldr-fail"
	os.Mkdir(location, 0100)
	cache := Cache{TTL: time.Minute, Location: location, Remote: remoteURL}
	err := cache.createAndLoad()
	assert.Error(t, err)
	cache.Purge()
}

func TestCreateCacheFolder(t *testing.T) {
	location := "./cache-create"
	cache := Cache{TTL: ttl, Location: location, Remote: remoteURL}
	cache.createCacheFolder()
	dir, err := os.Stat(location)
	assert.NoError(t, err, "There should be no error getting the directory")
	assert.Equal(t, true, dir.IsDir())
	cache.Purge()
}

func TestFetchRandomPage(t *testing.T) {
	tests := []struct {
		platform     string
		outPlatforms []string
	}{
		{"linux", []string{"linux", "common"}},
		{"sunos", []string{"sunos", "common"}},
		{"windows", []string{"windows", "common"}},
		{"osx", []string{"osx", "common"}},
		{"common", []string{"common"}},
	}

	for _, test := range tests {
		readCloser, pform, _ := cache.FetchRandomPage(test.platform)
		assert.Contains(t, test.outPlatforms, pform, fmt.Sprintf("Platform should be in: %s", test.outPlatforms))
		readCloser.Close()
	}
}

func TestRefresh(t *testing.T) {
	err := cache.Refresh()
	assert.NoError(t, err, "Refreshing Cache should not throw an error")

}

func TestGetCacheDir(t *testing.T) {
	HOME := os.Getenv("HOME")
	tests := []struct {
		input  string
		output string
	}{
		{"", path.Join(HOME, ".tldr")},
		{"test", "test"},
	}
	for _, test := range tests {
		out, _ := getCacheDir(test.input)
		assert.Equal(t, test.output, out, fmt.Sprintf("Expected: %s, Actual: %s", test.output, out))
	}
}

func TestListPages(t *testing.T) {
	tests := []struct {
		platform     string
		expectations []string
	}{
		{"common", []string{"curl", "ls"}},
		{"linux", []string{"dmesg", "alpine"}},
		{"osx", []string{"dmesg", "brew"}},
		{"sunos", []string{"dmesg", "stty"}},
		{"windows", []string{"rmdir", "mkdir"}},
	}

	for _, test := range tests {
		pages, _ := cache.ListPages(test.platform)

		for _, expectation := range test.expectations {
			assert.Contains(t, pages, expectation)
		}
	}
}

func TestAvailablePlatforms(t *testing.T) {
	platforms, _ := cache.AvailablePlatforms()
	assert.Len(t, platforms, 5, "There should be 5 available platforms")
	assert.Contains(t, platforms, "common", "Platforms should contain 'common'")
	assert.Contains(t, platforms, "linux", "Platforms should contain 'linux'")
	assert.Contains(t, platforms, "osx", "Platforms should contain 'osx'")
	assert.Contains(t, platforms, "sunos", "Platforms should contain 'sunos'")
	assert.Contains(t, platforms, "windows", "Platforms should contain 'windows'")
}

func TestIsPlatformValid(t *testing.T) {
	tests := []struct {
		platform string
		isValid  bool
	}{
		{"common", true},
		{"linux", true},
		{"osx", true},
		{"sunos", true},
		{"windows", true},
		{"fake", false},
	}

	for _, test := range tests {
		isValid, _ := cache.IsPlatformValid(test.platform)
		assert.Equal(t, test.isValid, isValid)
	}
}

func TestCacheInvalidation(t *testing.T) {
	location := "./test-cache-invalidation"
	Create(remoteURL, time.Millisecond, location)
	now := time.Now()
	Create(remoteURL, time.Millisecond, location)
	info, _ := os.Stat(location)
	assert.True(t, info.ModTime().After(now))
	os.RemoveAll(location)
}

func TestPurge(t *testing.T) {
	location := "./test-purge"
	os.Mkdir(location, 0755)
	purgeCache := Cache{
		Location: location,
		TTL:      time.Second,
	}
	purgeCache.Purge()
}
