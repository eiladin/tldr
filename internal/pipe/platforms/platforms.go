package platforms

import (
	"errors"
	"fmt"

	"github.com/eiladin/tldr/pkg/context"
)

var (
	errReadingPagesDir = errors.New("unable to read pages folder")
)

// Pipe for listing platforms
type Pipe struct{}

func (Pipe) String() string {
	return "listing platforms"
}

// Run the pipe
func (Pipe) Run(ctx *context.Context) error {
	for _, platform := range ctx.AvailablePlatforms {
		fmt.Fprintln(ctx.Writer, platform)
	}
	return nil
}
