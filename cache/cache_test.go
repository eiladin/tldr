package cache

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/user"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	remoteURL = "http://tldr-pages.github.com/assets/tldr.zip"
	ttl       = time.Hour * 24 * 7
	location  = "./tldr-tests"
)

var cache *Cache

func CreateCache() *Cache {
	cache, _ = Create(remoteURL, ttl, location)
	return cache
}

func DestroyCache() {
	os.RemoveAll(location)
}

func TestMain(m *testing.M) {
	CreateCache()
	code := m.Run()
	DestroyCache()
	os.Exit(code)
}

func TestFetch(t *testing.T) {
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
		readCloser, pform, _ := cache.Fetch(test.platform, test.page)
		assert.Equal(t, test.outPlatform, pform, fmt.Sprintf("Platform should match: %s", test.outPlatform))
		readCloser.Close()
	}

	_, _, err := cache.Fetch("linux", "qwaszx")
	assert.Error(t, err, "Should result in a not exist error")
}

func TestRandom(t *testing.T) {
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
		readCloser, pform, _ := cache.FetchRandom(test.platform)
		assert.Contains(t, test.outPlatforms, pform, fmt.Sprintf("Platform should be in: %s", test.outPlatforms))
		readCloser.Close()
	}
}

func TestRefresh(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cache.Refresh()

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()
	w.Close()
	os.Stdout = old
	out := <-outC
	assert.Equal(t, "Refreshing Cache ... Done\n", out, "There should a refresh cache message")

}

func TestGetCacheDir(t *testing.T) {
	usr, _ := user.Current()
	tests := []struct {
		input  string
		output string
	}{
		{"", usr.HomeDir + "/.tldr"},
		{"test", "test"},
	}
	for _, test := range tests {
		out, _ := getCacheDir(test.input)
		assert.Equal(t, test.output, out, fmt.Sprintf("Expected: %s, Actual: %s", test.output, out))
	}
}
