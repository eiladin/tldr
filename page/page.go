package page

import (
	"bufio"
	"io"
	"strings"

	"github.com/fatih/color"
)

var (
	titleText         = color.New(color.Bold)
	tagText           = color.New(color.Italic).Add(color.FgHiBlue)
	descriptionText   = color.New(color.FgYellow)
	exampleHeaderText = color.New(color.FgHiGreen)
	exampleText       = color.New(color.FgWhite)
	normalText        = color.New(color.FgWhite)
)

func formatLine(line string) (formattedLine string) {
	line = strings.TrimSpace(line)
	line = strings.Replace(line, "``", "", -1)

	for _, part := range strings.Split(line, "`") {
		if len(part) == 0 {
			continue
		}
		inTag := strings.HasPrefix(part, "{{")
		for _, segment := range strings.Split(part, "{{") {
			for _, piece := range strings.Split(segment, "}}") {
				if inTag {
					formattedLine += tagText.Sprint(piece)
				} else {
					formattedLine += exampleText.Sprint(piece)
				}
				inTag = !inTag
			}
		}
	}

	return
}

func render(markdown io.Reader) (string, error) {
	var rendered string
	scanner := bufio.NewScanner(markdown)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			rendered += titleText.Sprintln(line[2:])
		} else if strings.HasPrefix(line, ">") {
			rendered += descriptionText.Sprintln(line[2:])
		} else if strings.HasPrefix(line, "-") {
			rendered += exampleHeaderText.Sprintln(line)
			scanner.Scan()
			scanner.Scan()
			line = scanner.Text()
			line = formatLine(line)
			rendered += line + "\n"
		} else {
			rendered += normalText.Sprintln(line)
		}
	}
	return rendered, scanner.Err()
}

// Write the contents of markdown to dest
func Write(markdown io.Reader, dest io.Writer) error {
	out, err := render(markdown)
	if err != nil {
		return err
	}
	_, err = io.WriteString(dest, out)
	return err
}
