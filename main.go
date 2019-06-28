package main

import (
	"fmt"
	"time"

	"github.com/eiladin/tldr/cache"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	remoteURL = "http://tldr-pages.github.io/assets/tldr.zip"
	ttl       = time.Hour * 24 * 7
)

var (
	clear  = kingpin.Flag("clear-cache", fmt.Sprintf("Clear local cache then update from %s", remoteURL)).Short('c').Bool()
	update = kingpin.Flag("update", fmt.Sprintf("Clear local cache then update from %s", remoteURL)).Short('u').Bool()
	random = kingpin.Flag("raundom", "Random page for testing purposes").Short('r').Default("false").Bool()
	page   = kingpin.Arg("command", "Name of the command. (e.g. tldr tar)").Strings()
)

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("1.0").Author("Sami Khan")
	kingpin.Parse()
	cache.Create(remoteURL, ttl)
}
