package page

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

var bold = color.New(color.Bold)

func replaceTags(line *string) {
	re := regexp.MustCompile(`\{\{(.*?)\}\}`)
	keys := re.FindAllString(*line, -1)
	for _, ele := range keys {
		repl := strings.Trim(ele, "{{")
		repl = strings.Trim(repl, "}}")
		*line = strings.Replace(*line, ele, color.HiBlueString(repl), -1)
	}
}

func render(markdown io.Reader) (string, error) {
	var rendered string
	scanner := bufio.NewScanner(markdown)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			rendered += bold.Sprint(line[2:]) + "\n"
		} else if strings.HasPrefix(line, ">") {
			rendered += color.YellowString(line[2:]) + "\n"
		} else if strings.HasPrefix(line, "-") {
			rendered += color.CyanString(line) + "\n"
			scanner.Scan()
			scanner.Scan()
			line = scanner.Text()
			replaceTags(&line)
			rendered += color.WhiteString(strings.Trim(line, "`")) + "\n"
		} else {
			rendered += color.WhiteString(line) + "\n"
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
