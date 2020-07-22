package main

import (
	"os"

	"github.com/eiladin/tldr/cmd"
)

var version = "dev"

func main() {
	cmd.Execute(version, os.Args[1:])
}
