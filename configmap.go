package main

import (
	"fmt"
	"strings"
	"github.com/thanhpk/stringf"
)

type Config struct {
	Name       string
	Path       string
	Type       string // file, kv
	Value      string
	VaultPath  string
	VaultField string
	VaultValue *string
}

func getSubstitions(envs []string) map[string]string {
	m := make(map[string]string)
	for _, pair := range envs {// os.Environ() {
		pairsplit := strings.Split(pair, "=")
		name := strings.TrimSpace(pairsplit[0])
		if !strings.HasPrefix(name, "_") {
			continue
		}
		val := ""
		if len(pairsplit) > 1 {
			val = strings.Join(pairsplit[1:], "=")
		}
		name = strings.ToLower(name[1:])
		m[name] = val
	}
	return m
}


func extractPathAndField(key string, envs []string) (string, string) {
	arrs := strings.Split(key, "(")
	if len(arrs) < 2 {
		return arrs[0], ""
	}
	arrs[1] = strings.Join(arrs[1:], "(")
	path := strings.TrimSpace(arrs[0])
	field := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(arrs[1]), ")"))
	subs := getSubstitions(envs)
	path, field = stringf.Format(path, subs), stringf.Format(field, subs)
	return path, field
}

func ParseKey(data interface{}, envs []string) (string, string, string) {
	switch data := data.(type) {
	case map[interface{}]interface{}:
		for k, v := range data {
			path, field := extractPathAndField(strings.TrimSpace(toString(k)), envs)
			val := strings.TrimSpace(toString(v))
			return path, field, val
		}
		return "", "", ""
	default:
		return "", "", toString(data)
	}
}

func ParseConfigMap(obj map[interface{}]interface{}, envs []string) []Config {
	configs := make([]Config, 0)
	for k, v := range obj {
		c := Config{}
		key := toString(k)
		if strings.Contains(key, "/") {
			c.Type = "file"
			c.Path = key
		} else {
			c.Type = "kv"
			c.Name = strings.TrimSpace(key)
		}

		c.VaultPath, c.VaultField, c.Value = ParseKey(v, envs)
		configs = append(configs, c)
	}
	return configs
}

func parse(configs []Config, vaultvalues []*string) (string, error) {
	if len(configs) != len(vaultvalues) {
		return "", fmt.Errorf("len configs and len vaultvalues not match, got %d, %d", len(configs), len(vaultvalues))
	}

	out := strings.Builder{}
	var err error

	for i, c := range configs {
		if vaultvalues[i]  != nil {
			c.Value = *vaultvalues[i]
		}

		if c.Type == "kv" {
			_, err = out.Write([]byte(ExportKv(c)))
		} else if c.Type == "file" {
			err = WriteFile(c)
		} else {
			err = fmt.Errorf("unknow type %s", c.Type)
		}
		if err != nil {
			break
		}
	}
	return out.String(), err
}
