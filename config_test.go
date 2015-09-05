package pab

import "testing"

func TestConfigDefaultValues(t *testing.T) {
	c := NewConfig()
	if c.Name != "pab" {
		t.Fatal("name is not pab")
	}

	if c.Adapter != "shell" {
		t.Fatal("adapter name is not shell")
	}

	if c.Storage != "memory" {
		t.Fatal("storage is not memory")
	}
	if c.HTTPAddr != ":8001" {
		t.Fatal("http addr is not :8001")
	}
}

func TestEnv(t *testing.T) {
	c := NewConfig()

	if len(c.Env) != 0 {
		t.Fail()
	}
	c.Env["add_new"] = "test"
}
