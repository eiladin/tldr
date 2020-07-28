package zip

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
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
	defer os.RemoveAll(location)
	downloadZip(location + "/tldr.zip")
	files, err := Extract(path.Join(location, "tldr.zip"), location)
	assert.NoError(t, err)
	assert.NotEmpty(t, files, "zip should contain files")
}

func TestExtractNoSource(t *testing.T) {
	const location = "./tldr-zip-nosource-test"
	files, err := Extract(path.Join(location, "/tldr.zip"), location)
	assert.Empty(t, files)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), ErrOpeningReader.Error())
}

func TestExtractDirExists(t *testing.T) {
	const location = "./tldr-zip-bad-dir-test"
	os.Mkdir(location, 0100) //nolint:errcheck
	defer func() {
		os.RemoveAll(location)
		os.Remove("./tldr.zip")
	}()
	downloadZip("./tldr.zip")
	files, err := Extract("./tldr.zip", location)
	assert.NotEmpty(t, files)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), ErrCreateOutputDir.Error())
}
