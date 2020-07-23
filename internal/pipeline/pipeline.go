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

type Piper interface {
	fmt.Stringer
	Run(ctx *context.Context) error
}

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

var RenderPipeline = []Piper{
	invalidate.Pipe{},
	initialize.Pipe{},
	verify.Pipe{},
	random.Pipe{},
	page.Pipe{},
	render.Pipe{},
}

var ListCommandsPipeline = []Piper{
	invalidate.Pipe{},
	initialize.Pipe{},
	verify.Pipe{},
	commands.Pipe{},
}

var ListPlatformsPipeline = []Piper{
	invalidate.Pipe{},
	initialize.Pipe{},
	platforms.Pipe{},
}
