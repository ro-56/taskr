package tickets_test

import (
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ro-56/taskr/internal/tickets"
)

func TestClose_WithSummary_AppendsNotes(t *testing.T) {
	dir := t.TempDir()
	id := setupTicket(t, dir, tickets.AddOptions{Title: "summary test"})

	if err := tickets.Close(dir, id, "All done, shipped."); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := tickets.Show(dir, id)
	if err != nil {
		t.Fatalf("show after close: %v", err)
	}

	if !strings.Contains(result.Body, "## Notes") {
		t.Errorf("body missing ## Notes section: %q", result.Body)
	}
	if !strings.Contains(result.Body, "All done, shipped.") {
		t.Errorf("body missing summary text: %q", result.Body)
	}
}

func TestClose_AlreadyClosed_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	id := setupTicket(t, dir, tickets.AddOptions{Title: "double close"})

	if err := tickets.Close(dir, id, ""); err != nil {
		t.Fatalf("first close: %v", err)
	}

	err := tickets.Close(dir, id, "")
	if !errors.Is(err, tickets.ErrWrongStatus) {
		t.Errorf("expected ErrWrongStatus, got %v", err)
	}
}

func TestClose_SetsClosedAndArchives(t *testing.T) {
	dir := t.TempDir()
	id := setupTicket(t, dir, tickets.AddOptions{Title: "close me"})

	if err := tickets.Close(dir, id, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := tickets.Show(dir, id)
	if err != nil {
		t.Fatalf("show after close: %v", err)
	}
	if result.Frontmatter["status"] != "closed" {
		t.Errorf("status: got %v, want closed", result.Frontmatter["status"])
	}

	// original file gone, archive file exists
	if _, err := os.Stat(dir + "/.tickets/" + id + ".md"); !os.IsNotExist(err) {
		t.Error("original ticket file should not exist after close")
	}
	if _, err := os.Stat(dir + "/.tickets/archive/" + id + ".md"); err != nil {
		t.Errorf("archive file missing: %v", err)
	}
}

func TestStart_ClosedTicket_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	id := setupTicket(t, dir, tickets.AddOptions{Title: "closed ticket"})

	// close it first (manually set status to simulate closed state)
	if err := tickets.Close(dir, id, ""); err != nil {
		t.Fatalf("close setup: %v", err)
	}

	err := tickets.Start(dir, id)
	if !errors.Is(err, tickets.ErrWrongStatus) {
		t.Errorf("expected ErrWrongStatus, got %v", err)
	}
}

func TestStart_PartialID(t *testing.T) {
	dir := t.TempDir()
	id := setupTicket(t, dir, tickets.AddOptions{Title: "partial id"})

	if err := tickets.Start(dir, id[:10]); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := tickets.Show(dir, id)
	if result.Frontmatter["status"] != "in_progress" {
		t.Errorf("status: got %v, want in_progress", result.Frontmatter["status"])
	}
}

func TestStart_UpdatesTimestamp(t *testing.T) {
	dir := t.TempDir()
	id := setupTicket(t, dir, tickets.AddOptions{Title: "timestamp test"})

	before := time.Now().UTC().Add(-time.Second)

	if err := tickets.Start(dir, id); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := tickets.Show(dir, id)
	raw, _ := result.Frontmatter["updated"].(string)
	updated, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		t.Fatalf("updated not a valid RFC3339 timestamp: %q", raw)
	}
	if !updated.After(before) {
		t.Errorf("updated %v not after start time %v", updated, before)
	}
}

func TestStart_SetsInProgress(t *testing.T) {
	dir := t.TempDir()
	id := setupTicket(t, dir, tickets.AddOptions{Title: "start me"})

	if err := tickets.Start(dir, id); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := tickets.Show(dir, id)
	if err != nil {
		t.Fatalf("show after start: %v", err)
	}
	if got := result.Frontmatter["status"]; got != "in_progress" {
		t.Errorf("status: got %q, want %q", got, "in_progress")
	}
}
