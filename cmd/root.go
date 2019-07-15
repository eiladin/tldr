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
	Short:   "Simplified and community-driven man pages",
	Long:    `Simplified and community-driven man pages`,
	Version: "1.2.9",
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
	rootCmd.Flags().StringP("platform", "p", config.CurrentPlatform(), "Platform to show usage for (run 'tldr platforms' to see available platforms)")
	rootCmd.Flags().BoolP("color", "c", true, "Pretty Print (color and formatting)")
	rootCmd.Flags().BoolP("purge", "", false, "Clear local cache")
}

// ValidateArgs checks to make sure user input is valid before execution
func ValidateArgs(cmd *cobra.Command, args []string) error {
	update, _ := cmd.Flags().GetBool("update")
	random, _ := cmd.Flags().GetBool("random")
	purge, _ := cmd.Flags().GetBool("purge")
	if !random && !update && !purge && len(args) < 1 {
		return errors.New("requires a command")
	}
	return nil
}

type flags struct {
	update   bool
	platform string
	random   bool
	color    bool
	purge    bool
}

// FindPage will look up the page and display the simplified help for the user-provided command
func FindPage(cmd *cobra.Command, args []string) {
	u, _ := cmd.Flags().GetBool("update")
	p, _ := cmd.Flags().GetString("platform")
	r, _ := cmd.Flags().GetBool("random")
	c, _ := cmd.Flags().GetBool("color")
	pu, _ := cmd.Flags().GetBool("purge")

	f := flags{
		update:   u,
		platform: p,
		random:   r,
		color:    c,
		purge:    pu,
	}
	findPage(os.Stdout, f, cache.DefaultSettings, args...)
}

var logFatalf = log.Fatalf

func purgeCache(w io.Writer, settings cache.Cache) {
	fmt.Fprintf(w, "Clearing cache ... ")
	err := settings.Purge()
	if err != nil {
		logFatalf("ERROR: removing cache: %s", err)
	}
	fmt.Fprintf(w, "Done\n")
}

func initCache(settings cache.Cache) *cache.Cache {
	c, err := cache.Create(settings.Remote, settings.TTL, settings.Location)
	if err != nil {
		logFatalf("ERROR: creating cache: %s", err)
	}
	return c
}

func validatePlatform(cache *cache.Cache, platform string) {
	platformValid, availablePlatforms := cache.IsPlatformValid(platform)
	if !platformValid {
		logFatalf("ERROR: platform %s not found\nAvailable platforms: %s", platform, strings.Join(availablePlatforms, ", "))
	}
}

func updateCache(w io.Writer, cache *cache.Cache) {
	fmt.Fprint(w, "Refreshing Cache ... ")
	cache.Refresh()
	fmt.Fprintln(w, "Done")
}

func formatPlatform(platform string, foundPlatform string) string {
	if foundPlatform != platform {
		foundPlatform = fmt.Sprintf("%s (%s)", platform, foundPlatform)
	}
	return foundPlatform
}

func printRandomPage(w io.Writer, cache *cache.Cache, f flags) {
	c, foundPlatform, err := cache.FetchRandomPage(f.platform)
	if err != nil {
		logFatalf("ERROR: getting random page: %s", err)
	}
	foundPlatform = formatPlatform(f.platform, foundPlatform)
	defer c.Close()
	write(c, w, foundPlatform, f)
}

func printPage(w io.Writer, cache *cache.Cache, f flags, page string) {
	c, foundPlatform, err := cache.FetchPage(f.platform, page)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	foundPlatform = formatPlatform(f.platform, foundPlatform)
	defer c.Close()
	write(c, w, foundPlatform, f)
}

func write(c io.ReadCloser, w io.Writer, platform string, f flags) {
	if err := page.Write(c, w, platform, f.color); err != nil {
		logFatalf("ERROR: rendering page: %s", err)
	}
}

func findPage(w io.Writer, f flags, settings cache.Cache, args ...string) {
	if f.purge {
		purgeCache(w, settings)
		return
	}

	cache := initCache(settings)

	validatePlatform(cache, f.platform)

	if f.update {
		updateCache(w, cache)
	}

	if f.random {
		printRandomPage(w, cache, f)
	} else {
		cmd := strings.Join(args, "-")
		if f.update && cmd == "" {
			return
		}
		printPage(w, cache, f, cmd)
	}
}
