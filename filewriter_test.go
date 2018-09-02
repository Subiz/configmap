package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func TestWriteFile(t *testing.T) {
	dir := "/tmp/a342341XXH/V2108"
	os.Remove(dir)
	value := `1
2\a\n
"3`
	c := Config{Path: dir + "/a", Value: value}
	cmd := WriteFile(c, "bash")

	_, err := exec.Command("/bin/bash", "-c", cmd).Output()
	if err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadFile(dir + "/a")
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != value {
		t.Fatalf("got %s", data)
	}
}

func TestToPrintf(t *testing.T) {
	tcs := []struct {
		content, expect string
	}{
		{"hello", `printf "%s" "hello"`},
		{"hello1\a1\nx", `printf "%s\n%s" "hello1\a1" "x"`},
		{"hello1\a1\n", `printf "%s\n%s" "hello1\a1" ""`},
		{"hell\"o", `printf "%s" "hell\"o"`},
	}

	for _, c := range tcs {
		out := toPrintf(c.content)
		if out != c.expect {
			t.Errorf("expect %s, got %s.", c.expect, out)
		}
	}
}
