package page

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/eiladin/tldr/internal/renderer"
)

func render(markdown io.Reader, dest io.Writer, platform string, r renderer.PageRenderer) error {
	r.Init()
	scanner := bufio.NewScanner(markdown)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			io.WriteString(dest, r.RenderTitle(line[2:]))
			io.WriteString(dest, r.RenderPlatform(platform))
		} else if strings.HasPrefix(line, ">") {
			io.WriteString(dest, r.RenderDescription(line[2:]))
		} else if strings.HasPrefix(line, "-") {
			io.WriteString(dest, r.RenderExample(line))
			scanner.Scan()
			scanner.Scan()
			line = scanner.Text()
			io.WriteString(dest, r.RenderSyntax(line))
		} else {
			io.WriteString(dest, fmt.Sprintln(line))
		}
	}
	return scanner.Err()
}

// Write the contents of markdown to dest
func Write(markdown io.Reader, dest io.Writer, platform string, color bool) error {
	r := renderer.ColorRenderer{
		UseColor: color,
	}
	return render(markdown, dest, platform, r)
}
