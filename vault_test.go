package main

import "testing"

const g_testtoken = "ed5b9999-dfaf-d331-2010-481cb7398002"
const g_testpath = "secret/test_secret"
const g_addr = "https://vault.subiz.com"

func TestReadVault(t *testing.T) {
	data, err := ReadVault(g_addr, g_testtoken,
		[]string{g_testpath, g_testpath},
		[]string{"mot", "hai"})
	if err != nil {
		t.Fatal(err)
	}

	for i, d := range data {
		switch i {
		case 0:
			if string(d) != "hai" {
				t.Errorf("should be hai, got %s", d)
			}

		case 1:
			if string(d) != "11111" {
				t.Errorf("should be 11111, got %s", d)
			}
		default:
			t.Fatal("should not run this")
		}
	}
}
