package renderer

import (
	"fmt"
	"strings"
)

// UnformattedRenderer implements PageRenderer and prints without color and pretty formatting
type UnformattedRenderer struct{}

// RenderTitle returns an unformatted title
func (renderer UnformattedRenderer) RenderTitle(line string) string {
	return fmt.Sprintln(line)
}

// RenderPlatform returns an unfromatted platform
func (renderer UnformattedRenderer) RenderPlatform(line string) string {
	return fmt.Sprintln(line)
}

// RenderDescription returns an unformatted description
func (renderer UnformattedRenderer) RenderDescription(line string) string {
	return fmt.Sprintln(line)
}

// RenderExample returns an unformatted example header
func (renderer UnformattedRenderer) RenderExample(line string) string {
	return fmt.Sprintln(line)
}

// RenderSyntax returns unformatted example syntax
func (renderer UnformattedRenderer) RenderSyntax(line string) string {
	line = strings.Replace(line, "`", "", -1)
	return fmt.Sprintln(line)
}
