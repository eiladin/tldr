package pagerenderer

import "fmt"

// NormalRenderer implements PageRenderer and prints without color and pretty formatting
type NormalRenderer struct{}

// RenderTitle returns an unformatted title
func (renderer NormalRenderer) RenderTitle(line string) string {
	return fmt.Sprintln(line)
}

// RenderDescription returns an unformatted description
func (renderer NormalRenderer) RenderDescription(line string) string {
	return fmt.Sprintln(line)
}

// RenderExample returns an unformatted example header
func (renderer NormalRenderer) RenderExample(line string) string {
	return fmt.Sprintln(line)
}

// RenderSyntax returns unformatted example syntax
func (renderer NormalRenderer) RenderSyntax(line string) string {
	return fmt.Sprintln(line)
}
