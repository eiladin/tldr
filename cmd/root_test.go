package cmd

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/eiladin/tldr/cache"
	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
)

var settings = cache.Cache{
	Remote:   "http://tldr-pages.github.com/assets/tldr.zip",
	TTL:      time.Minute,
	Location: "./tldr-cmd-test",
}

func TestPurgeCache(t *testing.T) {
	settings.Location = "./tldr-purge-test"
	os.Mkdir(settings.Location, 0755)
	var b bytes.Buffer
	purgeCache(&b, settings)
	assert.Equal(t, "Clearing cache ... Done\n", b.String())
	f, err := os.Stat(settings.Location)
	assert.Nil(t, f)
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}

func TestInitCache(t *testing.T) {
	settings.Location = "./tldr-init-test"
	c := initCache(settings)
	assert.NotNil(t, c)
	f, err := os.Stat(settings.Location)
	assert.NotNil(t, f)
	assert.NoError(t, err)
	os.RemoveAll(settings.Location)
}

func TestValidatePlatform(t *testing.T) {
	tests := []struct {
		platform  string
		shouldErr bool
	}{
		{"common", false},
		{"linux", false},
		{"osx", false},
		{"sunos", false},
		{"windows", false},
		{"darwin", true},
	}

	settings.Location = "./tldr-validate-platform-test"
	c := initCache(settings)
	origLogFatalf := logFatalf
	defer func() { logFatalf = origLogFatalf }()
	errors := []string{}
	logFatalf = func(format string, args ...interface{}) {
		if len(args) > 0 {
			errors = append(errors, fmt.Sprintf(format, args))
		} else {
			errors = append(errors, format)
		}
	}

	for _, test := range tests {
		errors = []string{}
		validatePlatform(c, test.platform)
		if test.shouldErr {
			assert.True(t, len(errors) == 1)
		} else {
			assert.True(t, len(errors) == 0)
		}
	}
	os.RemoveAll(settings.Location)
}

func TestUpdateCache(t *testing.T) {
	settings.Location = "./tldr-update-cache-test"
	var b bytes.Buffer
	updateCache(&b, &settings)
	assert.Equal(t, "Refreshing Cache ... Done\n", b.String())
	os.RemoveAll(settings.Location)
}

func TestFormatPlatform(t *testing.T) {
	tests := []struct {
		platform      string
		foundPlatform string
		expected      string
	}{
		{"linux", "linux", "linux"},
		{"linux", "common", "linux (common)"},
	}

	for _, test := range tests {
		actual := formatPlatform(test.platform, test.foundPlatform)
		assert.Equal(t, test.expected, actual)
	}
}

func TestPrintRandomPage(t *testing.T) {
	color.NoColor = false
	tests := []struct {
		f        flags
		expected string
	}{
		{flags{false, "linux", true, true, false}, "\x1b"},
		{flags{false, "linux", true, false, false}, "linux"},
	}

	settings.Location = "./tldr-random-page-test"
	cache := initCache(settings)

	for _, test := range tests {
		var b bytes.Buffer
		printRandomPage(&b, cache, test.f)
		assert.Contains(t, b.String(), test.expected)
	}
	var buf bytes.Buffer
	purgeCache(&buf, *cache)
}

func TestPrintPage(t *testing.T) {
	color.NoColor = false
	tests := []struct {
		f        flags
		page     string
		expected []string
	}{
		{flags{false, "linux", false, true, false}, "git-pull", []string{"\x1b", "git-pull"}},
		{flags{false, "linux", false, false, false}, "git-pull", []string{"{{branch}}", "git-pull"}},
		{flags{false, "linux", false, false, false}, "not-found", []string{"This page (not-found) does not exist yet!"}},
	}

	settings.Location = "./tldr-print-page-test"
	cache := initCache(settings)

	for _, test := range tests {
		var b bytes.Buffer
		printPage(&b, cache, test.f, test.page)
		for _, expected := range test.expected {
			assert.Contains(t, b.String(), expected)
		}
	}
	var buf bytes.Buffer
	purgeCache(&buf, *cache)
}

func TestFindPage(t *testing.T) {
	color.NoColor = false
	settings.Location = "./tldr-root-cmd-test"

	tests := []struct {
		f            flags
		args         []string
		expectations []string
	}{
		{flags{true, "linux", false, false, false}, []string{}, []string{"Refreshing Cache"}},
		{flags{false, "linux", false, false, false}, []string{"git", "pull"}, []string{"git-pull", "linux", "common"}},
		{flags{false, "linux", true, false, false}, []string{}, []string{"linux"}},
		{flags{false, "linux", true, true, false}, []string{}, []string{"\x1b"}},
		{flags{false, "linux", false, false, false}, []string{"qwaszx"}, []string{"This page (qwaszx) does not exist yet!"}},
		{flags{false, "linux", false, false, true}, []string{}, []string{"Clearing cache ... Done\n"}},
	}

	for _, test := range tests {
		var b bytes.Buffer
		findPage(&b, test.f, settings, test.args...)
		out := b.String()
		for _, expectation := range test.expectations {
			assert.Contains(t, out, expectation)
		}
	}
}
