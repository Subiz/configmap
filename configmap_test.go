package main

import (
	"encoding/json"
//	"fmt"
	"testing"
	"gopkg.in/yaml.v2"
	"bytes"
)

func TestExtractPathAndField(t *testing.T) {
	tcs := []struct{
		in, path, field string
	}{{"secret/stripe(api_dev)", "secret/stripe", "api_dev"},
		{"secret/stripe(api_dev", "secret/stripe", "api_dev"},
		{"secret/stripe(api_dev  )", "secret/stripe", "api_dev"},
		{"  secret/stripe(api_dev  )", "secret/stripe", "api_dev"},
		{"  secret/stripe((api_dev  )  ", "secret/stripe", "(api_dev"},
	}

	for _, c := range tcs {
		epath, efield := extractPathAndField(c.in, nil)
		if epath != c.path || efield != c.field {
			t.Errorf("should be equal, got %s, %s, %s", c.in, epath, efield)
		}
	}
}

func TestParseKey(t *testing.T) {
	tcs := []struct {
		obj interface{}
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
# 0
stripe_apikey:
  secret/stripe({dev}_apikey): "222222222223333333333333"

# 1 - default vaule
/workspace/x:
  "stripe ke": asdlkfjkalsjdfkljasdklfj

# 2
s3_apikey: default value
`
	obj := make(map[interface{}]interface{})
	err := yaml.Unmarshal([]byte(data), &obj)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	configs := ParseConfigMap(obj, nil)
	expects := []Config{{
		Name: "stripe_apikey",
		Path: "",
		Type: "kv",
		Value: "222222222223333333333333",
		VaultPath: "secret/stripe",
		VaultField: "{dev}_apikey",
	}, {
		Name: "",
		Path: "/workspace/x",
		Type: "file",
		Value: "asdlkfjkalsjdfkljasdklfj",
		VaultPath: "stripe ke",
		VaultField: "",
	}, {
		Name: "s3_apikey",
		Path: "",
		Type: "kv",
		Value: "default value",
		VaultPath: "",
		VaultField: "",
	}}

	if len(configs) != len(expects) {
		t.Fatalf("expect %d keys, got %d keys", len(expects), len(configs))
	}

	for i := range configs {
		c, e := configs[i], expects[i]
		if !compareConfig(c, e) {
			t.Fatalf("should be equal, expect %v, got %v", jsonify(e), jsonify(c))
		}
	}
}

func jsonify(a Config) string {
	ab, _ := json.Marshal(&a)
	return string(ab)
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
		"hv": "haivan",
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
