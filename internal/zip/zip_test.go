package zip

import (
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
)

const remoteURL = "http://tldr-pages.github.com/assets/tldr.zip"

func downloadZip(location string) {
	dir, _ := os.Create(location)
	resp, err := http.Get(remoteURL)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	defer resp.Body.Close()
	io.Copy(dir, resp.Body) //nolint:errcheck
}

func TestExtract(t *testing.T) {
	const location = "./tldr-zip-test"
	os.Mkdir(location, 0755) //nolint:errcheck
	downloadZip(location + "/tldr.zip")
	files, err := Extract(location+"/tldr.zip", location)
	assert.NoError(t, err)
	assert.NotEmpty(t, files, "zip should contain files")
	os.RemoveAll(location)
}

func TestExtractNoSource(t *testing.T) {
	const location = "./tldr-zip-nosource-test"
	files, err := Extract(location+"/tldr.zip", location)
	assert.Empty(t, files)
	assert.True(t, xerrors.Is(err, ErrOpeningReader))
}

func TestExtractDirExists(t *testing.T) {
	const location = "./tldr-zip-bad-dir-test"
	os.Mkdir(location, 0100) //nolint:errcheck
	downloadZip("./tldr.zip")
	files, err := Extract("./tldr.zip", location)
	assert.NotEmpty(t, files)
	assert.True(t, xerrors.Is(err, ErrCreateOutputDir))
	os.RemoveAll(location)
	os.RemoveAll("./tldr.zip")
}