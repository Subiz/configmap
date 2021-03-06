package main

import (
	"fmt"
	"testing"
)

func ps(s string) *string { return &s }

func TestFile(t *testing.T) {
	tcs := []struct {
		desc     string
		filepath string
		paths    []string
		keys     []string
		expect   []*string
		err      error
	}{{
		"normal",
		"./test/config1.yaml",
		[]string{"stripe", "stripe", "gcp"},
		[]string{"apikey", "account", "file"},
		[]*string{ps("1"), ps("2"), ps("3")},
		nil,
	}, {
		"normal2",
		"./test/config1.yaml",
		[]string{"stripe", "stripe", "gcp"},
		[]string{"apikey", "account", "file2"},
		[]*string{ps("1"), ps("2"), nil},
		nil,
	}, {
		"file not found",
		"./test/confi1.yaml",
		[]string{"stripe"},
		[]string{"apikey"},
		[]*string{},
		fmt.Errorf("open ./test/confi1.yaml: no such file or directory"),
	}, {
		"invalid yaml",
		"./test/config_invalid.yaml",
		[]string{"stripe"},
		[]string{"apikey"},
		[]*string{},
		fmt.Errorf("yaml: line 1: did not find expected node content"),
	}, {
		"number",
		"./test/config1.yaml",
		[]string{"s2"},
		[]string{"password"},
		[]*string{ps("2")},
		nil,
	}}

	for _, tc := range tcs {
		out, err := readFile(tc.filepath, tc.paths, tc.keys)
		if !compareErr(err, tc.err) {
			t.Errorf("[%s] expect err {%v} got {%v}", tc.desc, tc.err, err)
		}
		if !comparePString(out, tc.expect) {
			t.Errorf("[%s] expect %v, got %v", tc.desc, tc.expect, out)
		}
	}
}

func compareErr(e1, e2 error) bool {
	if e1 == e2 {
		return true
	}

	if e1 == nil || e2 == nil {
		return false
	}

	return e1.Error() == e2.Error()
}

func comparePString(a, b []*string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] == nil {
			if b[i] != nil {
				return false
			}
			continue
		}

		if b[i] == nil {
			return false
		}

		if *a[i] != *b[i] {
			return false
		}
	}
	return true
}
