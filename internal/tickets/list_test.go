package tickets_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/ro-56/taskr/internal/tickets"
)

func TestList_ReturnsActiveTickets(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	id1, _ := tickets.Add(dir, tickets.AddOptions{Title: "Alpha"})
	id2, _ := tickets.Add(dir, tickets.AddOptions{Title: "Beta"})

	results, err := tickets.List(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 tickets, got %d", len(results))
	}

	byID := make(map[string]tickets.TicketSummary)
	for _, r := range results {
		byID[r.ID] = r
	}

	for _, id := range []string{id1, id2} {
		r, ok := byID[id]
		if !ok {
			t.Errorf("ticket %s not in results", id)
			continue
		}
		if r.Status != "open" {
			t.Errorf("ticket %s: status = %q, want open", id, r.Status)
		}
	}
	if byID[id1].Title != "Alpha" {
		t.Errorf("title: got %q, want Alpha", byID[id1].Title)
	}
}

func TestList_ExcludesArchive(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	id, _ := tickets.Add(dir, tickets.AddOptions{Title: "Archived"})

	src := fmt.Sprintf("%s/.tickets/%s.md", dir, id)
	dst := fmt.Sprintf("%s/.tickets/archive/%s.md", dir, id)
	if err := os.Rename(src, dst); err != nil {
		t.Fatal(err)
	}

	results, err := tickets.List(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.ID == id {
			t.Errorf("archived ticket %s should not appear in List", id)
		}
	}
}

func TestListTags_SortedUnique(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	tickets.Add(dir, tickets.AddOptions{Title: "A", Tags: []string{"zebra", "alpha"}})
	tickets.Add(dir, tickets.AddOptions{Title: "B", Tags: []string{"alpha", "beta"}})

	tags, err := tickets.ListTags(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"alpha", "beta", "zebra"}
	if len(tags) != len(want) {
		t.Fatalf("tags: got %v, want %v", tags, want)
	}
	for i, w := range want {
		if tags[i] != w {
			t.Errorf("tags[%d]: got %q, want %q", i, tags[i], w)
		}
	}
}

func TestListTags_EmptyWhenNoTags(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	tickets.Add(dir, tickets.AddOptions{Title: "No tags here"})

	tags, err := tickets.ListTags(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags) != 0 {
		t.Errorf("expected no tags, got %v", tags)
	}
}

func TestListCounts_IncludesZeros(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")

	counts, err := tickets.ListCounts(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, status := range []string{"open", "in_progress", "closed"} {
		if _, ok := counts[status]; !ok {
			t.Errorf("status %q missing from counts", status)
		}
	}
	if counts["open"] != 0 || counts["in_progress"] != 0 || counts["closed"] != 0 {
		t.Errorf("expected all zeros, got %v", counts)
	}
}

func TestListCounts_CorrectCounts(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")

	tickets.Add(dir, tickets.AddOptions{Title: "Open 1"})
	tickets.Add(dir, tickets.AddOptions{Title: "Open 2"})

	inProgressID, _ := tickets.Add(dir, tickets.AddOptions{Title: "In Progress"})
	tickets.Start(dir, inProgressID)

	closedID, _ := tickets.Add(dir, tickets.AddOptions{Title: "Closed"})
	tickets.Close(dir, closedID, "")

	counts, err := tickets.ListCounts(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if counts["open"] != 2 {
		t.Errorf("open: got %d, want 2", counts["open"])
	}
	if counts["in_progress"] != 1 {
		t.Errorf("in_progress: got %d, want 1", counts["in_progress"])
	}
	if counts["closed"] != 1 {
		t.Errorf("closed: got %d, want 1", counts["closed"])
	}
}

func TestListByStatus_Open(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")

	openID, _ := tickets.Add(dir, tickets.AddOptions{Title: "Open ticket"})
	inProgressID, _ := tickets.Add(dir, tickets.AddOptions{Title: "In progress"})
	tickets.Start(dir, inProgressID)

	results, err := tickets.ListByStatus(dir, "open")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d: %v", len(results), results)
	}
	if results[0].ID != openID {
		t.Errorf("ID: got %q, want %q", results[0].ID, openID)
	}
	if results[0].Status != "open" {
		t.Errorf("status: got %q, want open", results[0].Status)
	}
}

func TestListByStatus_ClosedIncludesArchive(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")

	id1, _ := tickets.Add(dir, tickets.AddOptions{Title: "Will close"})
	id2, _ := tickets.Add(dir, tickets.AddOptions{Title: "Will close 2"})
	tickets.Close(dir, id1, "")
	tickets.Close(dir, id2, "")
	tickets.Add(dir, tickets.AddOptions{Title: "Still open"})

	results, err := tickets.ListByStatus(dir, "closed")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 closed tickets, got %d: %v", len(results), results)
	}
	for _, r := range results {
		if r.Status != "closed" {
			t.Errorf("ticket %s: status = %q, want closed", r.ID, r.Status)
		}
	}
}
