package cmd

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/eiladin/tldr/cache"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

var settings = cache.Cache{
	Remote:   "http://tldr-pages.github.com/assets/tldr.zip",
	TTL:      time.Minute,
	Location: "./tldr-cmd-test",
}

func getCobraCommand() *cobra.Command {
	c := &cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {},
	}

	c.Flags().BoolP("update", "u", false, "")
	c.Flags().BoolP("random", "r", false, "")
	c.Flags().BoolP("purge", "", false, "")
	c.Flags().StringP("platform", "p", "linux", "")
	c.Flags().BoolP("color", "c", false, "")
	return c
}

func TestValidateArgs(t *testing.T) {
	tests := []struct {
		update    string
		random    string
		purge     string
		args      []string
		shouldErr bool
	}{
		{"false", "false", "false", []string{}, true},
		{"true", "false", "false", []string{}, false},
		{"false", "true", "false", []string{}, false},
		{"false", "false", "true", []string{}, false},
		{"false", "false", "false", []string{"arg1"}, false},
	}

	c := getCobraCommand()

	for _, test := range tests {
		c.Flag("update").Value.Set(test.update)
		c.Flag("random").Value.Set(test.random)
		c.Flag("purge").Value.Set(test.purge)
		err := ValidateArgs(c, test.args)
		if test.shouldErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestCreateFlags(t *testing.T) {
	tests := []struct {
		update   string
		random   string
		purge    string
		platform string
		color    string
		flags    flags
	}{
		{"false", "false", "false", "linux", "false", flags{update: false, random: false, purge: false, platform: "linux", color: false}},
		{"true", "false", "false", "linux", "false", flags{update: true, random: false, purge: false, platform: "linux", color: false}},
		{"false", "true", "false", "linux", "false", flags{update: false, random: true, purge: false, platform: "linux", color: false}},
		{"false", "false", "true", "linux", "false", flags{update: false, random: false, purge: true, platform: "linux", color: false}},
		{"false", "false", "false", "osx", "false", flags{update: false, random: false, purge: false, platform: "osx", color: false}},
		{"false", "false", "false", "linux", "true", flags{update: false, random: false, purge: false, platform: "linux", color: true}},
	}

	c := getCobraCommand()

	for _, test := range tests {
		c.Flag("update").Value.Set(test.update)
		c.Flag("random").Value.Set(test.random)
		c.Flag("purge").Value.Set(test.purge)
		c.Flag("platform").Value.Set(test.platform)
		c.Flag("color").Value.Set(test.color)

		f := createFlags(c)
		assert.Equal(t, test.flags.update, f.update, fmt.Sprintf("update flag expected %t, but got %t", test.flags.update, f.update))
		assert.Equal(t, test.flags.random, f.random, fmt.Sprintf("random flag expected %t, but got %t", test.flags.random, f.random))
		assert.Equal(t, test.flags.purge, f.purge, fmt.Sprintf("purge flag expected %t, but got %t", test.flags.purge, f.purge))
		assert.Equal(t, test.flags.platform, f.platform, fmt.Sprintf("platform flag expected %s, but got %s", test.flags.platform, f.platform))
		assert.Equal(t, test.flags.color, f.color, fmt.Sprintf("color flag expected %t, but got %t", test.flags.color, f.color))
	}

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
