package main

import (
	"sort"
	"bytes"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"strings"
	"testing"
)

func TestExtractPathAndField(t *testing.T) {
	tcs := []struct {
		envs            []string
		in, path, field string
	}{{nil, "secret/stripe(api_dev)", "secret/stripe", "api_dev"},
		{nil, "secret/stripe(api_dev", "secret/stripe", "api_dev"},
		{nil, "secret/stripe(api_dev  )", "secret/stripe", "api_dev"},
		{nil, "  secret/stripe(api_dev  )", "secret/stripe", "api_dev"},
		{nil, "  secret/stripe((api_dev  )  ", "secret/stripe", "(api_dev"},
		{[]string{"_env=dev", "_key=4"}, "secret/{env}(2{key})", "secret/dev", "24"},
	}

	for _, c := range tcs {
		epath, efield := extractPathAndField(c.in, c.envs)
		if epath != c.path || efield != c.field {
			t.Errorf("should be equal, got %s, %s, %s", c.in, epath, efield)
		}
	}
}

func TestParseKey(t *testing.T) {
	tcs := []struct {
		obj                interface{}
		path, field, value string
	}{{map[interface{}]interface{}{"secret/stripe({dev}_apikey) ": "2222"}, "secret/stripe", "{dev}_apikey", "2222"},
		{map[interface{}]interface{}{"stripe ke": "123"}, "stripe ke", "", "123"},
		{"default value ", "", "", "default value "},
	}

	for _, c := range tcs {
		epath, efield, evalue := ParseKey(c.obj, nil)
		if epath != c.path || efield != c.field || evalue != c.value {
			t.Fatalf("wrong %v, %s, %s, %s.", c.obj, epath, efield, evalue)
		}
	}
}

func TestParseObject(t *testing.T) {
	var data = `
---
# 0
stripe_apikey:
  secret/stripe({dev}_apikey): "222222222223333333333333"

# 1 - default vaule
/workspace/x:
  "stripe ke": asdlkfjkalsjdfkljasdklfj

# 2
s3_apikey: default value
---
# 3
version: 1
---
# 4
a:
  b(c): 4

`
	obj := make(map[interface{}]interface{})
	dec := yaml.NewDecoder(strings.NewReader(data))
	for dec.Decode(&obj) == nil {
	}

	//err := yaml.Unmarshal([]byte(data), &x)

	configs := ParseConfigMap(obj, nil)
	expects := []Config{{
		Name:       "stripe_apikey",
		Path:       "",
		Type:       "kv",
		Value:      "222222222223333333333333",
		VaultPath:  "secret/stripe",
		VaultField: "{dev}_apikey",
	}, {
		Name:       "",
		Path:       "/workspace/x",
		Type:       "file",
		Value:      "asdlkfjkalsjdfkljasdklfj",
		VaultPath:  "stripe ke",
		VaultField: "",
	}, {
		Name:  "s3_apikey",
		Type:  "kv",
		Value: "default value",
	}, {
		Name:  "version",
		Type:  "kv",
		Value: "1",
	}, {
		Name:       "a",
		Type:       "kv",
		Value:      "4",
		VaultPath:  "b",
		VaultField: "c",
	}}

	if len(configs) != len(expects) {
		t.Fatalf("expect %d keys, got %d keys", len(expects), len(configs))
	}

	if !compareConfigArr(expects, configs) {
		t.Fatalf("should be equal, expect %v, got %v", jsonify(expects), jsonify(configs))
	}

}

func jsonify(a interface{}) string {
	ab, _ := json.Marshal(&a)
	return string(ab)
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
		t.Errorf("Expect %v, Got %v", expect, out)
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
