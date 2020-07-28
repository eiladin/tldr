package pipeline

import (
	"fmt"

	"github.com/eiladin/tldr/internal/middleware"
	"github.com/eiladin/tldr/internal/pipe/commands"
	"github.com/eiladin/tldr/internal/pipe/initialize"
	"github.com/eiladin/tldr/internal/pipe/invalidate"
	"github.com/eiladin/tldr/internal/pipe/page"
	"github.com/eiladin/tldr/internal/pipe/platforms"
	"github.com/eiladin/tldr/internal/pipe/random"
	"github.com/eiladin/tldr/internal/pipe/render"
	"github.com/eiladin/tldr/internal/pipe/verify"
	"github.com/eiladin/tldr/pkg/context"
)

// Piper defines a pipe, which can be part of a pipeline (a series of pipes).
type Piper interface {
	fmt.Stringer
	Run(ctx *context.Context) error
}

// Execute runs a given pipeline with logging and error handling
func Execute(ctx *context.Context, pipeline []Piper) (*context.Context, error) {
	for _, pipe := range pipeline {
		if err := middleware.Logging(
			pipe.String(),
			middleware.ErrHandler(pipe.Run),
		)(ctx); err != nil {
			return ctx, err
		}
	}
	return ctx, nil
}

// RenderPipeline contains all pipe implementations to render pages in order
var RenderPipeline = []Piper{
	invalidate.Pipe{},
	initialize.Pipe{},
	verify.Pipe{},
	random.Pipe{},
	page.Pipe{},
	render.Pipe{},
}

// ListCommandsPipeline contains all pipe implementations to list commands in order
var ListCommandsPipeline = []Piper{
	invalidate.Pipe{},
	initialize.Pipe{},
	verify.Pipe{},
	commands.Pipe{},
}

// ListPlatformsPipeline contains all pipe implementations to list platforms in order
var ListPlatformsPipeline = []Piper{
	invalidate.Pipe{},
	initialize.Pipe{},
	platforms.Pipe{},
}
