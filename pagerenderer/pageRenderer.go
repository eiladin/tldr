package pagerenderer

// PageRenderer is an interface for rendering the output of a page
type PageRenderer interface {
	RenderTitle(string) string
	RenderPlatform(string) string
	RenderDescription(string) string
	RenderExample(string) string
	RenderSyntax(string) string
}
