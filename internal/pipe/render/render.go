package render

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/eiladin/tldr/internal/renderer"
	"github.com/eiladin/tldr/pkg/context"
)

type Pipe struct{}

func (Pipe) String() string {
	return "getting page"
}

func (Pipe) Run(ctx *context.Context) error {
	r := renderer.ColorRenderer{
		UseColor: ctx.Color,
	}
	r.Init()
	closer, err := os.Open(ctx.Page)
	if err != nil {
		return err
	}
	defer closer.Close()
	render(closer, r, ctx)
	return nil
}

func render(markdown io.Reader, r renderer.PageRenderer, ctx *context.Context) error {
	r.Init()
	scanner := bufio.NewScanner(markdown)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			io.WriteString(ctx.Writer, r.RenderTitle(line[2:]))
			io.WriteString(ctx.Writer, r.RenderPlatform(ctx.RenderPlatform()))
		} else if strings.HasPrefix(line, ">") {
			io.WriteString(ctx.Writer, r.RenderDescription(line[2:]))
		} else if strings.HasPrefix(line, "-") {
			io.WriteString(ctx.Writer, r.RenderExample(line))
			scanner.Scan()
			scanner.Scan()
			line = scanner.Text()
			io.WriteString(ctx.Writer, r.RenderSyntax(line))
		} else {
			io.WriteString(ctx.Writer, fmt.Sprintln(line))
		}
	}
	return scanner.Err()
}
