package tickets_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ro-56/taskr/internal/tickets"
)

// writeRawTicket writes a ticket file directly with the given YAML frontmatter fields.
func writeRawTicket(t *testing.T, dir, id string, inArchive bool, fields map[string]string) {
	t.Helper()
	var sb strings.Builder
	sb.WriteString("---\n")
	for k, v := range fields {
		fmt.Fprintf(&sb, "%s: %s\n", k, v)
	}
	sb.WriteString("---\n")
	subdir := filepath.Join(dir, ".tickets")
	if inArchive {
		subdir = filepath.Join(subdir, "archive")
	}
	if err := os.WriteFile(filepath.Join(subdir, id+".md"), []byte(sb.String()), 0644); err != nil {
		t.Fatal(err)
	}
}

func containsViolation(violations []tickets.Violation, id, substr string) bool {
	for _, v := range violations {
		if v.TicketID == id && strings.Contains(v.Message, substr) {
			return true
		}
	}
	return false
}

// validFields returns a map of all required fields with valid values for a given ID.
func validFields(id string) map[string]string {
	return map[string]string{
		"id":       id,
		"title":    "Test ticket",
		"status":   "open",
		"type":     "task",
		"priority": "2",
		"mode":     "hitl",
		"created":  "2026-01-01T00:00:00Z",
		"updated":  "2026-01-01T00:00:00Z",
	}
}

// --- Tracer bullet ---

func TestCheck_CleanRepo_NoViolations(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	tickets.Add(dir, tickets.AddOptions{Title: "A"})

	violations, err := tickets.Check(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Errorf("expected no violations, got: %v", violations)
	}
}

// --- Field value violations ---

func TestCheck_InvalidStatus(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	id := "TKT-aabbccdd"
	f := validFields(id)
	f["status"] = "wip"
	writeRawTicket(t, dir, id, false, f)

	violations, err := tickets.Check(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsViolation(violations, id, "status") {
		t.Errorf("expected status violation, got: %v", violations)
	}
}

func TestCheck_InvalidType(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	id := "TKT-aabbccdd"
	f := validFields(id)
	f["type"] = "story"
	writeRawTicket(t, dir, id, false, f)

	violations, err := tickets.Check(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsViolation(violations, id, "type") {
		t.Errorf("expected type violation, got: %v", violations)
	}
}

func TestCheck_InvalidPriority(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	id := "TKT-aabbccdd"
	f := validFields(id)
	f["priority"] = "5"
	writeRawTicket(t, dir, id, false, f)

	violations, err := tickets.Check(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsViolation(violations, id, "priority") {
		t.Errorf("expected priority violation, got: %v", violations)
	}
}

func TestCheck_InvalidMode(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	id := "TKT-aabbccdd"
	f := validFields(id)
	f["mode"] = "auto"
	writeRawTicket(t, dir, id, false, f)

	violations, err := tickets.Check(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsViolation(violations, id, "mode") {
		t.Errorf("expected mode violation, got: %v", violations)
	}
}

// --- Required fields ---

func TestCheck_MissingRequiredField(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	id := "TKT-aabbccdd"
	f := validFields(id)
	delete(f, "title")
	writeRawTicket(t, dir, id, false, f)

	violations, err := tickets.Check(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsViolation(violations, id, "title") {
		t.Errorf("expected missing title violation, got: %v", violations)
	}
}

// --- ID integrity ---

func TestCheck_IDFilenameMismatch(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	fileID := "TKT-aabbccdd"
	f := validFields(fileID)
	f["id"] = "TKT-11223344" // different from filename
	writeRawTicket(t, dir, fileID, false, f)

	violations, err := tickets.Check(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsViolation(violations, "TKT-11223344", "does not match filename") {
		t.Errorf("expected ID/filename mismatch violation, got: %v", violations)
	}
}

func TestCheck_IDFormatMismatch(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	// Use a different prefix or wrong format
	id := "TKT-aabbccdd"
	f := validFields(id)
	f["id"] = "WRONG-aabbccdd"
	// Write file with WRONG- prefix so id matches filename but wrong prefix
	wrongID := "WRONG-aabbccdd"
	writeRawTicket(t, dir, wrongID, false, f)
	// Note: we need the file to be named WRONG-aabbccdd.md but id also WRONG-aabbccdd
	// Re-write with matching id and filename to isolate the prefix check
	f2 := validFields(wrongID)
	writeRawTicket(t, dir, wrongID, false, f2)

	violations, err := tickets.Check(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsViolation(violations, wrongID, "does not match expected format") {
		t.Errorf("expected ID format violation, got: %v", violations)
	}
}

// --- File location ---

func TestCheck_ClosedTicketInActiveDir(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	id := "TKT-aabbccdd"
	f := validFields(id)
	f["status"] = "closed"
	writeRawTicket(t, dir, id, false, f) // closed but not in archive

	violations, err := tickets.Check(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsViolation(violations, id, "archive") {
		t.Errorf("expected location violation for closed ticket in active dir, got: %v", violations)
	}
}

func TestCheck_OpenTicketInArchiveDir(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	id := "TKT-aabbccdd"
	f := validFields(id)
	f["status"] = "open"
	writeRawTicket(t, dir, id, true, f) // open but in archive

	violations, err := tickets.Check(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsViolation(violations, id, "archive") {
		t.Errorf("expected location violation for open ticket in archive dir, got: %v", violations)
	}
}

// --- Dependencies ---

func TestCheck_MissingDependency(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	id := "TKT-aabbccdd"
	f := validFields(id)
	f["dependencies"] = "[TKT-deadbeef]"
	writeRawTicket(t, dir, id, false, f)

	violations, err := tickets.Check(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsViolation(violations, id, "TKT-deadbeef") {
		t.Errorf("expected missing dependency violation, got: %v", violations)
	}
}

func TestCheck_DependencyCycle(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID, _ := tickets.Add(dir, tickets.AddOptions{Title: "A"})
	bID, _ := tickets.Add(dir, tickets.AddOptions{Title: "B"})
	cID, _ := tickets.Add(dir, tickets.AddOptions{Title: "C"})

	tickets.Link(dir, aID, bID) // A → B
	tickets.Link(dir, bID, cID) // B → C
	// Manually inject C → A to create a cycle without going through Link's cycle check
	result, _ := tickets.Show(dir, cID)
	_ = result
	// Use Link but intercept: we need to write the cycle directly
	// Write C's file with A in dependencies bypassing Link's cycle guard
	cPath := filepath.Join(dir, ".tickets", cID+".md")
	raw, _ := os.ReadFile(cPath)
	content := strings.Replace(string(raw), "dependencies: []\n", fmt.Sprintf("dependencies:\n  - %s\n", aID), 1)
	os.WriteFile(cPath, []byte(content), 0644)

	violations, err := tickets.Check(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	hasCycleViolation := false
	for _, v := range violations {
		if strings.Contains(v.Message, "cycle") {
			hasCycleViolation = true
			break
		}
	}
	if !hasCycleViolation {
		t.Errorf("expected cycle violation, got: %v", violations)
	}
}

// --- Links ---

func TestCheck_MissingLink(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	id := "TKT-aabbccdd"
	f := validFields(id)
	f["links"] = "[TKT-deadbeef]"
	writeRawTicket(t, dir, id, false, f)

	violations, err := tickets.Check(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsViolation(violations, id, "TKT-deadbeef") {
		t.Errorf("expected missing link violation, got: %v", violations)
	}
}

func TestCheck_AsymmetricLink(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID := "TKT-aabbccdd"
	bID := "TKT-11223344"
	// A links B, but B does not link A
	fa := validFields(aID)
	fa["links"] = fmt.Sprintf("[%s]", bID)
	writeRawTicket(t, dir, aID, false, fa)
	writeRawTicket(t, dir, bID, false, validFields(bID))

	violations, err := tickets.Check(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsViolation(violations, aID, "not symmetric") {
		t.Errorf("expected asymmetric link violation, got: %v", violations)
	}
}

// --- Uninitialized ---

func TestCheck_Uninitialized_ReturnsError(t *testing.T) {
	dir := t.TempDir()

	_, err := tickets.Check(dir)
	if err == nil {
		t.Fatal("expected error for uninitialized directory, got nil")
	}
}
