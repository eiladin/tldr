package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/eiladin/tldr/cache"
	"github.com/stretchr/testify/assert"
)

func TestFindPage(t *testing.T) {
	settings := cache.Cache{
		Remote:   "http://tldr-pages.github.com/assets/tldr.zip",
		Ttl:      time.Minute,
		Location: "./tldr-test",
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
		fmt.Println(out)
		for _, expectation := range test.expectations {
			assert.Contains(t, out, expectation)
		}
	}

	os.RemoveAll("./tldr-test")
}
