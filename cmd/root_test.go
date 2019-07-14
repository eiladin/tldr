package cmd

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/eiladin/tldr/cache"
	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
)

func TestFindPage(t *testing.T) {
	color.NoColor = false

	settings := cache.Cache{
		Remote:   "http://tldr-pages.github.com/assets/tldr.zip",
		TTL:      time.Minute,
		Location: "./tldr-root-cmd-test",
	}

	tests := []struct {
		update       bool
		platform     string
		random       bool
		color        bool
		args         []string
		expectations []string
	}{
		{true, "linux", false, false, []string{}, []string{"Refreshing Cache"}},
		{false, "linux", false, false, []string{"git", "pull"}, []string{"git-pull", "linux", "common"}},
		{false, "linux", true, false, []string{}, []string{"linux"}},
		{false, "linux", true, true, []string{}, []string{"\x1b"}},
		{false, "linux", false, false, []string{"qwaszx"}, []string{"This page (qwaszx) does not exist yet!"}},
	}

	for _, test := range tests {
		var b bytes.Buffer
		findPage(&b, test.update, test.platform, test.random, test.color, settings, test.args...)
		out := b.String()
		for _, expectation := range test.expectations {
			assert.Contains(t, out, expectation)
		}
	}

	os.RemoveAll(settings.Location)
}
