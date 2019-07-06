package pagerenderer

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	titleText         = color.New(color.Bold)
	platformText      = color.New(color.FgHiBlack)
	tagText           = color.New(color.Italic).Add(color.FgHiBlue)
	descriptionText   = color.New(color.FgYellow)
	exampleHeaderText = color.New(color.FgHiGreen)
	exampleText       = color.New(color.FgWhite)
)

// ColorRenderer implements Renderer and prints with color and formatting
type ColorRenderer struct{}

func formatSyntaxLine(line string) string {
	formattedLine := ""
	line = strings.TrimSpace(line)
	line = strings.Replace(line, "`", "", -1)

	inTag := strings.HasPrefix(line, "{{")
	for _, segment := range strings.Split(line, "{{") {
		for _, piece := range strings.Split(segment, "}}") {
			if inTag {
				formattedLine += tagText.Sprint(piece)
			} else {
				formattedLine += exampleText.Sprint(piece)
			}
			inTag = !inTag
		}
	}

	return fmt.Sprintln(formattedLine)
}

// RenderTitle returns a formatted title
func (renderer ColorRenderer) RenderTitle(line string) string {
	return titleText.Sprintln(line)
}

// RenderPlatform returns a formatted platform
func (renderer ColorRenderer) RenderPlatform(line string) string {
	return platformText.Sprintln(line)
}

// RenderDescription returns a formatted description
func (renderer ColorRenderer) RenderDescription(line string) string {
	return descriptionText.Sprintln(line)
}

// RenderExample returns a formatted example header
func (renderer ColorRenderer) RenderExample(line string) string {
	return exampleHeaderText.Sprintln(line)
}

// RenderSyntax returns formatted example syntax
func (renderer ColorRenderer) RenderSyntax(line string) string {
	return formatSyntaxLine(line)
}
