package main

import (
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
	return newname.String()
}

func ExportKv(c Config) string {
	if c.Type != "kv" {
		return ""
	}

	// remove all space, all newline, unicode character in name
	c.Name = toBashName(c.Name)
	if c.Name == "" {
		return ""
	}

	c.Value = strings.Replace(c.Value, `"`, `\"`, -1)
	return c.Name + `="` + c.Value + `"
`
}
