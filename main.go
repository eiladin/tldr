package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/eiladin/tldr/cache"
	"github.com/eiladin/tldr/config"
	"github.com/eiladin/tldr/renderer"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	remoteURL = "http://tldr-pages.github.io/assets/tldr.zip"
	ttl       = time.Hour * 24 * 7
)

var (
	update   = kingpin.Flag("update", fmt.Sprintf("Clear local cache then update from %s", remoteURL)).Short('u').Bool()
	platform = kingpin.Flag("platform", "Platform to show usage for (linux, osx, sunos, common)").Short('p').String()
	random   = kingpin.Flag("random", "Random page for testing purposes").Short('r').Default("false").Bool()
	pages    = kingpin.Arg("command", "Name of the command. (e.g. tldr tar)").Strings()
)

func main() {
	kingpin.HelpFlag.Short('h')
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("1.0").Author("Sami Khan")
	kingpin.CommandLine.Help = "Everyday help for everyday commands"

	cache, err := cache.Create(remoteURL, ttl)
	if err != nil {
		log.Fatalf("ERROR: creating cache: %s", err)
	}

	kingpin.Parse()

	osName := config.OSName()
	if *platform != "" {
		osName = *platform
	}

	if *update {
		fmt.Println("Refreshing Cache")
		cache.Refresh()
	} else if *random {
		markdown, err := cache.FetchRandom(osName)
		if err != nil {
			fmt.Println(err)
		}
		if markdown != nil {
			defer markdown.Close()
		}
		err = renderer.Write(markdown, os.Stdout)
		if err != nil {
			log.Fatalf("ERROR: rendering markdown: %s", err)
		}
	} else {
		page := ""
		for i, l := range *pages {
			if i == len(*pages)-1 {
				page = page + l
				break
			} else {
				page = page + l + "-"
			}
		}
		if page == "" {
			kingpin.Fatalf("required argument 'command' not provided, try --help")
			return
		}
		markdown, err := cache.Fetch(osName, page)
		if err != nil {
			fmt.Println(err)
			return
		}
		if markdown != nil {
			defer markdown.Close()
		}
		err = renderer.Write(markdown, os.Stdout)
		if err != nil {
			log.Fatalf("ERROR: rendering markdown: %s", err)
		}
	}
}
