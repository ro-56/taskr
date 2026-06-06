package tickets_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/ro-56/taskr/internal/tickets"
)

func TestAdd_ReturnsValidID(t *testing.T) {
	dir := t.TempDir()
	if err := tickets.Init(dir, "TKT"); err != nil {
		t.Fatal(err)
	}

	id, err := tickets.Add(dir, tickets.AddOptions{Title: "Fix login bug"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// must match TKT-xxxxxxxx (8 lowercase hex chars)
	matched, _ := regexp.MatchString(`^TKT-[0-9a-f]{8}$`, id)
	if !matched {
		t.Errorf("id %q does not match PREFIX-xxxxxxxx format", id)
	}
}

func TestAdd_CreatesTicketFile(t *testing.T) {
	dir := t.TempDir()
	if err := tickets.Init(dir, "TKT"); err != nil {
		t.Fatal(err)
	}

	id, err := tickets.Add(dir, tickets.AddOptions{Title: "Fix login bug"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	path := filepath.Join(dir, ".tickets", id+".md")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected ticket file at %s", path)
	}
}

func TestAdd_FrontmatterFields(t *testing.T) {
	dir := t.TempDir()
	if err := tickets.Init(dir, "TKT"); err != nil {
		t.Fatal(err)
	}

	p := 1
	id, err := tickets.Add(dir, tickets.AddOptions{
		Title:    "Fix login bug",
		Type:     "bug",
		Priority: &p,
		Mode:     "afk",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fm := readFrontmatter(t, filepath.Join(dir, ".tickets", id+".md"))

	checks := map[string]string{
		"id":     id,
		"title":  "Fix login bug",
		"status": "open",
		"type":   "bug",
		"mode":   "afk",
	}
	for field, want := range checks {
		if got := fmt.Sprint(fm[field]); got != want {
			t.Errorf("field %q: got %q, want %q", field, got, want)
		}
	}
	if got := fmt.Sprint(fm["priority"]); got != "1" {
		t.Errorf("priority: got %q, want %q", got, "1")
	}
}

func TestAdd_Defaults(t *testing.T) {
	dir := t.TempDir()
	if err := tickets.Init(dir, "TKT"); err != nil {
		t.Fatal(err)
	}

	id, err := tickets.Add(dir, tickets.AddOptions{Title: "Default ticket"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fm := readFrontmatter(t, filepath.Join(dir, ".tickets", id+".md"))

	if got := fmt.Sprint(fm["type"]); got != "task" {
		t.Errorf("default type: got %q, want %q", got, "task")
	}
	if got := fmt.Sprint(fm["priority"]); got != "2" {
		t.Errorf("default priority: got %q, want %q", got, "2")
	}
	if got := fmt.Sprint(fm["mode"]); got != "hitl" {
		t.Errorf("default mode: got %q, want %q", got, "hitl")
	}
}

func TestAdd_Timestamps(t *testing.T) {
	dir := t.TempDir()
	if err := tickets.Init(dir, "TKT"); err != nil {
		t.Fatal(err)
	}

	before := time.Now().UTC().Truncate(time.Second)
	id, err := tickets.Add(dir, tickets.AddOptions{Title: "ts test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	after := time.Now().UTC().Add(time.Second)

	fm := readFrontmatter(t, filepath.Join(dir, ".tickets", id+".md"))
	for _, field := range []string{"created", "updated"} {
		raw, ok := fm[field].(string)
		if !ok || raw == "" {
			t.Errorf("field %q missing or empty", field)
			continue
		}
		ts, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			t.Errorf("field %q not a valid RFC3339 timestamp: %v", field, err)
			continue
		}
		if ts.Before(before) || ts.After(after) {
			t.Errorf("field %q timestamp %v out of range [%v, %v]", field, ts, before, after)
		}
	}
}

func TestAdd_EmptyDepsAndLinks(t *testing.T) {
	dir := t.TempDir()
	if err := tickets.Init(dir, "TKT"); err != nil {
		t.Fatal(err)
	}

	id, err := tickets.Add(dir, tickets.AddOptions{Title: "deps test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fm := readFrontmatter(t, filepath.Join(dir, ".tickets", id+".md"))
	for _, field := range []string{"dependencies", "links"} {
		val, ok := fm[field]
		if !ok {
			t.Errorf("field %q missing", field)
			continue
		}
		slice, ok := val.([]any)
		if !ok || len(slice) != 0 {
			t.Errorf("field %q: want empty list, got %v", field, val)
		}
	}
}

func TestAdd_Tags(t *testing.T) {
	dir := t.TempDir()
	if err := tickets.Init(dir, "TKT"); err != nil {
		t.Fatal(err)
	}

	id, err := tickets.Add(dir, tickets.AddOptions{
		Title: "tagged ticket",
		Tags:  []string{"auth", "backend"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fm := readFrontmatter(t, filepath.Join(dir, ".tickets", id+".md"))
	tags, ok := fm["tags"].([]any)
	if !ok {
		t.Fatalf("tags field is not a list: %v", fm["tags"])
	}
	if len(tags) != 2 || fmt.Sprint(tags[0]) != "auth" || fmt.Sprint(tags[1]) != "backend" {
		t.Errorf("tags: got %v, want [auth backend]", tags)
	}
}

func TestAdd_PriorityCritical(t *testing.T) {
	dir := t.TempDir()
	if err := tickets.Init(dir, "TKT"); err != nil {
		t.Fatal(err)
	}

	p := 0 // critical
	id, err := tickets.Add(dir, tickets.AddOptions{Title: "critical ticket", Priority: &p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fm := readFrontmatter(t, filepath.Join(dir, ".tickets", id+".md"))
	if got := fmt.Sprint(fm["priority"]); got != "0" {
		t.Errorf("priority: got %q, want %q", got, "0")
	}
}

func TestAdd_NotInitialized(t *testing.T) {
	dir := t.TempDir()

	_, err := tickets.Add(dir, tickets.AddOptions{Title: "Fix login bug"})
	if err != tickets.ErrNotInitialized {
		t.Errorf("expected ErrNotInitialized, got %v", err)
	}
}

// readFrontmatter parses YAML frontmatter from a markdown file.
func readFrontmatter(t *testing.T, path string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read %s: %v", path, err)
	}
	content := string(raw)
	if !strings.HasPrefix(content, "---\n") {
		t.Fatalf("file does not start with frontmatter: %s", path)
	}
	end := strings.Index(content[4:], "\n---\n")
	if end < 0 {
		t.Fatalf("frontmatter closing delimiter not found in %s", path)
	}
	fmRaw := content[4 : end+4]
	var out map[string]any
	if err := yaml.Unmarshal([]byte(fmRaw), &out); err != nil {
		t.Fatalf("yaml parse error in %s: %v", path, err)
	}
	return out
}
