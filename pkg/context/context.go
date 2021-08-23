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

// Cache struct for working with the cache
type Cache struct {
	Location string
	Remote   string
	TTL      time.Duration
}

const (
	// OperationListCommands operation
	OperationListCommands = "list-commands"
	// OperationListPlatforms operation
	OperationListPlatforms = "list-platforms"
	// OperationRenderPage operation
	OperationRenderPage = "render-page"
)

// Context carries state through pipes
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

// RenderPlatform formats the platform as `requested (actual)`
func (ctx *Context) RenderPlatform() string {
	if ctx.FoundPlatform != ctx.Platform {
		return fmt.Sprintf("%s (%s)", ctx.Platform, ctx.FoundPlatform)
	}
	return ctx.FoundPlatform
}

// New creates a new context
func New() *Context {
	location, _ := getCacheDir()
	ctx := Context{
		PagesDirectory: "pages",
		PageSuffix:     ".md",
		Writer:         os.Stdout,
		Cache: Cache{
			Location: location,
			TTL:      time.Hour * 24 * 7,
			Remote:   "https://tldr.sh/assets/tldr.zip",
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
