package main

import (
	"io/ioutil"
	"os"
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

func TestKvWriterPrepare(t *testing.T) {
	tcs := []struct {
		name, value, expect string
	}{{"A", "B", `A="B"
`},
		{"C2f  ", "C ", `C2f="C "
`},
		{"C2f  ", "C \"", `C2f="C \""
`},
	}
	for _, c := range tcs {
		out := prepareKv(Config{Name: c.name, Value: c.value, Type: "kv"})
		if out != c.expect {
			t.Errorf("expect %s, got %s", c.expect, out)
		}
	}
}

func TestKvWrite(t *testing.T) {
	configs := []Config{
		{Type: "kv", Name: "A", Value: "B"},
		{Type: "kv", Name: "C  ", Value: `'Thanh Van "`},
		{Type: "kv", Name: "44\n5", Value: "Mot\n Hai\n Ba"},
		{Type: "kv", Name: "", Value: ""},
		{Type: "kv", Name: "", Value: ""},
	}

	tmpfile, err := ioutil.TempFile("", "testkvwriter")
	if err != nil {
		t.Fatal(err)
	}
	filepath:= tmpfile.Name()
	tmpfile.Close()

	defer os.Remove(tmpfile.Name()) // clean up
	for _, c := range configs {
		c.Path = filepath
		err := WriteKv(c)
		if err != nil {
			t.Fatal(err)
		}
	}

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		t.Fatal(err)
	}

	expect := `A="B"
C="'Thanh Van \""
445="Mot
 Hai
 Ba"
`
	if string(data) != expect {
		t.Errorf("expect /%s/, got /%s/", expect, string(data))
	}
}
