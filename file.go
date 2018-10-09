package main

import (
	"encoding/json"
	"io/ioutil"
)

// Read kv in file, return nil if not found
func readFile(filepath string, paths, keys []string) ([]*string, error) {
	dat, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	if err := json.Unmarshal(dat, &m); err != nil {
		return nil, err
	}

	data := make([]*string, len(keys))
	for i, k := range keys {
		path, _ := m[paths[i]].(map[string]interface{})
		data[i] = toPString(path[k])
	}
	return data, err
}
