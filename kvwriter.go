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

func exportKvShell(c Config) string {
	c.Name = toBashName(c.Name)
	if c.Name == "" {
		return ""
	}
	c.Value = strings.Replace(c.Value, `"`, `\"`, -1)
	c.Value = strings.Replace(c.Value, `\`, `\\`, -1)
	return fmt.Sprintf(`export %s="%s"
`, c.Name, c.Value)
}

func exportKvDocker(c Config) string {
	c.Name = toBashName(c.Name)
	if c.Name == "" {
		return ""
	}
	c.Value = strings.Replace(c.Value, `"`, `\"`, -1)
	c.Value = strings.Replace(c.Value, `\`, `\\`, -1)
	return fmt.Sprintf(`ENV %s "%s"
`, c.Name, c.Value)
}

func ExportKv(c Config, format string) string {
	if c.Type != "kv" {
		return ""
	}
	if format == "docker" {
		return exportKvDocker(c)
	} else {
		return exportKvShell(c)
	}
}
