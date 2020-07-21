package renderer

import (
	"fmt"
	"strings"

	"github.com/logrusorgru/aurora"
)

var (
	au = aurora.NewAurora(false)
)

func titleText(arg interface{}) aurora.Value {
	return au.Bold(au.White(arg))
}

func platformText(arg interface{}) aurora.Value {
	return au.BrightBlack(arg)
}

func tagText(arg interface{}) aurora.Value {
	return au.White(arg)
}

func descriptionText(arg interface{}) aurora.Value {
	return au.White(arg)
}

func exampleHeaderText(arg interface{}) aurora.Value {
	return au.Green(arg)
}

func exampleText(arg interface{}) aurora.Value {
	return au.Blue(arg)
}

// ColorRenderer implements Renderer and prints with color and formatting
type ColorRenderer struct {
	UseColor bool
}

func (renderer ColorRenderer) Init() {
	au = aurora.NewAurora(renderer.UseColor)
}

func formatSyntaxLine(line string) string {
	formattedLine := "  "
	line = strings.TrimSpace(line)
	line = strings.Replace(line, "`", "", -1)
	line = strings.Replace(line, "{{}}", "", -1)
	line = strings.Replace(line, "}}{{", "", -1)

	inTag := strings.HasPrefix(line, "{{")
	for _, segment := range strings.Split(line, "{{") {
		for _, piece := range strings.Split(segment, "}}") {
			if inTag {
				formattedLine += fmt.Sprint(tagText(piece))
			} else {
				formattedLine += fmt.Sprint(exampleText(piece))
			}
			inTag = !inTag
		}
	}

	return fmt.Sprintln(formattedLine)
}

// RenderTitle returns a formatted title
func (renderer ColorRenderer) RenderTitle(line string) string {
	return fmt.Sprintln(titleText(line))
}

// RenderPlatform returns a formatted platform
func (renderer ColorRenderer) RenderPlatform(line string) string {
	return fmt.Sprintln(platformText(line))
}

// RenderDescription returns a formatted description
func (renderer ColorRenderer) RenderDescription(line string) string {
	return fmt.Sprintln(descriptionText(line))
}

// RenderExample returns a formatted example header
func (renderer ColorRenderer) RenderExample(line string) string {
	return fmt.Sprintln(exampleHeaderText(line))
}

// RenderSyntax returns formatted example syntax
func (renderer ColorRenderer) RenderSyntax(line string) string {
	return formatSyntaxLine(line)
}
