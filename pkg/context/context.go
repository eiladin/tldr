package context

import (
	ctx "context"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/mitchellh/go-homedir"
)

type Cache struct {
	Location string
	Remote   string
	TTL      time.Duration
}

const (
	OperationListCommands  = "list-commands"
	OperationListPlatforms = "list-platforms"
	OperationRenderPage    = "render-page"
)

type Context struct {
	ctx.Context
	PurgeCache         bool
	Platform           string
	Random             bool
	Cache              Cache
	PagesDirectory     string
	PageSuffix         string
	FoundPlatform      string
	Reader             io.ReadCloser
	Page               string
	Args               string
	AvailablePlatforms []string
	Color              bool
	Writer             io.Writer
	Operation          string
}

func (ctx *Context) RenderPlatform() string {
	if ctx.FoundPlatform != ctx.Platform {
		return fmt.Sprintf("%s (%s)", ctx.Platform, ctx.FoundPlatform)
	}
	return ctx.FoundPlatform
}

func New() *Context {
	location, _ := getCacheDir()
	ctx := Context{
		PagesDirectory: "pages",
		PageSuffix:     ".md",
		Writer:         os.Stdout,
		Cache: Cache{
			Location: location,
			TTL:      time.Hour * 24 * 7,
			Remote:   "http://tldr-pages.github.com/assets/tldr.zip",
		},
		Operation: OperationRenderPage,
	}
	return &ctx
}

func getCacheDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return path.Join(home, ".tldr"), nil
}
