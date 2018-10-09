package main

import (
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

// Read kv in file, return nil if not found
func readFile(filepath string, paths, keys []string) ([]*string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m := make(map[interface{}]interface{})
	dec := yaml.NewDecoder(f)
	for {
		err := dec.Decode(&m)
		if err == nil {
			continue
		}

		if err == io.EOF {
			break
		}

		return nil, err
	}

	data := make([]*string, len(keys))
	for i, k := range keys {
		path, _ := m[paths[i]].(map[interface{}]interface{})
		data[i] = toPString(path[k])
	}
	return data, err
}
