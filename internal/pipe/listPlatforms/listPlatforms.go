package listPlatforms

import (
	"fmt"

	"github.com/eiladin/tldr/pkg/context"
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
