package random

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/eiladin/tldr/internal/pipe"
	"github.com/eiladin/tldr/pkg/context"
	"github.com/stretchr/testify/assert"
)

func initTest(dir string, cases []test) {
	d := path.Join(dir, "pages")
	for _, c := range cases {
		os.MkdirAll(path.Join(d, c.platform), 0755)
		for _, e := range c.expectations {
			os.Create(path.Join(d, c.platform, fmt.Sprintf("%s.md", e)))
		}
	}
}

func cleanTest(dir string) {
	os.RemoveAll(dir)
}

type test struct {
	platform     string
	expectations []string
}

func TestString(t *testing.T) {
	p := Pipe{}
	assert.NotEmpty(t, p.String())
}

func TestRandom(t *testing.T) {
	defer cleanTest("./test-cache")

	cases := []test{
		{"linux", []string{"dmesg", "alpine"}},
		{"osx", []string{"dmesg", "brew"}},
		{"sunos", []string{"dmesg", "stty"}},
		{"windows", []string{"rmdir", "mkdir"}},
		{"common", []string{"ls", "curl"}},
	}
	initTest("./test-cache", cases)

	ctx := context.New()
	ctx.Cache.Location = "./test-cache"
	ctx.Cache.TTL = time.Hour
	ctx.Random = true

	for _, c := range cases {
		ctx.Platform = c.platform
		err := Pipe{}.Run(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, ctx.Page)
		if ctx.FoundPlatform == "common" {
			assert.Contains(t, ctx.Page, "common")
		} else {
			assert.Contains(t, ctx.Page, c.platform)
		}
	}
}

func TestSkip(t *testing.T) {
	ctx := context.New()
	err := Pipe{}.Run(ctx)
	assert.Error(t, err)
	assert.True(t, pipe.IsSkip(err))
}

func TestDirError(t *testing.T) {
	ctx := context.New()
	ctx.Cache.Location = "."
	ctx.PagesDirectory = "pages"
	ctx.Platform = "none"
	ctx.Random = true
	err := Pipe{}.Run(ctx)
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}
