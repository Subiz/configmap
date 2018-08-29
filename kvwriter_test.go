package main

import (
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
		name, value, expect string
	}{{"A", "B", `A="B"
`},
		{"C2f  ", "C ", `C2f="C "
`},
		{"C2f  ", "C \"", `C2f="C \""
`},
		{"44", `1
2
3
4`, `44="1
2
3
4"
`},
	}
	for _, c := range tcs {
		out := ExportKv(Config{Name: c.name, Value: c.value, Type: "kv"})
		if out != c.expect {
			t.Errorf("expect %s, got %s", c.expect, out)
		}
	}
}
