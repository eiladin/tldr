package renderer

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	titleText         = color.New(color.Bold, color.FgWhite)
	platformText      = color.New(color.FgHiBlack)
	tagText           = color.New(color.Italic, color.FgHiBlue)
	descriptionText   = color.New(color.FgHiYellow)
	exampleHeaderText = color.New(color.FgHiGreen)
	exampleText       = color.New(color.FgWhite)
)

// ColorRenderer implements Renderer and prints with color and formatting
type ColorRenderer struct{}

func colorize(str string, clr *color.Color) string {
	return clr.Sprint(str)
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
				formattedLine += colorize(piece, tagText)
			} else {
				formattedLine += colorize(piece, exampleText)
			}
			inTag = !inTag
		}
	}

	return fmt.Sprintln(formattedLine)
}

// RenderTitle returns a formatted title
func (renderer ColorRenderer) RenderTitle(line string) string {
	return fmt.Sprintln(colorize(line, titleText))
}

// RenderPlatform returns a formatted platform
func (renderer ColorRenderer) RenderPlatform(line string) string {
	return fmt.Sprintln(colorize(line, platformText))
}

// RenderDescription returns a formatted description
func (renderer ColorRenderer) RenderDescription(line string) string {
	return fmt.Sprintln(colorize(line, descriptionText))
}

// RenderExample returns a formatted example header
func (renderer ColorRenderer) RenderExample(line string) string {
	return fmt.Sprintln(colorize(line, exampleHeaderText))
}

// RenderSyntax returns formatted example syntax
func (renderer ColorRenderer) RenderSyntax(line string) string {
	return formatSyntaxLine(line)
}
