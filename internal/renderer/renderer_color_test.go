package renderer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestColorRenderTitle(t *testing.T) {
	renderer := New(true)
	expected := "\x1b[1;37mTitle\x1b[0m\n"
	actual := renderer.RenderTitle("Title")
	assert.Equal(t, expected, actual)
}

func TestColorRenderDescription(t *testing.T) {
	renderer := New(true)
	expected := "\x1b[37mDescription\x1b[0m\n"
	actual := renderer.RenderDescription("Description")
	assert.Equal(t, expected, actual)
}

func TestColorRenderPlatform(t *testing.T) {
	renderer := New(true)
	expected := "\x1b[90mPlatform\x1b[0m\n"
	actual := renderer.RenderPlatform("Platform")
	assert.Equal(t, expected, actual)
}

func TestColorRenderExample(t *testing.T) {
	renderer := New(true)
	expected := "\x1b[32mExample\x1b[0m\n"
	actual := renderer.RenderExample("Example")
	assert.Equal(t, expected, actual)
}

func TestColorRenderSyntax(t *testing.T) {
	renderer := New(true)
	expected := "  \x1b[34mSyntax \x1b[0m\x1b[37mexample\x1b[0m\x1b[34m\x1b[0m\n"
	actual := renderer.RenderSyntax("Syntax {{example}}")
	assert.Equal(t, expected, actual)
}

func TestFormatSyntaxLine(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Syntax {{example}}", "  \x1b[34mSyntax \x1b[0m\x1b[37mexample\x1b[0m\x1b[34m\x1b[0m\n"},
		{"Syntax {{example}}{{2}}", "  \x1b[34mSyntax \x1b[0m\x1b[37mexample2\x1b[0m\x1b[34m\x1b[0m\n"},
		{"Empty {{}}", "  \x1b[34mEmpty \x1b[0m\n"},
	}

	r := New(true)
	for _, test := range tests {
		assert.Equal(t, test.expected, r.formatSyntaxLine(test.input))
	}
}
