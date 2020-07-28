package invalidate

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

func TestPurgeCache(t *testing.T) {
	defer cleanTest("./test-cache")
	cases := []test{
		{"linux", []string{"dmesg", "alpine"}},
		{"osx", []string{"dmesg", "brew"}},
		{"sunos", []string{"dmesg", "stty"}},
		{"windows", []string{"rmdir", "mkdir"}},
	}
	initTest("./test-cache", cases)

	ctx := context.New()
	ctx.Cache.Location = "./test-cache"
	ctx.Cache.TTL = time.Hour
	ctx.PurgeCache = true
	err := Pipe{}.Run(ctx)
	assert.NoError(t, err)
	assert.NoDirExists(t, "./test-cache")
}

func TestExpiredCache(t *testing.T) {
	defer cleanTest("./expired-cache")
	cases := []test{
		{"linux", []string{"dmesg", "alpine"}},
		{"osx", []string{"dmesg", "brew"}},
		{"sunos", []string{"dmesg", "stty"}},
		{"windows", []string{"rmdir", "mkdir"}},
	}
	initTest("./expired-cache", cases)
	time.Sleep(time.Millisecond * 3)

	ctx := context.New()
	ctx.Cache.Location = "./expired-cache"
	ctx.Cache.TTL = time.Millisecond
	err := Pipe{}.Run(ctx)
	assert.NoError(t, err)
	assert.NoDirExists(t, "./expired-cache")
}

func TestSkip(t *testing.T) {
	ctx := context.New()
	ctx.Cache.Location = "./skip-test"
	ctx.Cache.TTL = time.Hour
	err := Pipe{}.Run(ctx)
	assert.True(t, pipe.IsSkip(err))
}
