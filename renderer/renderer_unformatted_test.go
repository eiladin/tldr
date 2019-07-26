package renderer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnformattedRenderTitle(t *testing.T) {
	renderer := new(UnformattedRenderer)
	expected := "Title\n"
	actual := renderer.RenderTitle("Title")
	assert.Equal(t, expected, actual)
}

func TestUnformattedRenderDescription(t *testing.T) {
	renderer := new(UnformattedRenderer)
	expected := "Description\n"
	actual := renderer.RenderDescription("Description")
	assert.Equal(t, expected, actual)
}

func TestUnformattedRenderPlatform(t *testing.T) {
	renderer := new(UnformattedRenderer)
	expected := "Platform\n"
	actual := renderer.RenderPlatform("Platform")
	assert.Equal(t, expected, actual)
}

func TestUnformattedRenderExample(t *testing.T) {
	renderer := new(UnformattedRenderer)
	expected := "Example\n"
	actual := renderer.RenderExample("Example")
	assert.Equal(t, expected, actual)
}

func TestUnformattedRenderSyntax(t *testing.T) {
	renderer := new(UnformattedRenderer)
	expected := "  Syntax {{example}}\n"
	actual := renderer.RenderSyntax("Syntax {{example}}")
	assert.Equal(t, expected, actual)
}
