package page

import (
	"bytes"
	"strings"
	"testing"

	"github.com/eiladin/tldr/internal/renderer"
	"github.com/stretchr/testify/assert"
)

type testRenderer struct{}

func (renderer testRenderer) Init() {}

// RenderTitle returns an unformatted title
func (renderer testRenderer) RenderTitle(line string) string {
	return line
}

// RenderPlatform returns an unfromatted platform
func (renderer testRenderer) RenderPlatform(line string) string {
	return line
}

// RenderDescription returns an unformatted description
func (renderer testRenderer) RenderDescription(line string) string {
	return line
}

// RenderExample returns an unformatted example header
func (renderer testRenderer) RenderExample(line string) string {
	return line
}

// RenderSyntax returns unformatted example syntax
func (renderer testRenderer) RenderSyntax(line string) string {
	return line
}

func TestRender(t *testing.T) {
	tests := []struct {
		input    string
		platform string
		renderer renderer.PageRenderer
		output   string
	}{
		{"# Title", "linux", new(testRenderer), "Titlelinux"},
		{"> Description", "linux", new(testRenderer), "Description"},
		{"- Example Header", "linux", new(testRenderer), "- Example Header"},
		{"normal line", "linux", new(testRenderer), "normal line\n"},
		{"- Header\n\n`test {{tag}}`", "linux", new(testRenderer), "- Header`test {{tag}}`"},
	}

	for _, test := range tests {
		markdown := strings.NewReader(test.input)
		var b bytes.Buffer
		render(markdown, &b, test.platform, test.renderer)
		assert.Equal(t, test.output, b.String())
	}
}

func TestWrite(t *testing.T) {
	tests := []struct {
		input    string
		platform string
		color    bool
		output   string
	}{
		{"# Title", "linux", true, "\x1b[1;37mTitle\x1b[0m\n\x1b[90mlinux\x1b[0m\n"},
		{"> Description", "linux", true, "\x1b[37mDescription\x1b[0m\n"},
		{"- Example Header", "linux", true, "\x1b[32m- Example Header\x1b[0m\n  \x1b[34m\x1b[0m\n"},
		{"normal line", "linux", true, "normal line\n"},
		{"- Header\n\n`test {{tag}}`", "linux", true, "\x1b[32m- Header\x1b[0m\n  \x1b[34mtest \x1b[0m\x1b[37mtag\x1b[0m\x1b[34m\x1b[0m\n"},
		{"# Title", "linux", false, "Title\nlinux\n"},
		{"> Description", "linux", false, "Description\n"},
		{"- Example Header", "linux", false, "- Example Header\n  \n"},
		{"normal line", "linux", false, "normal line\n"},
		{"- Header\n\n`test {{tag}}`", "linux", false, "- Header\n  test tag\n"},
	}

	for _, test := range tests {
		markdown := strings.NewReader(test.input)
		var b bytes.Buffer
		Write(markdown, &b, test.platform, test.color)
		assert.Equal(t, test.output, b.String())
	}
}
