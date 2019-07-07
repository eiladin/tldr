package zip

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const remoteURL = "http://tldr-pages.github.com/assets/tldr.zip"

func TestExtract(t *testing.T) {
	const location = "./tldr-test"
	os.Mkdir(location, 0755)
	dir, _ := os.Create(location + "/tldr.zip")
	resp, err := http.Get(remoteURL)
	assert.NoError(t, err, "Error downloading zip")
	defer resp.Body.Close()
	io.Copy(dir, resp.Body)
	files, _ := Extract(location+"/tldr.zip", location)
	assert.NotEmpty(t, files, "zip should contain files")
	os.RemoveAll(location)
}
