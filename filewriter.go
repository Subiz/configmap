package main

import (
	"path"
	"strings"
)

// convert content to pair (format, param) used in `printf format param`
// example:
// + toPrintf("hello")    // printf %s "hello"
// + toPrintf("hel\nlo")  // printf %s\n%s "hel" "lo"
func toPrintf(content string) string {
	content = unescapeString(content)

	line := strings.Split(content, "\n")
	content = "\"" + strings.Replace(content, "\n", "\" \"", -1) + "\""
	formats := make([]string, 0)
	for range line {
		formats = append(formats, "%s")
	}
	format := strings.Join(formats, "\\n")
	return "printf \"" + format + "\" " + content
}

func writeFile(prefix, path, dir, content string) string {
	printfcmd := toPrintf(content)

	return prefix + " mkdir -p " + dir + " && " + printfcmd + " > " + path + "\n"
}

func WriteFile(c Config, format string) string {
	dir := path.Dir(c.Path)
	if format == "docker" {
		return writeFile("RUN", c.Path, dir, c.Value)
	} else {
		return writeFile("", c.Path, dir, c.Value)
	}
}
