package render

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/eiladin/tldr/internal/pipe"
	"github.com/eiladin/tldr/pkg/context"
	"github.com/stretchr/testify/assert"
)

var testData = []byte("# cat\n\n" +
	"> Print and concatenate files.\n\n" +
	"- Print the contents of a file to the standard output:\n\n" +
	"`cat {{file}}`\n\n" +
	"- Concatenate several files into the target file:\n\n" +
	"`cat {{file1}} {{file2}} > {{target_file}}`\n\n" +
	"- Append several files into the target file:\n\n" +
	"`cat {{file1}} {{file2}} >> {{target_file}}`\n\n" +
	"- Number all output lines:\n\n" +
	"`cat -n {{file}}`\n\n" +
	"- Display non-printable and whitespace characters (with `M-` prefix if non-ASCII):\n\n" +
	"`cat -v -t -e {{file}}`")

var expectation = "cat\n" +
	"linux\n\n" +
	"Print and concatenate files.\n\n" +
	"- Print the contents of a file to the standard output:\n" +
	"  cat file\n\n" +
	"- Concatenate several files into the target file:\n" +
	"  cat file1 file2 > target_file\n\n" +
	"- Append several files into the target file:\n" +
	"  cat file1 file2 >> target_file\n\n" +
	"- Number all output lines:\n" +
	"  cat -n file\n\n" +
	"- Display non-printable and whitespace characters (with `M-` prefix if non-ASCII):\n" +
	"  cat -v -t -e file\n"

func TestString(t *testing.T) {
	p := Pipe{}
	assert.NotEmpty(t, p.String())
}

func TestRender(t *testing.T) {
	var b bytes.Buffer
	os.MkdirAll("./test-cache/pages/linux", 0755)

	err := ioutil.WriteFile("./test-cache/pages/linux/cat.md", testData, 0644)
	assert.NoError(t, err)
	ctx := context.New()
	ctx.Cache.Location = "./test-cache"
	ctx.Cache.TTL = time.Hour
	ctx.Color = false
	ctx.Writer = &b
	ctx.Platform = "linux"
	ctx.FoundPlatform = "linux"
	ctx.Page = "./test-cache/pages/linux/cat.md"
	err = Pipe{}.Run(ctx)
	out := b.String()
	assert.NoError(t, err)
	assert.Equal(t, out, expectation)
	os.RemoveAll("./test-cache")
}

func TestSkip(t *testing.T) {
	ctx := context.New()
	err := Pipe{}.Run(ctx)
	assert.True(t, pipe.IsSkip(err))
}

func TestDirError(t *testing.T) {
	ctx := context.New()
	ctx.Page = "not-found.page"
	err := Pipe{}.Run(ctx)
	assert.True(t, os.IsNotExist(err))
}
