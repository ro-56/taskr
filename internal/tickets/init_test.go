package tickets_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/rodrigo/taskr/internal/tickets"
)

func TestInit_CustomPrefix(t *testing.T) {
	dir := t.TempDir()

	if err := tickets.Init(dir, "FOO"); err != nil {
		t.Fatalf("Init returned error: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(dir, ".tickets", "config.json"))
	var cfg map[string]string
	json.Unmarshal(data, &cfg)

	if cfg["prefix"] != "FOO" {
		t.Errorf("expected prefix FOO, got %q", cfg["prefix"])
	}
}

func TestInit_AlreadyInitialized(t *testing.T) {
	dir := t.TempDir()

	if err := tickets.Init(dir, "TKT"); err != nil {
		t.Fatalf("first Init failed: %v", err)
	}

	err := tickets.Init(dir, "TKT")
	if err != tickets.ErrAlreadyInitialized {
		t.Errorf("expected ErrAlreadyInitialized, got %v", err)
	}
}

func TestDerivePrefix(t *testing.T) {
	cases := []struct{ in, want string }{
		{"my-project", "MYPROJECT"},
		{"hello world 123", "HELLOWORLD123"},
		{"taskr", "TASKR"},
		{"My_App.v2", "MYAPPV2"},
	}
	for _, c := range cases {
		got := tickets.DerivePrefix(c.in)
		if got != c.want {
			t.Errorf("DerivePrefix(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestInit_CreatesDirectoriesAndConfig(t *testing.T) {
	dir := t.TempDir()

	if err := tickets.Init(dir, ""); err != nil {
		t.Fatalf("Init returned error: %v", err)
	}

	for _, sub := range []string{".tickets", filepath.Join(".tickets", "archive")} {
		if _, err := os.Stat(filepath.Join(dir, sub)); os.IsNotExist(err) {
			t.Errorf("expected directory %s to exist", sub)
		}
	}

	data, err := os.ReadFile(filepath.Join(dir, ".tickets", "config.json"))
	if err != nil {
		t.Fatalf("config.json missing: %v", err)
	}

	var cfg map[string]string
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("config.json is not valid JSON: %v", err)
	}

	if cfg["prefix"] == "" {
		t.Error("expected prefix field in config.json")
	}
}
