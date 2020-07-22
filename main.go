package main

import (
	"os"

	"github.com/eiladin/tldr/cmd"
)

func main() {
	cmd.Execute(os.Args[1:])
}
