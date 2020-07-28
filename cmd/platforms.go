package cmd

import (
	"os"

	"github.com/eiladin/tldr/internal/pipeline"
	"github.com/eiladin/tldr/pkg/context"
	"github.com/spf13/cobra"
)

type platformsCmd struct {
	cmd *cobra.Command
}

func newPlatformsCmd() *platformsCmd {
	c := &platformsCmd{}
	cmd := &cobra.Command{
		Use:           "platforms",
		Short:         "List available platforms.",
		SilenceErrors: true,
		SilenceUsage:  true,
		Run: func(cmd *cobra.Command, args []string) {
			_, err := listPlatforms(args...)
			if err != nil {
				os.Exit(1)
			}
		},
	}
	c.cmd = cmd
	return c
}

func listPlatforms(args ...string) (*context.Context, error) {
	ctx := context.New()
	setupPlatformsContext(ctx)
	return pipeline.Execute(ctx, pipeline.ListPlatformsPipeline)
}

func setupPlatformsContext(ctx *context.Context) *context.Context {
	ctx.Operation = context.OperationListPlatforms
	return ctx
}
