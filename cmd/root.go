package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/eiladin/tldr/internal/config"
	"github.com/eiladin/tldr/internal/pipeline"
	"github.com/eiladin/tldr/pkg/context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type options struct {
	update   bool
	platform string
	random   bool
	color    bool
	purge    bool
}

var opts options

var rootCmd = &cobra.Command{
	Use:     "tldr",
	Short:   "Simplified and community-driven man pages",
	Long:    `Simplified and community-driven man pages`,
	Version: "1.3.11",
	Args:    validateArgs,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := findPage(args...)
		if err != nil {
			log.Fatalf(err.Error())
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var debug bool

func init() {
	cobra.OnInitialize(func() {
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}
	})
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug mode")
	rootCmd.Flags().BoolVarP(&opts.update, "update", "u", false, "update local cache")
	rootCmd.Flags().BoolVarP(&opts.random, "random", "r", false, "random page for testing purposes")
	rootCmd.Flags().StringVarP(&opts.platform, "platform", "p", config.CurrentPlatform(), "platform to show usage for (run 'tldr platforms' to see available platforms)")
	rootCmd.Flags().BoolVarP(&opts.color, "color", "c", true, "pretty print (color and formatting)")
	rootCmd.Flags().BoolVarP(&opts.purge, "purge", "", false, "clear local cache")
}

// ValidateArgs checks to make sure user input is valid before execution
func validateArgs(cmd *cobra.Command, args []string) error {
	update, _ := cmd.Flags().GetBool("update")
	random, _ := cmd.Flags().GetBool("random")
	purge, _ := cmd.Flags().GetBool("purge")
	if !random && !update && !purge && len(args) < 1 {
		return errors.New("requires a command")
	}
	return nil
}

func findPage(args ...string) (*context.Context, error) {
	ctx := context.New()
	setupContext(ctx, args...)
	return pipeline.Execute(ctx, pipeline.RenderPipeline)
}

func setupContext(ctx *context.Context, args ...string) {
	ctx.PurgeCache = opts.purge
	ctx.Platform = opts.platform
	ctx.Random = opts.random
	ctx.Color = opts.color
	ctx.Args = strings.Join(args, "-")
}
