package main

import (
	"fmt"
	"strings"
)

func toBashName(name string) string {
	const chars = "abcdefghijklmnopqrstuvwxyz_0123456789"
	newname := strings.Builder{}
	for _, r := range name {
		if !strings.Contains(chars, strings.ToLower(string(r))) {
			continue
		}
		newname.Write([]byte(string(r)))
	}
	return strings.TrimSpace(newname.String())
}

func unescapeString(content string) string {
	content = strings.Replace(content, `\`, `\\`, -1)
	content = strings.Replace(content, `"`, `\"`, -1)
	content = strings.Replace(content, "\a", `\a`, -1)
	content = strings.Replace(content, "\b", `\b`, -1)
	content = strings.Replace(content, "\r", `\r`, -1)
	content = strings.Replace(content, "\f", `\f`, -1)
	content = strings.Replace(content, "\t", `\t`, -1)
	content = strings.Replace(content, "\v", `\v`, -1)
	return content
}

func exportKv(prefix string, c Config) string {
	c.Name = toBashName(c.Name)
	if c.Name == "" {
		return ""
	}
	c.Value = unescapeString(c.Value)
	return fmt.Sprintf(`%s %s="%s"
`, prefix, c.Name, c.Value)
}

func ExportKv(c Config, format string) string {
	if c.Type != "kv" {
		return ""
	}
	if format == "docker" {
		return exportKv("ENV", c)
	} else {
		return exportKv("export", c)
	}
}
