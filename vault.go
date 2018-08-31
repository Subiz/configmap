package main

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"sync"
	"strings"
)

// Read kv in vault, return nil if not found
func readVault(addr, token string, paths, fields []string) ([]*string, error) {
	client, err := api.NewClient(&api.Config{Address: addr})
	if err != nil {
		return nil, err
	}
	client.SetToken(token)

	if len(paths) != len(fields) {
		return nil, fmt.Errorf("paths and fields should have same len, got %d, %d", len(paths), len(fields))
	}

	data := make([]*string, len(paths))
	errs := make([]error, len(paths))
	wg := &sync.WaitGroup{}
	wg.Add(len(paths))
	for i := range paths {
		go func(i int) {
			defer wg.Done()
			path, field := strings.TrimSpace(paths[i]), strings.TrimSpace(fields[i])
			if path == "" && field == "" {
				return
			}
			secretValues, err := client.Logical().Read(path)
			if err != nil {
				errs[i] = err
				return
			}
			if secretValues == nil {
				return
			}
			d := secretValues.Data[field]
			if s, b := d.(string); !b {
				errs[i] = fmt.Errorf("secret is not string, got %v", d)
			} else {
				data[i] = &s
			}
		}(i)
	}

	wg.Wait()
	var errstr string
	for i, err := range errs {
		if err != nil {
			errstr = errstr + fmt.Sprintf("Cannot read %s(%s): %s\n", paths[i], fields[i], err.Error())
		}
	}
	if errstr != "" {
		err = fmt.Errorf("%s", errstr)
	}
	return data, err
}
