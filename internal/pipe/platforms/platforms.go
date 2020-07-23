package platforms

import (
	"errors"
	"fmt"

	"github.com/eiladin/tldr/pkg/context"
)

var (
	errReadingPagesDir = errors.New("unable to read pages folder")
)

type Pipe struct{}

func (Pipe) String() string {
	return "listing platforms"
}

func (Pipe) Run(ctx *context.Context) error {
	for _, platform := range ctx.AvailablePlatforms {
		fmt.Fprintln(ctx.Writer, platform)
	}
	return nil
}
