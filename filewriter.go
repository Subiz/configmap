package main

import (
	"io/ioutil"
	"os"
	"path"
)

func WriteFile(c Config) error {
	dir := path.Dir(c.Path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		println("here")
		return err
	}
	return ioutil.WriteFile(c.Path, []byte(c.Value), 0700)
}
