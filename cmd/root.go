package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/eiladin/tldr/cache"
	"github.com/eiladin/tldr/config"
	"github.com/eiladin/tldr/page"
	"github.com/spf13/cobra"
)

const (
	remoteURL = "http://tldr-pages.github.com/assets/tldr.zip"
	ttl       = time.Hour * 24 * 7
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:     "tldr",
	Short:   "Everyday help for everyday commands",
	Long:    `Everyday help for everyday commands`,
	Version: "1.1",
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
	rootCmd.Flags().BoolP("update", "u", false, fmt.Sprintf("Clear local cache and update from %s", remoteURL))
	rootCmd.Flags().BoolP("random", "r", false, "Random page for testing purposes.")
	rootCmd.Flags().StringP("platform", "p", config.OSName(), "Platform to show usage for (linux, osx, sunos, windows, common)")
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
	findPage(update, platform, random, color, args...)
}

func findPage(update bool, platform string, random bool, color bool, args ...string) {
	cache, err := cache.Create(remoteURL, ttl, "")
	if err != nil {
		log.Fatalf("ERROR: creating cache: %s", err)
	}

	if update {
		cache.Refresh()
	}

	var markdown io.ReadCloser
	pform := platform
	if random {
		markdown, pform, err = cache.FetchRandom(platform)
		if err != nil {
			log.Fatalf("ERROR: getting random page: %s", err)
		}
	} else {
		cmd := strings.Join(args, "-")
		if update && cmd == "" {
			return
		}
		markdown, pform, err = cache.Fetch(platform, cmd)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if pform != platform {
		pform = fmt.Sprintf("%s (%s)", platform, pform)
	}

	defer markdown.Close()
	if err = page.Write(markdown, os.Stdout, pform, color); err != nil {
		log.Fatalf("ERROR: rendering page: %s", err)
	}
}
