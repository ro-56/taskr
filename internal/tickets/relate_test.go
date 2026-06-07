package tickets_test

import (
	"errors"
	"testing"
	"time"

	"github.com/ro-56/taskr/internal/tickets"
)

func TestRelate_AddsBidirectionalLinks(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID, _ := tickets.Add(dir, tickets.AddOptions{Title: "A"})
	bID, _ := tickets.Add(dir, tickets.AddOptions{Title: "B"})

	if err := tickets.Relate(dir, aID, bID); err != nil {
		t.Fatalf("Relate: %v", err)
	}

	aResult, _ := tickets.Show(dir, aID)
	bResult, _ := tickets.Show(dir, bID)

	aLinks := aResult.Frontmatter["links"]
	bLinks := bResult.Frontmatter["links"]

	if !containsID(aLinks, bID) {
		t.Errorf("A.links does not contain B: %v", aLinks)
	}
	if !containsID(bLinks, aID) {
		t.Errorf("B.links does not contain A: %v", bLinks)
	}
}

func TestRelate_AlreadyRelated_NoOp(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID, _ := tickets.Add(dir, tickets.AddOptions{Title: "A"})
	bID, _ := tickets.Add(dir, tickets.AddOptions{Title: "B"})

	tickets.Relate(dir, aID, bID)
	err := tickets.Relate(dir, aID, bID)
	if !errors.Is(err, tickets.ErrAlreadyRelated) {
		t.Errorf("expected ErrAlreadyRelated, got %v", err)
	}

	aResult, _ := tickets.Show(dir, aID)
	bResult, _ := tickets.Show(dir, bID)

	if countID(aResult.Frontmatter["links"], bID) != 1 {
		t.Errorf("A.links should have exactly one entry for B: %v", aResult.Frontmatter["links"])
	}
	if countID(bResult.Frontmatter["links"], aID) != 1 {
		t.Errorf("B.links should have exactly one entry for A: %v", bResult.Frontmatter["links"])
	}
}

func TestRelate_SelfRelate_Error(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID, _ := tickets.Add(dir, tickets.AddOptions{Title: "A"})

	err := tickets.Relate(dir, aID, aID)
	if !errors.Is(err, tickets.ErrSelfRelate) {
		t.Errorf("expected ErrSelfRelate, got %v", err)
	}
}

func TestRelate_AcceptsBareNumberIDs(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID, _ := tickets.Add(dir, tickets.AddOptions{Title: "A"}) // TKT-1
	bID, _ := tickets.Add(dir, tickets.AddOptions{Title: "B"}) // TKT-2

	if err := tickets.Relate(dir, "1", "2"); err != nil {
		t.Fatalf("Relate with bare-number IDs: %v", err)
	}

	aResult, _ := tickets.Show(dir, aID)
	if !containsID(aResult.Frontmatter["links"], bID) {
		t.Errorf("A.links does not contain B after bare-number relate: %v", aResult.Frontmatter["links"])
	}
}

func TestRelate_UpdatesBothTimestamps(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID, _ := tickets.Add(dir, tickets.AddOptions{Title: "A"})
	bID, _ := tickets.Add(dir, tickets.AddOptions{Title: "B"})

	before := time.Now().UTC().Truncate(time.Second)
	tickets.Relate(dir, aID, bID)

	for _, id := range []string{aID, bID} {
		result, _ := tickets.Show(dir, id)
		updatedStr, _ := result.Frontmatter["updated"].(string)
		updated, err := time.Parse(time.RFC3339, updatedStr)
		if err != nil {
			t.Fatalf("could not parse updated for %s: %v", id, err)
		}
		if updated.Before(before) {
			t.Errorf("%s updated %v is before pre-call time %v", id, updated, before)
		}
	}
}

func TestRelate_ThreeTickets_LinksAllPairs(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID, _ := tickets.Add(dir, tickets.AddOptions{Title: "A"})
	bID, _ := tickets.Add(dir, tickets.AddOptions{Title: "B"})
	cID, _ := tickets.Add(dir, tickets.AddOptions{Title: "C"})

	if err := tickets.Relate(dir, aID, bID, cID); err != nil {
		t.Fatalf("Relate 3 tickets: %v", err)
	}

	aResult, _ := tickets.Show(dir, aID)
	bResult, _ := tickets.Show(dir, bID)
	cResult, _ := tickets.Show(dir, cID)

	if !containsID(aResult.Frontmatter["links"], bID) || !containsID(aResult.Frontmatter["links"], cID) {
		t.Errorf("A.links should contain B and C: %v", aResult.Frontmatter["links"])
	}
	if !containsID(bResult.Frontmatter["links"], aID) || !containsID(bResult.Frontmatter["links"], cID) {
		t.Errorf("B.links should contain A and C: %v", bResult.Frontmatter["links"])
	}
	if !containsID(cResult.Frontmatter["links"], aID) || !containsID(cResult.Frontmatter["links"], bID) {
		t.Errorf("C.links should contain A and B: %v", cResult.Frontmatter["links"])
	}
}

func TestRelate_PartiallyRelated_LinksRemainingPairs(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID, _ := tickets.Add(dir, tickets.AddOptions{Title: "A"})
	bID, _ := tickets.Add(dir, tickets.AddOptions{Title: "B"})
	cID, _ := tickets.Add(dir, tickets.AddOptions{Title: "C"})

	tickets.Relate(dir, aID, bID) // A↔B already exists

	if err := tickets.Relate(dir, aID, bID, cID); err != nil {
		t.Fatalf("expected success when some pairs new, got: %v", err)
	}

	aResult, _ := tickets.Show(dir, aID)
	cResult, _ := tickets.Show(dir, cID)

	if !containsID(aResult.Frontmatter["links"], cID) {
		t.Errorf("A.links should contain C: %v", aResult.Frontmatter["links"])
	}
	if !containsID(cResult.Frontmatter["links"], aID) || !containsID(cResult.Frontmatter["links"], bID) {
		t.Errorf("C.links should contain A and B: %v", cResult.Frontmatter["links"])
	}
	if countID(aResult.Frontmatter["links"], bID) != 1 {
		t.Errorf("A.links should have exactly one entry for B: %v", aResult.Frontmatter["links"])
	}
}

func TestRelate_AllPairsAlreadyRelated_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID, _ := tickets.Add(dir, tickets.AddOptions{Title: "A"})
	bID, _ := tickets.Add(dir, tickets.AddOptions{Title: "B"})
	cID, _ := tickets.Add(dir, tickets.AddOptions{Title: "C"})

	tickets.Relate(dir, aID, bID, cID)

	err := tickets.Relate(dir, aID, bID, cID)
	if !errors.Is(err, tickets.ErrAlreadyRelated) {
		t.Errorf("expected ErrAlreadyRelated when all pairs exist, got %v", err)
	}
}

// helpers

func containsID(v any, id string) bool {
	raw, _ := v.([]any)
	for _, r := range raw {
		if s, ok := r.(string); ok && s == id {
			return true
		}
	}
	return false
}

func countID(v any, id string) int {
	raw, _ := v.([]any)
	n := 0
	for _, r := range raw {
		if s, ok := r.(string); ok && s == id {
			n++
		}
	}
	return n
}
