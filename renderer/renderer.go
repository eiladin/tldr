package renderer

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

var (
	bold   = color.New(color.Bold)
	blue   = color.New(color.FgHiBlue)
	cyan   = color.New(color.FgHiCyan)
	white  = color.New(color.FgWhite)
	yellow = color.New(color.FgYellow)
)

func replaceTags(line *string) {
	re := regexp.MustCompile(`\{\{(.*?)\}\}`)
	keys := re.FindAllString(*line, -1)
	for _, ele := range keys {
		repl := strings.Trim(ele, "{{")
		repl = strings.Trim(repl, "}}")
		*line = strings.Replace(*line, ele, blue.Sprint(repl), -1)
	}
}

func render(markdown io.Reader) (string, error) {
	var rendered string
	var renderingExample bool
	scanner := bufio.NewScanner(markdown)
	for scanner.Scan() {
		line := scanner.Text()
		if renderingExample {
			scanner.Scan()
			line = scanner.Text()
			replaceTags(&line)
			rendered += white.Sprint(strings.Trim(line, "`")) + "\n"
			renderingExample = false
		} else if strings.HasPrefix(line, "#") {
			rendered += bold.Sprint(line[2:]) + "\n"
		} else if strings.HasPrefix(line, ">") {
			rendered += yellow.Sprint(line[2:]) + "\n"
		} else if strings.HasPrefix(line, "-") {
			rendered += cyan.Sprintln(line)
			renderingExample = true
		} else {
			rendered += line + "\n"
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
