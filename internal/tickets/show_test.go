package tickets_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/ro-56/taskr/internal/tickets"
)

// --- helpers ---

func setupTicket(t *testing.T, dir string, opts tickets.AddOptions) string {
	t.Helper()
	if err := tickets.Init(dir, "TKT"); err != nil && err != tickets.ErrAlreadyInitialized {
		t.Fatal(err)
	}
	id, err := tickets.Add(dir, opts)
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func appendBody(t *testing.T, path, body string) {
	t.Helper()
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if _, err := f.WriteString("\n" + body); err != nil {
		t.Fatal(err)
	}
}

// writeTicketWithDeps rewrites a ticket file's dependencies frontmatter field.
func writeTicketWithDeps(t *testing.T, dir, id string, deps []string) {
	t.Helper()
	path := filepath.Join(dir, ".tickets", id+".md")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	s := string(raw)
	start := strings.Index(s, "---\n")
	end := strings.Index(s[4:], "\n---\n")
	if start < 0 || end < 0 {
		t.Fatalf("no frontmatter in %s", path)
	}
	fmRaw := s[4 : end+4]
	var fm map[string]any
	if err := yaml.Unmarshal([]byte(fmRaw), &fm); err != nil {
		t.Fatal(err)
	}
	fm["dependencies"] = deps
	out, err := yaml.Marshal(fm)
	if err != nil {
		t.Fatal(err)
	}
	body := s[end+9:]
	content := fmt.Sprintf("---\n%s---\n%s", out, body)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

// --- tests ---

func TestShow_IncludesBody(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	id, _ := tickets.Add(dir, tickets.AddOptions{Title: "body test"})

	appendBody(t, fmt.Sprintf("%s/.tickets/%s.md", dir, id), "This is the ticket body.\n")

	result, err := tickets.Show(dir, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Body != "This is the ticket body.\n" {
		t.Errorf("body: got %q, want %q", result.Body, "This is the ticket body.\n")
	}
}

func TestShow_PartialID_ResolvesToFull(t *testing.T) {
	dir := t.TempDir()
	id := setupTicket(t, dir, tickets.AddOptions{Title: "partial match"})

	// use first 10 chars of the ID (e.g. "TKT-a3f8bc" from "TKT-a3f8bc2d")
	partial := id[:10]

	result, err := tickets.Show(dir, partial)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != id {
		t.Errorf("ID: got %q, want %q", result.ID, id)
	}
}

func TestShow_FindsArchivedTicket(t *testing.T) {
	dir := t.TempDir()
	id := setupTicket(t, dir, tickets.AddOptions{Title: "archived"})

	// move ticket to archive manually
	src := fmt.Sprintf("%s/.tickets/%s.md", dir, id)
	dst := fmt.Sprintf("%s/.tickets/archive/%s.md", dir, id)
	if err := os.Rename(src, dst); err != nil {
		t.Fatal(err)
	}

	result, err := tickets.Show(dir, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != id {
		t.Errorf("ID: got %q, want %q", result.ID, id)
	}
}

func TestShow_AmbiguousPartialID_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")

	// create two tickets with the same prefix up to a short partial
	// We can't control IDs, but we can use "TKT-" as the partial — all tickets match
	tickets.Add(dir, tickets.AddOptions{Title: "first"})
	tickets.Add(dir, tickets.AddOptions{Title: "second"})

	_, err := tickets.Show(dir, "TKT-")
	if err == nil {
		t.Fatal("expected error for ambiguous partial ID, got nil")
	}

	var ambErr *tickets.ErrAmbiguous
	if !errors.As(err, &ambErr) {
		t.Fatalf("expected ErrAmbiguous, got %T: %v", err, err)
	}
	if len(ambErr.Matches) < 2 {
		t.Errorf("expected at least 2 matches, got %d: %v", len(ambErr.Matches), ambErr.Matches)
	}
}

func TestShow_UnknownID_ReturnsNotFound(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")

	_, err := tickets.Show(dir, "TKT-00000000")
	if !errors.Is(err, tickets.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestShow_DependsOn_Empty(t *testing.T) {
	dir := t.TempDir()
	id := setupTicket(t, dir, tickets.AddOptions{Title: "no deps"})

	result, err := tickets.Show(dir, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.DependsOn) != 0 {
		t.Errorf("DependsOn: expected empty, got %v", result.DependsOn)
	}
	if len(result.RequiredBy) != 0 {
		t.Errorf("RequiredBy: expected empty, got %v", result.RequiredBy)
	}
}

func TestShow_DependsOn_WithDeps(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")

	depID, _ := tickets.Add(dir, tickets.AddOptions{Title: "dependency"})
	mainID, _ := tickets.Add(dir, tickets.AddOptions{Title: "main ticket"})

	// write mainID's frontmatter with depID in dependencies
	writeTicketWithDeps(t, dir, mainID, []string{depID})

	result, err := tickets.Show(dir, mainID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.DependsOn) != 1 {
		t.Fatalf("DependsOn: expected 1 entry, got %d: %v", len(result.DependsOn), result.DependsOn)
	}
	if result.DependsOn[0].ID != depID {
		t.Errorf("DependsOn[0].ID: got %q, want %q", result.DependsOn[0].ID, depID)
	}
	if result.DependsOn[0].Status != "open" {
		t.Errorf("DependsOn[0].Status: got %q, want %q", result.DependsOn[0].Status, "open")
	}
}

func TestShow_RequiredBy(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")

	baseID, _ := tickets.Add(dir, tickets.AddOptions{Title: "base ticket"})
	dependentID, _ := tickets.Add(dir, tickets.AddOptions{Title: "dependent ticket"})

	// dependentID lists baseID as a dependency
	writeTicketWithDeps(t, dir, dependentID, []string{baseID})

	result, err := tickets.Show(dir, baseID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.RequiredBy) != 1 {
		t.Fatalf("RequiredBy: expected 1, got %d: %v", len(result.RequiredBy), result.RequiredBy)
	}
	if result.RequiredBy[0].ID != dependentID {
		t.Errorf("RequiredBy[0].ID: got %q, want %q", result.RequiredBy[0].ID, dependentID)
	}
	if result.RequiredBy[0].Status != "open" {
		t.Errorf("RequiredBy[0].Status: got %q, want %q", result.RequiredBy[0].Status, "open")
	}
}

func TestShow_FullID_ReturnsFrontmatter(t *testing.T) {
	dir := t.TempDir()
	id := setupTicket(t, dir, tickets.AddOptions{Title: "Fix login bug", Type: "bug"})

	result, err := tickets.Show(dir, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != id {
		t.Errorf("ID: got %q, want %q", result.ID, id)
	}
	if got := fmt.Sprint(result.Frontmatter["title"]); got != "Fix login bug" {
		t.Errorf("title: got %q, want %q", got, "Fix login bug")
	}
	if got := fmt.Sprint(result.Frontmatter["type"]); got != "bug" {
		t.Errorf("type: got %q, want %q", got, "bug")
	}
	if got := fmt.Sprint(result.Frontmatter["status"]); got != "open" {
		t.Errorf("status: got %q, want %q", got, "open")
	}
}
