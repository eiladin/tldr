package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/eiladin/tldr/internal/config"
	"github.com/eiladin/tldr/internal/pipeline"
	"github.com/eiladin/tldr/pkg/context"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type options struct {
	update   bool
	platform string
	random   bool
	color    bool
	purge    bool
}

func Execute(args []string) {
	newRootCmd().Execute(args)
}

func (cmd *rootCmd) Execute(args []string) {
	cmd.cmd.SetArgs(args)

	if err := cmd.cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type rootCmd struct {
	cmd   *cobra.Command
	debug bool
	opts  options
}

func newRootCmd() *rootCmd {
	var root = &rootCmd{}
	var cmd = &cobra.Command{
		Use:           "tldr",
		Short:         "Simplified and community-driven man pages",
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       "1.3.11",
		Args: func(cmd *cobra.Command, args []string) error {
			return validateArgs(root.opts, args)
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if root.debug {
				log.SetLevel(log.DebugLevel)
				log.Debug("debug logs enabled")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			_, err := findPage(root.opts, args...)
			if err != nil {
				log.Fatalf(err.Error())
			}
		},
	}

	cmd.PersistentFlags().BoolVarP(&root.debug, "debug", "d", false, "debug mode")
	cmd.Flags().BoolVarP(&root.opts.update, "update", "u", false, "update local cache")
	cmd.Flags().BoolVarP(&root.opts.random, "random", "r", false, "random page for testing purposes")
	cmd.Flags().StringVarP(&root.opts.platform, "platform", "p", config.CurrentPlatform(), "platform to show usage for (run 'tldr platforms' to see available platforms)")
	cmd.Flags().BoolVarP(&root.opts.color, "color", "c", true, "pretty print (color and formatting)")
	cmd.Flags().BoolVarP(&root.opts.purge, "purge", "", false, "clear local cache")

	cmd.AddCommand(
		newCompletionCmd().cmd,
		newCommandsCmd().cmd,
		newPlatformsCmd().cmd,
	)

	root.cmd = cmd
	return root
}

func findPage(opts options, args ...string) (*context.Context, error) {
	ctx := context.New()
	setupContext(ctx, opts, args...)
	return pipeline.Execute(ctx, pipeline.RenderPipeline)
}

func setupContext(ctx *context.Context, opts options, args ...string) {
	ctx.PurgeCache = opts.purge
	ctx.Platform = opts.platform
	ctx.Random = opts.random
	ctx.Color = opts.color
	ctx.Args = strings.Join(args, "-")
}

func validateArgs(opts options, args []string) error {
	if !opts.random && !opts.update && !opts.purge && len(args) < 1 {
		return errors.New("Command required")
	}
	return nil

}
