package tickets_test

import (
	"errors"
	"testing"
	"time"

	"github.com/ro-56/taskr/internal/tickets"
)

// --- Bare-number IDs ---

func TestLink_AcceptsBareNumberIDs(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID, _ := tickets.Add(dir, tickets.AddOptions{Title: "A"}) // TKT-1
	bID, _ := tickets.Add(dir, tickets.AddOptions{Title: "B"}) // TKT-2

	if err := tickets.Link(dir, "1", "2"); err != nil {
		t.Fatalf("Link with bare-number IDs: %v", err)
	}

	result, _ := tickets.Show(dir, aID)
	if len(result.DependsOn) != 1 || result.DependsOn[0].ID != bID {
		t.Errorf("DependsOn after bare-number link: got %v", result.DependsOn)
	}
}

func TestUnlink_AcceptsBareNumberIDs(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID, _ := tickets.Add(dir, tickets.AddOptions{Title: "A"}) // TKT-1
	bID, _ := tickets.Add(dir, tickets.AddOptions{Title: "B"}) // TKT-2

	tickets.Link(dir, aID, bID)
	if err := tickets.Unlink(dir, "1", "2"); err != nil {
		t.Fatalf("Unlink with bare-number IDs: %v", err)
	}

	result, _ := tickets.Show(dir, aID)
	if len(result.DependsOn) != 0 {
		t.Errorf("DependsOn after bare-number unlink: got %v", result.DependsOn)
	}
}

func TestPrune_AcceptsBareNumberID(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID, _ := tickets.Add(dir, tickets.AddOptions{Title: "A"}) // TKT-1
	bID, _ := tickets.Add(dir, tickets.AddOptions{Title: "B"}) // TKT-2

	tickets.Link(dir, bID, aID)
	if err := tickets.Prune(dir, "1"); err != nil {
		t.Fatalf("Prune with bare-number ID: %v", err)
	}

	result, _ := tickets.Show(dir, bID)
	if len(result.DependsOn) != 0 {
		t.Errorf("B.DependsOn after bare-number prune: got %v", result.DependsOn)
	}
}

// --- Prune ---

func TestPrune_ClearsOwnDependencies(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID, _ := tickets.Add(dir, tickets.AddOptions{Title: "A"})
	bID, _ := tickets.Add(dir, tickets.AddOptions{Title: "B"})

	tickets.Link(dir, aID, bID) // A depends on B

	if err := tickets.Prune(dir, aID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := tickets.Show(dir, aID)
	if len(result.DependsOn) != 0 {
		t.Errorf("DependsOn: expected empty after prune, got %v", result.DependsOn)
	}
}

func TestPrune_RemovesFromOtherTickets(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID, _ := tickets.Add(dir, tickets.AddOptions{Title: "A"})
	bID, _ := tickets.Add(dir, tickets.AddOptions{Title: "B"})
	cID, _ := tickets.Add(dir, tickets.AddOptions{Title: "C"})

	tickets.Link(dir, bID, aID) // B depends on A
	tickets.Link(dir, cID, aID) // C depends on A

	if err := tickets.Prune(dir, aID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	bResult, _ := tickets.Show(dir, bID)
	if len(bResult.DependsOn) != 0 {
		t.Errorf("B.DependsOn: expected empty after prune of A, got %v", bResult.DependsOn)
	}
	cResult, _ := tickets.Show(dir, cID)
	if len(cResult.DependsOn) != 0 {
		t.Errorf("C.DependsOn: expected empty after prune of A, got %v", cResult.DependsOn)
	}
}

// --- Unlink ---

func TestUnlink_RemovesDependency(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID, _ := tickets.Add(dir, tickets.AddOptions{Title: "A"})
	bID, _ := tickets.Add(dir, tickets.AddOptions{Title: "B"})

	tickets.Link(dir, aID, bID)
	if err := tickets.Unlink(dir, aID, bID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := tickets.Show(dir, aID)
	if len(result.DependsOn) != 0 {
		t.Errorf("DependsOn: expected empty after unlink, got %v", result.DependsOn)
	}
}

func TestUnlink_UpdatesTimestamp(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID, _ := tickets.Add(dir, tickets.AddOptions{Title: "A"})
	bID, _ := tickets.Add(dir, tickets.AddOptions{Title: "B"})

	tickets.Link(dir, aID, bID)
	before := time.Now().UTC().Truncate(time.Second)
	tickets.Unlink(dir, aID, bID)

	result, _ := tickets.Show(dir, aID)
	updatedStr, _ := result.Frontmatter["updated"].(string)
	updated, err := time.Parse(time.RFC3339, updatedStr)
	if err != nil {
		t.Fatalf("could not parse updated: %v", err)
	}
	if updated.Before(before) {
		t.Errorf("updated %v is before pre-call time %v", updated, before)
	}
}

// --- Link ---

func TestLink_UpdatesTimestamp(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID, _ := tickets.Add(dir, tickets.AddOptions{Title: "A"})
	bID, _ := tickets.Add(dir, tickets.AddOptions{Title: "B"})

	before := time.Now().UTC().Truncate(time.Second)

	tickets.Link(dir, aID, bID)

	result, _ := tickets.Show(dir, aID)
	updatedStr, _ := result.Frontmatter["updated"].(string)
	updated, err := time.Parse(time.RFC3339, updatedStr)
	if err != nil {
		t.Fatalf("could not parse updated: %v", err)
	}
	// RFC3339 is second-precision; updated must be at or after the pre-call second
	if updated.Before(before) {
		t.Errorf("updated %v is before pre-call time %v", updated, before)
	}
}

func TestLink_IndirectCycle_Error(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID, _ := tickets.Add(dir, tickets.AddOptions{Title: "A"})
	bID, _ := tickets.Add(dir, tickets.AddOptions{Title: "B"})
	cID, _ := tickets.Add(dir, tickets.AddOptions{Title: "C"})

	tickets.Link(dir, aID, bID) // A → B
	tickets.Link(dir, bID, cID) // B → C
	err := tickets.Link(dir, cID, aID) // C → A (cycle: A→B→C→A)
	var cycleErr *tickets.ErrCycle
	if !errors.As(err, &cycleErr) {
		t.Fatalf("expected ErrCycle, got %v", err)
	}
	// path should contain at least cID, aID, bID, cID or similar
	if len(cycleErr.Path) < 3 {
		t.Errorf("cycle path too short: %v", cycleErr.Path)
	}
}

func TestLink_DirectCycle_Error(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID, _ := tickets.Add(dir, tickets.AddOptions{Title: "A"})
	bID, _ := tickets.Add(dir, tickets.AddOptions{Title: "B"})

	tickets.Link(dir, aID, bID) // A → B
	err := tickets.Link(dir, bID, aID) // B → A (cycle)
	var cycleErr *tickets.ErrCycle
	if !errors.As(err, &cycleErr) {
		t.Fatalf("expected ErrCycle, got %v", err)
	}
	if len(cycleErr.Path) < 2 {
		t.Errorf("cycle path too short: %v", cycleErr.Path)
	}
}

func TestLink_AlreadyLinked_NoOp(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID, _ := tickets.Add(dir, tickets.AddOptions{Title: "A"})
	bID, _ := tickets.Add(dir, tickets.AddOptions{Title: "B"})

	tickets.Link(dir, aID, bID)
	err := tickets.Link(dir, aID, bID)
	if !errors.Is(err, tickets.ErrAlreadyLinked) {
		t.Errorf("expected ErrAlreadyLinked, got %v", err)
	}

	result, _ := tickets.Show(dir, aID)
	if len(result.DependsOn) != 1 {
		t.Errorf("DependsOn: expected 1 entry after duplicate link, got %d", len(result.DependsOn))
	}
}

func TestLink_SelfLink_Error(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID, _ := tickets.Add(dir, tickets.AddOptions{Title: "A"})

	err := tickets.Link(dir, aID, aID)
	if !errors.Is(err, tickets.ErrSelfLink) {
		t.Errorf("expected ErrSelfLink, got %v", err)
	}
}

func TestLink_AddsDependency(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID, _ := tickets.Add(dir, tickets.AddOptions{Title: "A"})
	bID, _ := tickets.Add(dir, tickets.AddOptions{Title: "B"})

	if err := tickets.Link(dir, aID, bID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := tickets.Show(dir, aID)
	if len(result.DependsOn) != 1 || result.DependsOn[0].ID != bID {
		t.Errorf("DependsOn: got %v, want [{%s open}]", result.DependsOn, bID)
	}
}
