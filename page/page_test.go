package page

import (
	"bytes"
	"strings"
	"testing"

	"github.com/eiladin/tldr/renderer"
	"github.com/stretchr/testify/assert"
)

type testRenderer struct{}

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
