package main

import (
	"os"
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

func prepareKv(c Config) string {
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

func WriteKv(c Config) error {
	data := prepareKv(c)
	if data == "" {
		return nil
	}

	f, err := os.OpenFile(c.Path, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	if _, err = f.WriteString(data); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}
