package tickets_test

import (
	"testing"
	"time"

	"github.com/ro-56/taskr/internal/tickets"
)

func TestReady_SortsByPriorityAsc(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	p0 := intPtr(0)
	p2 := intPtr(2)
	p3 := intPtr(3)

	midID, _ := tickets.Add(dir, tickets.AddOptions{Title: "mid", Priority: p2})
	critID, _ := tickets.Add(dir, tickets.AddOptions{Title: "crit", Priority: p0})
	lowID, _ := tickets.Add(dir, tickets.AddOptions{Title: "low", Priority: p3})

	result, err := tickets.Ready(dir, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ids := summaryIDs(result)

	pos := func(id string) int {
		for i, v := range ids {
			if v == id {
				return i
			}
		}
		return -1
	}
	if pos(critID) >= pos(midID) || pos(midID) >= pos(lowID) {
		t.Errorf("expected crit < mid < low; order: %v", ids)
	}
}

func TestReady_ModeFilterIncludesOnlyMatchingMode(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	afkID, _ := tickets.Add(dir, tickets.AddOptions{Title: "afk task", Mode: "afk"})
	_, _ = tickets.Add(dir, tickets.AddOptions{Title: "hitl task", Mode: "hitl"})

	result, err := tickets.Ready(dir, "afk")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ids := summaryIDs(result)
	if len(ids) != 1 || ids[0] != afkID {
		t.Errorf("expected only afk ticket, got %v", ids)
	}
}

func TestReady_TiesBrokenByUpdatedDesc(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	p1 := intPtr(1)

	olderID, _ := tickets.Add(dir, tickets.AddOptions{Title: "older", Priority: p1})
	time.Sleep(1100 * time.Millisecond) // ensure different RFC3339 second
	newerID, _ := tickets.Add(dir, tickets.AddOptions{Title: "newer", Priority: p1})

	result, err := tickets.Ready(dir, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ids := summaryIDs(result)

	pos := func(id string) int {
		for i, v := range ids {
			if v == id {
				return i
			}
		}
		return -1
	}
	if pos(newerID) >= pos(olderID) {
		t.Errorf("newer (higher updated) should appear before older; order: %v", ids)
	}
}

func intPtr(n int) *int { return &n }

func TestReady_ExcludesClosedTickets(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	id, _ := tickets.Add(dir, tickets.AddOptions{Title: "done"})
	tickets.Close(dir, id, "")

	result, err := tickets.Ready(dir, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty after close, got %v", summaryIDs(result))
	}
}

func TestReady_ExcludesTicketBlockedByOpenDep(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	depID, _ := tickets.Add(dir, tickets.AddOptions{Title: "dep"})
	id, _ := tickets.Add(dir, tickets.AddOptions{Title: "main"})
	tickets.Link(dir, id, depID) // id depends on depID (still open)

	result, err := tickets.Ready(dir, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ids := summaryIDs(result)
	if contains(ids, id) {
		t.Errorf("blocked ticket %s should not appear in ready list", id)
	}
}

func TestReady_ExcludesTicketBlockedByInProgressDep(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	depID, _ := tickets.Add(dir, tickets.AddOptions{Title: "dep"})
	id, _ := tickets.Add(dir, tickets.AddOptions{Title: "main"})
	tickets.Link(dir, id, depID)
	tickets.Start(dir, depID) // dep is in_progress

	result, err := tickets.Ready(dir, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ids := summaryIDs(result)
	if contains(ids, id) {
		t.Errorf("blocked ticket %s should not appear in ready list", id)
	}
}

func TestReady_IncludesInProgressWithAllClosedDeps(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	depID, _ := tickets.Add(dir, tickets.AddOptions{Title: "dep"})
	id, _ := tickets.Add(dir, tickets.AddOptions{Title: "main"})
	tickets.Link(dir, id, depID)
	tickets.Close(dir, depID, "")

	result, err := tickets.Ready(dir, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ids := summaryIDs(result)
	if !contains(ids, id) {
		t.Errorf("expected %s in ready list, got %v", id, ids)
	}
	if contains(ids, depID) {
		t.Errorf("closed dep %s should not appear in ready list", depID)
	}
}

func TestReady_IncludesOpenTicketWithNoDeps(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	id, _ := tickets.Add(dir, tickets.AddOptions{Title: "do something"})

	result, err := tickets.Ready(dir, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0].ID != id {
		t.Errorf("expected [%s], got %v", id, result)
	}
}

func TestReady_EmptyWhenNoTickets(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")

	result, err := tickets.Ready(dir, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty, got %v", result)
	}
}

// --- Blocked ---

func TestBlocked_EmptyWhenNoTickets(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")

	result, err := tickets.Blocked(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty, got %v", result)
	}
}

func TestBlocked_IncludesTicketWithOpenDep(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	depID, _ := tickets.Add(dir, tickets.AddOptions{Title: "dep"})
	id, _ := tickets.Add(dir, tickets.AddOptions{Title: "main"})
	tickets.Link(dir, id, depID)

	result, err := tickets.Blocked(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ids := summaryIDs(result)
	if !contains(ids, id) {
		t.Errorf("expected %s in blocked list, got %v", id, ids)
	}
}

func TestBlocked_IncludesTicketWithInProgressDep(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	depID, _ := tickets.Add(dir, tickets.AddOptions{Title: "dep"})
	id, _ := tickets.Add(dir, tickets.AddOptions{Title: "main"})
	tickets.Link(dir, id, depID)
	tickets.Start(dir, depID)

	result, err := tickets.Blocked(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ids := summaryIDs(result)
	if !contains(ids, id) {
		t.Errorf("expected %s in blocked list, got %v", id, ids)
	}
}

func TestBlocked_ExcludesClosedTickets(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	depID, _ := tickets.Add(dir, tickets.AddOptions{Title: "dep"})
	id, _ := tickets.Add(dir, tickets.AddOptions{Title: "main"})
	tickets.Link(dir, id, depID)
	tickets.Close(dir, id, "") // close the blocked ticket itself

	result, err := tickets.Blocked(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ids := summaryIDs(result)
	if contains(ids, id) {
		t.Errorf("closed ticket %s should not appear in blocked list", id)
	}
}

// helpers

func summaryIDs(ss []tickets.TicketSummary) []string {
	ids := make([]string, len(ss))
	for i, s := range ss {
		ids[i] = s.ID
	}
	return ids
}

func contains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}
