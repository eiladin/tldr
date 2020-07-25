package page

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

func initTest(cases []test) {
	dir := "./test-cache/pages"
	for _, c := range cases {
		os.MkdirAll(path.Join(dir, c.platform), 0755)
		for _, e := range c.expectations {
			os.Create(path.Join(dir, c.platform, fmt.Sprintf("%s.md", e)))
		}
	}
}

func cleanTest() {
	os.RemoveAll("./test-cache")
}

type test struct {
	platform     string
	expectations []string
}

func TestString(t *testing.T) {
	p := Pipe{}
	assert.NotEmpty(t, p.String())
}

func TestPage(t *testing.T) {
	defer cleanTest()
	cases := []test{
		{"linux", []string{"dmesg", "alpine"}},
		{"osx", []string{"dmesg", "brew"}},
		{"sunos", []string{"dmesg", "stty"}},
		{"windows", []string{"rmdir", "mkdir"}},
	}
	initTest(cases)

	for _, c := range cases {
		ctx := context.New()
		ctx.Cache.Location = "./test-cache"
		ctx.Cache.TTL = time.Minute
		ctx.Platform = c.platform
		for _, expectation := range c.expectations {
			ctx.Args = expectation
			err := Pipe{}.Run(ctx)
			assert.NoError(t, err)
			assert.Contains(t, ctx.Page, expectation)
		}
	}
}

func TestRandom(t *testing.T) {
	ctx := context.New()
	ctx.Cache.Location = "./test-cache"
	ctx.Cache.TTL = time.Minute
	ctx.Platform = "linux"
	ctx.Random = true
	err := Pipe{}.Run(ctx)
	assert.Error(t, err)
	assert.True(t, pipe.IsSkip(err))
}

func TestSkip(t *testing.T) {
	ctx := context.New()
	ctx.Cache.Location = "./test-cache"
	ctx.Cache.TTL = time.Minute
	ctx.Platform = "linux"
	err := Pipe{}.Run(ctx)
	assert.Error(t, err)
	assert.True(t, pipe.IsSkip(err))
}

func TestNotFound(t *testing.T) {
	ctx := context.New()
	ctx.Cache.Location = "./test-cache"
	ctx.Cache.TTL = time.Minute
	ctx.Platform = "linux"
	ctx.Args = "not found"
	err := Pipe{}.Run(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Submit new pages here")
}
