package main

import (
	"testing"
	"os"
	"io/ioutil"
)

func TestWriteFile(t *testing.T) {
	dir := "/tmp/a342341XXH/V2108"
	os.Remove(dir)
	c := Config{Path: dir + "/a", Value: "123"}
	if err := WriteFile(c); err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadFile(dir + "/a")
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "123" {
		t.Fatalf("should be 123, got %s", data)
	}

	// write to exist file should cause error
	c = Config{Path: dir, Value: "321"}
	err = WriteFile(c)
	if err == nil {
		t.Fatal("should be err, got nil")
	}
}
