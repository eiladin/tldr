package cmd

import (
	"bytes"
	"io"
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
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		findPage(test.update, test.platform, test.random, test.color, settings, test.args...)

		outC := make(chan string)
		go func() {
			var buf bytes.Buffer
			io.Copy(&buf, r)
			outC <- buf.String()
		}()
		w.Close()
		out := <-outC
		os.Stdout = old
		for _, expectation := range test.expectations {
			assert.Contains(t, out, expectation)
		}
	}

	os.RemoveAll(settings.Location)
}
