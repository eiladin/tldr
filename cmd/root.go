package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/eiladin/tldr/cache"
	"github.com/eiladin/tldr/config"
	"github.com/eiladin/tldr/page"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "tldr",
	Short:   "Everyday help for everyday commands",
	Long:    `Everyday help for everyday commands`,
	Version: "1.2.1",
	Args:    ValidateArgs,
	Run:     FindPage,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("update", "u", false, fmt.Sprintf("Clear local cache and update from %s", cache.DefaultSettings.Remote))
	rootCmd.Flags().BoolP("random", "r", false, "Random page for testing purposes.")
	rootCmd.Flags().StringP("platform", "p", config.OSName(), "Platform to show usage for (run 'tldr platforms' to see available platforms)")
	rootCmd.Flags().BoolP("color", "c", true, "Pretty Print (color and formatting)")
}

// ValidateArgs checks to make sure user input is valid before execution
func ValidateArgs(cmd *cobra.Command, args []string) error {
	update, _ := cmd.Flags().GetBool("update")
	random, _ := cmd.Flags().GetBool("random")
	if !random && !update && len(args) < 1 {
		return errors.New("requires a command")
	}
	return nil
}

// FindPage will look up the page and display the simplified help for the user-provided command
func FindPage(cmd *cobra.Command, args []string) {
	update, _ := cmd.Flags().GetBool("update")
	platform, _ := cmd.Flags().GetString("platform")
	random, _ := cmd.Flags().GetBool("random")
	color, _ := cmd.Flags().GetBool("color")
	findPage(os.Stdout, update, platform, random, color, cache.DefaultSettings, args...)
}

func findPage(writer io.Writer, update bool, platform string, random bool, color bool, settings cache.Cache, args ...string) {
	cache, err := cache.Create(settings.Remote, settings.TTL, settings.Location)
	platformValid := cache.IsPlatformValid(platform)
	if !platformValid {
		availablePlatforms, _ := cache.AvailablePlatforms()
		log.Fatalf("ERROR: platform %s not found\nAvailable platforms: %s", platform, strings.Join(availablePlatforms, ", "))
	}
	if err != nil {
		log.Fatalf("ERROR: creating cache: %s", err)
	}

	if update {
		cache.Refresh(writer)
	}

	var markdown io.ReadCloser
	pform := platform
	if random {
		markdown, pform, err = cache.FetchRandomPage(platform)
		if err != nil {
			log.Fatalf("ERROR: getting random page: %s", err)
		}
	} else {
		cmd := strings.Join(args, "-")
		if update && cmd == "" {
			return
		}
		markdown, pform, err = cache.FetchPage(platform, cmd)
		if err != nil {
			fmt.Fprintln(writer, err)
			return
		}
	}

	if pform != platform {
		pform = fmt.Sprintf("%s (%s)", platform, pform)
	}

	defer markdown.Close()
	if err = page.Write(markdown, writer, pform, color); err != nil {
		log.Fatalf("ERROR: rendering page: %s", err)
	}
}
