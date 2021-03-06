package main

import (
	"bytes"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"sort"
	"strings"
	"testing"
)

func TestExtractPathAndField(t *testing.T) {
	tcs := []struct {
		envs            []string
		in, path, field string
	}{{nil, "secret/stripe.api_dev", "secret/stripe", "api_dev"},
		{nil, "secret/stripe.api_dev", "secret/stripe", "api_dev"},
		{nil, "secret/stripe.api_dev .", "secret/stripe", "api_dev"},
		{nil, "  secret/stripe.api_dev  .", "secret/stripe", "api_dev"},
		{nil, "  secret/stripe.api_dev    ", "secret/stripe", "api_dev"},
		{[]string{"_env=dev", "_key=4"}, "secret/{env}.2{key}", "secret/dev", "24"},
	}

	for _, c := range tcs {
		epath, efield := extractPathAndField(c.in, c.envs)
		if epath != c.path || efield != c.field {
			t.Errorf("should be equal, got %s, %s, %s", c.in, epath, efield)
		}
	}
}

func TestParseObject(t *testing.T) {
	var data = `
---
# 0
stripe_apikey: secret/stripe.{dev}_apikey

# 1 - default vaule
/workspace/x: stripe.ke

# 2
s3_apikey:
---
# 3
version: 1
---
# 4
a: b.c

`
	obj := make(map[interface{}]interface{})
	dec := yaml.NewDecoder(strings.NewReader(data))
	for dec.Decode(&obj) == nil {
	}

	configs := ParseConfigMap(obj, nil)
	expects := []Config{{
		Name:        "stripe_apikey",
		Path:        "",
		Type:        "kv",
		Value:       "",
		ConfigPath:  "secret/stripe",
		ConfigField: "{dev}_apikey",
	}, {
		Name:        "",
		Path:        "/workspace/x",
		Type:        "file",
		Value:       "",
		ConfigPath:  "stripe",
		ConfigField: "ke",
	}, {
		Name:  "s3_apikey",
		Type:  "kv",
		Value: "",
	}, {
		Name:  "version",
		Type:  "kv",
		Value: "",
		ConfigPath: "1",
	}, {
		Name:        "a",
		Type:        "kv",
		Value:       "",
		ConfigPath:  "b",
		ConfigField: "c",
	}}

	if len(configs) != len(expects) {
		t.Fatalf("expect %d keys, got %d keys", len(expects), len(configs))
	}

	if !compareConfigArr(expects, configs) {
		t.Fatalf("should be equal, expect %v, got %v", jsonify(expects), jsonify(configs))
	}
}

func compareConfigArr(a, b []Config) bool {
	sort.Sort(ByConfigNameAndPath(a))
	sort.Sort(ByConfigNameAndPath(b))

	for i := range a {
		ai, bi := a[i], b[i]
		if !compareConfig(ai, bi) {
			return false
		}
	}
	return true
}

func compareConfig(a, b Config) bool {
	ab, _ := json.Marshal(&a)
	bb, _ := json.Marshal(&b)
	return bytes.Compare(ab, bb) == 0
}

func TestGetSubstitution(t *testing.T) {
	envs := []string{
		`A=4`,
		`B=5`,
		`_HV=haivan`,
		`_H_V====`,
	}

	expect := map[string]string{
		"hv":  "haivan",
		"h_v": "===",
	}
	out := getSubstitions(envs)
	if !compareMap(expect, out) {
		t.Errorf("expect %v, got %v", expect, out)
	}
}

func compareMap(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		vb, f := b[k]
		if !f {
			return false
		}
		if v != vb {
			return false
		}
	}
	return true
}

func jsonify(a interface{}) string {
	ab, _ := json.Marshal(a)
	return string(ab)
}
