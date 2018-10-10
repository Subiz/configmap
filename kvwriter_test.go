package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestToBashName(t *testing.T) {
	tcs := []struct {
		name, expect string
	}{{"A", "A"},
		{" ", ""},
		{"界 aa Y 界", "aaY"},
		{"a + $ 3_1-1", "a3_11"},
		{"44", "44"},
	}

	for _, c := range tcs {
		out := toBashName(c.name)
		if out != c.expect {
			t.Errorf("expect %s, got %s", c.expect, out)
		}
	}
}

func TestKvWriterExport(t *testing.T) {
	tcs := []struct {
		name, value string
	}{{"A", "B"},
		{"C2f  ", "C "},
		{"C2f  ", "C \""},
	}

	for _, c := range tcs {
		command := exportKv(Config{Name: c.name, Value: c.value, Type: "kv"})
		op, err := exec.Command("/bin/bash", "-c", "export "+command+"\n"+`printf "%s" "$`+strings.TrimSpace(c.name)+`"`).Output()
		if err != nil {
			t.Error(err)
		}

		if string(op) != c.value {
			t.Errorf("should be %s., got %v.", c.value, op)
		}
	}
}
