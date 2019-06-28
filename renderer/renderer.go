package renderer

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

var (
	bold  = color.New(color.Bold)
	blue  = color.New(color.FgHiBlue)
	red   = color.New(color.FgRed)
	cyan  = color.New(color.FgCyan)
	white = color.New(color.FgWhite)
)

func Render(markdown io.Reader) (string, error) {
	var rendered string
	var renderingExample bool
	scanner := bufio.NewScanner(markdown)
	for scanner.Scan() {
		line := scanner.Text()
		if renderingExample {
			scanner.Scan()
			line = scanner.Text()

			re := regexp.MustCompile(`\{\{.*\}\}`)
			fmt.Printf("Pattern: %s", re.String())
			fmt.Println("Matched: ", re.MatchString(line))

			// line = strings.Replace(line, "{{"+arg+"}}", )

			// line = strings.Replace(line, "{{", BLUE, -1)
			// line = strings.Replace(line, "}}", RED, -1)
			// rendered += "\t" + RED + strings.Trim(line, "`") + RESET + "\n"

			renderingExample = false
		} else if strings.HasPrefix(line, "#") {
			// Heading
			rendered += bold.Sprint(line[2:]) + "\n"
		} else if strings.HasPrefix(line, ">") {
			// Quote
			rendered += line[2:] + "\n"
		} else if strings.HasPrefix(line, "-") {
			// Example
			rendered += blue.Sprintln(line)
			renderingExample = true
		} else {
			rendered += line + "\n"
		}
	}
	return rendered, scanner.Err()
}

func Write(markdown io.Reader, dest io.Writer) error {
	out, err := Render(markdown)
	if err != nil {
		return err
	}
	_, err = io.WriteString(dest, out)
	return err
}
