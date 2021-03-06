package verify

import (
	"fmt"
	"strings"

	"github.com/eiladin/tldr/pkg/context"
)

// Pipe for verifying the selected platform
type Pipe struct{}

func (Pipe) String() string {
	return "verifying selected platform"
}

// Run the pipe
func (Pipe) Run(ctx *context.Context) error {
	valid := false
	if ctx.Operation == context.OperationListCommands && ctx.Platform == "all" {
		valid = true
	} else {
		for _, p := range ctx.AvailablePlatforms {
			if ctx.Platform == p {
				valid = true
				break
			}
		}
	}

	if !valid {
		return fmt.Errorf("ERROR: platform %s not found\nAvailable platforms: %s", ctx.Platform, strings.Join(ctx.AvailablePlatforms, ", "))
	}
	return nil
}
