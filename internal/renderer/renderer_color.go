package renderer

import (
	"fmt"
	"strings"

	"github.com/logrusorgru/aurora"
)

// ColorRenderer implements Renderer and prints with color and formatting
type ColorRenderer struct {
	Color aurora.Aurora
}

// New inits the color renderer
func New(useColor bool) ColorRenderer {
	au := aurora.NewAurora(useColor)
	return ColorRenderer{
		Color: au,
	}
}

func (r ColorRenderer) formatSyntaxLine(line string) string {
	formattedLine := "  "
	line = strings.TrimSpace(line)
	line = strings.Replace(line, "`", "", -1)
	line = strings.Replace(line, "{{}}", "", -1)
	line = strings.Replace(line, "}}{{", "", -1)

	inTag := strings.HasPrefix(line, "{{")
	for _, segment := range strings.Split(line, "{{") {
		for _, piece := range strings.Split(segment, "}}") {
			if inTag {
				formattedLine += fmt.Sprint(r.Color.White(piece))
			} else {
				formattedLine += fmt.Sprint(r.Color.Blue(piece))
			}
			inTag = !inTag
		}
	}

	return fmt.Sprintln(formattedLine)
}

// RenderTitle returns a formatted title
func (r ColorRenderer) RenderTitle(line string) string {
	return fmt.Sprintln(r.Color.Bold(r.Color.White(line)))
}

// RenderPlatform returns a formatted platform
func (r ColorRenderer) RenderPlatform(line string) string {
	return fmt.Sprintln(r.Color.BrightBlack(line))
}

// RenderDescription returns a formatted description
func (r ColorRenderer) RenderDescription(line string) string {
	return fmt.Sprintln(r.Color.White(line))
}

// RenderExample returns a formatted example header
func (r ColorRenderer) RenderExample(line string) string {
	return fmt.Sprintln(r.Color.Green(line))
}

// RenderSyntax returns formatted example syntax
func (r ColorRenderer) RenderSyntax(line string) string {
	return r.formatSyntaxLine(line)
}
