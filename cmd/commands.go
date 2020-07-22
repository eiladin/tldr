package cmd

import (
	"log"

	"github.com/eiladin/tldr/internal/config"
	"github.com/eiladin/tldr/internal/pipeline"
	"github.com/eiladin/tldr/pkg/context"
	"github.com/spf13/cobra"
)

type commandsCmd struct {
	cmd  *cobra.Command
	opts commandsOptions
}

type commandsOptions struct {
	platform string
}

func newCommandsCmd() *commandsCmd {
	var c = &commandsCmd{}
	var cmd = &cobra.Command{
		Use:           "commands",
		Short:         "list all commands for the selected platform.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(cmd *cobra.Command, args []string) {
			_, err := listCommands(c.opts, args...)
			if err != nil {
				log.Fatalf(err.Error())
			}
		},
	}

	cmd.Flags().StringVarP(&c.opts.platform, "platform", "p", config.CurrentPlatform(), "Platform to show usage for (run 'tldr platforms' to see available platforms)")
	c.cmd = cmd
	return c
}

func listCommands(options commandsOptions, args ...string) (*context.Context, error) {
	ctx := context.New()
	setupCommandContext(ctx, options)
	return pipeline.Execute(ctx, pipeline.ListCommandsPipeline)
}

func setupCommandContext(ctx *context.Context, options commandsOptions) *context.Context {
	ctx.Operation = context.OperationListCommands
	ctx.Platform = options.platform
	return ctx
}
