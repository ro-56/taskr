package tickets_test

import (
	"testing"

	"github.com/ro-56/taskr/internal/tickets"
)

// --- RenderDepTree ---

func TestRenderDepTree_SingleNode(t *testing.T) {
	node := &tickets.DepTreeNode{ID: "TKT-aabbccdd", Status: "open"}
	got := tickets.RenderDepTree(node)
	want := "TKT-aabbccdd (open)\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRenderDepTree_MultipleChildren(t *testing.T) {
	node := &tickets.DepTreeNode{
		ID:     "TKT-root1234",
		Status: "in_progress",
		Children: []*tickets.DepTreeNode{
			{ID: "TKT-aaaa0001", Status: "closed"},
			{ID: "TKT-bbbb0002", Status: "open"},
		},
	}
	got := tickets.RenderDepTree(node)
	want := "TKT-root1234 (in_progress)\n" +
		"├── TKT-aaaa0001 (closed)\n" +
		"└── TKT-bbbb0002 (open)\n"
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestRenderDepTree_NestedChildren(t *testing.T) {
	node := &tickets.DepTreeNode{
		ID:     "TKT-root1234",
		Status: "in_progress",
		Children: []*tickets.DepTreeNode{
			{ID: "TKT-aaaa0001", Status: "closed"},
			{
				ID:     "TKT-bbbb0002",
				Status: "open",
				Children: []*tickets.DepTreeNode{
					{ID: "TKT-cccc0003", Status: "open"},
				},
			},
		},
	}
	got := tickets.RenderDepTree(node)
	want := "TKT-root1234 (in_progress)\n" +
		"├── TKT-aaaa0001 (closed)\n" +
		"└── TKT-bbbb0002 (open)\n" +
		"    └── TKT-cccc0003 (open)\n"
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

// --- DepTree ---

func TestDepTree_FirstDegree_ChildrenNotExpanded(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID, _ := tickets.Add(dir, tickets.AddOptions{Title: "A"})
	bID, _ := tickets.Add(dir, tickets.AddOptions{Title: "B"})
	cID, _ := tickets.Add(dir, tickets.AddOptions{Title: "C"})
	tickets.Link(dir, aID, bID) // A → B
	tickets.Link(dir, bID, cID) // B → C

	node, err := tickets.DepTree(dir, aID, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(node.Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(node.Children))
	}
	child := node.Children[0]
	if child.ID != bID {
		t.Errorf("child ID: got %q, want %q", child.ID, bID)
	}
	if len(child.Children) != 0 {
		t.Errorf("full=false must not expand grandchildren, got %d", len(child.Children))
	}
}

func TestDepTree_Full_ExpandsRecursively(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID, _ := tickets.Add(dir, tickets.AddOptions{Title: "A"})
	bID, _ := tickets.Add(dir, tickets.AddOptions{Title: "B"})
	cID, _ := tickets.Add(dir, tickets.AddOptions{Title: "C"})
	tickets.Link(dir, aID, bID) // A → B
	tickets.Link(dir, bID, cID) // B → C

	node, err := tickets.DepTree(dir, aID, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(node.Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(node.Children))
	}
	b := node.Children[0]
	if len(b.Children) != 1 {
		t.Fatalf("expected B to have 1 child, got %d", len(b.Children))
	}
	c := b.Children[0]
	if c.ID != cID {
		t.Errorf("grandchild ID: got %q, want %q", c.ID, cID)
	}
}

func TestDepTree_AcceptsPartialID(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID, _ := tickets.Add(dir, tickets.AddOptions{Title: "A"})

	node, err := tickets.DepTree(dir, aID[:10], false)
	if err != nil {
		t.Fatalf("partial ID should resolve: %v", err)
	}
	if node.ID != aID {
		t.Errorf("ID: got %q, want %q", node.ID, aID)
	}
}

func TestDepTree_ArchivedRoot(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID, _ := tickets.Add(dir, tickets.AddOptions{Title: "A"})
	tickets.Start(dir, aID)
	tickets.Close(dir, aID, "")

	node, err := tickets.DepTree(dir, aID, false)
	if err != nil {
		t.Fatalf("archived ticket should work: %v", err)
	}
	if node.Status != "closed" {
		t.Errorf("status: got %q, want %q", node.Status, "closed")
	}
}

func TestDepTree_NoDeps_RootOnly(t *testing.T) {
	dir := t.TempDir()
	tickets.Init(dir, "TKT")
	aID, _ := tickets.Add(dir, tickets.AddOptions{Title: "A"})

	node, err := tickets.DepTree(dir, aID, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node.ID != aID {
		t.Errorf("root ID: got %q, want %q", node.ID, aID)
	}
	if node.Status != "open" {
		t.Errorf("root status: got %q, want %q", node.Status, "open")
	}
	if len(node.Children) != 0 {
		t.Errorf("expected no children, got %d", len(node.Children))
	}
}
