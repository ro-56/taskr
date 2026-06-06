package tickets_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/ro-56/taskr/internal/tickets"
)

func ptr[T any](v T) *T { return &v }

func TestUpdate_Title(t *testing.T) {
	dir := t.TempDir()
	id := setupTicket(t, dir, tickets.AddOptions{Title: "old title"})

	if err := tickets.Update(dir, id, tickets.UpdateOptions{Title: ptr("new title")}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := tickets.Show(dir, id)
	if got := result.Frontmatter["title"]; got != "new title" {
		t.Errorf("title: got %q, want %q", got, "new title")
	}
}

func TestUpdate_Priority(t *testing.T) {
	dir := t.TempDir()
	id := setupTicket(t, dir, tickets.AddOptions{Title: "prio test"})

	if err := tickets.Update(dir, id, tickets.UpdateOptions{Priority: ptr(0)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := tickets.Show(dir, id)
	if got := fmt.Sprint(result.Frontmatter["priority"]); got != "0" {
		t.Errorf("priority: got %q, want %q", got, "0")
	}
}

func TestUpdate_InvalidPriority(t *testing.T) {
	dir := t.TempDir()
	id := setupTicket(t, dir, tickets.AddOptions{Title: "invalid prio"})

	for _, p := range []int{-1, 4} {
		err := tickets.Update(dir, id, tickets.UpdateOptions{Priority: ptr(p)})
		if !errors.Is(err, tickets.ErrInvalidPriority) {
			t.Errorf("priority %d: expected ErrInvalidPriority, got %v", p, err)
		}
	}
}

func TestUpdate_Mode(t *testing.T) {
	dir := t.TempDir()
	id := setupTicket(t, dir, tickets.AddOptions{Title: "mode test", Mode: "hitl"})

	if err := tickets.Update(dir, id, tickets.UpdateOptions{Mode: ptr("afk")}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := tickets.Show(dir, id)
	if got := result.Frontmatter["mode"]; got != "afk" {
		t.Errorf("mode: got %q, want %q", got, "afk")
	}
}

func TestUpdate_InvalidMode(t *testing.T) {
	dir := t.TempDir()
	id := setupTicket(t, dir, tickets.AddOptions{Title: "invalid mode"})

	err := tickets.Update(dir, id, tickets.UpdateOptions{Mode: ptr("auto")})
	if !errors.Is(err, tickets.ErrInvalidMode) {
		t.Errorf("expected ErrInvalidMode, got %v", err)
	}
}

func TestUpdate_Tags(t *testing.T) {
	dir := t.TempDir()
	id := setupTicket(t, dir, tickets.AddOptions{Title: "tags test", Tags: []string{"old"}})

	if err := tickets.Update(dir, id, tickets.UpdateOptions{Tags: []string{"p0", "auth"}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := tickets.Show(dir, id)
	tags, _ := result.Frontmatter["tags"].([]any)
	if len(tags) != 2 || fmt.Sprint(tags[0]) != "p0" || fmt.Sprint(tags[1]) != "auth" {
		t.Errorf("tags: got %v, want [p0 auth]", tags)
	}
}

func TestUpdate_MultipleFlags(t *testing.T) {
	dir := t.TempDir()
	id := setupTicket(t, dir, tickets.AddOptions{Title: "multi"})

	opts := tickets.UpdateOptions{
		Title:    ptr("updated title"),
		Priority: ptr(1),
		Mode:     ptr("afk"),
		Tags:     []string{"p1", "backend"},
	}
	if err := tickets.Update(dir, id, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := tickets.Show(dir, id)
	fm := result.Frontmatter
	if got := fm["title"]; got != "updated title" {
		t.Errorf("title: got %q", got)
	}
	if got := fmt.Sprint(fm["priority"]); got != "1" {
		t.Errorf("priority: got %q", got)
	}
	if got := fm["mode"]; got != "afk" {
		t.Errorf("mode: got %q", got)
	}
	tags, _ := fm["tags"].([]any)
	if len(tags) != 2 {
		t.Errorf("tags: got %v", tags)
	}
}

func TestUpdate_UnspecifiedFieldsUnchanged(t *testing.T) {
	dir := t.TempDir()
	p := 1
	id := setupTicket(t, dir, tickets.AddOptions{
		Title:    "unchanged",
		Type:     "bug",
		Priority: &p,
		Mode:     "afk",
		Tags:     []string{"x"},
	})

	if err := tickets.Update(dir, id, tickets.UpdateOptions{Title: ptr("new")}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := tickets.Show(dir, id)
	fm := result.Frontmatter
	if got := fmt.Sprint(fm["priority"]); got != "1" {
		t.Errorf("priority changed: got %q", got)
	}
	if got := fm["mode"]; got != "afk" {
		t.Errorf("mode changed: got %q", got)
	}
	if got := fm["type"]; got != "bug" {
		t.Errorf("type changed: got %q", got)
	}
	tags, _ := fm["tags"].([]any)
	if len(tags) != 1 || fmt.Sprint(tags[0]) != "x" {
		t.Errorf("tags changed: got %v", tags)
	}
}

func TestUpdate_UpdatesTimestamp(t *testing.T) {
	dir := t.TempDir()
	id := setupTicket(t, dir, tickets.AddOptions{Title: "ts test"})

	before := time.Now().UTC().Add(-time.Second)
	if err := tickets.Update(dir, id, tickets.UpdateOptions{Title: ptr("new")}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := tickets.Show(dir, id)
	raw, _ := result.Frontmatter["updated"].(string)
	ts, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		t.Fatalf("updated not valid RFC3339: %q", raw)
	}
	if !ts.After(before) {
		t.Errorf("updated %v not after start time %v", ts, before)
	}
}

func TestUpdate_PartialID(t *testing.T) {
	dir := t.TempDir()
	id := setupTicket(t, dir, tickets.AddOptions{Title: "partial"})

	if err := tickets.Update(dir, id[:10], tickets.UpdateOptions{Title: ptr("patched")}); err != nil {
		t.Fatalf("unexpected error with partial ID: %v", err)
	}

	result, _ := tickets.Show(dir, id)
	if got := result.Frontmatter["title"]; got != "patched" {
		t.Errorf("title: got %q, want %q", got, "patched")
	}
}

func TestUpdate_Body_WithTitle(t *testing.T) {
	dir := t.TempDir()
	id := setupTicket(t, dir, tickets.AddOptions{Title: "old"})

	opts := tickets.UpdateOptions{
		Title: ptr("new title"),
		Body:  ptr("## Notes\n\ndetails"),
	}
	if err := tickets.Update(dir, id, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := tickets.Show(dir, id)
	if got := result.Frontmatter["title"]; got != "new title" {
		t.Errorf("title: got %q", got)
	}
	if got := result.Body; got != "## Notes\n\ndetails" {
		t.Errorf("body: got %q", got)
	}
}

func TestUpdate_Body_Preserved(t *testing.T) {
	dir := t.TempDir()
	id := setupTicket(t, dir, tickets.AddOptions{Title: "preserve body"})
	tickets.Update(dir, id, tickets.UpdateOptions{Body: ptr("keep this")})

	if err := tickets.Update(dir, id, tickets.UpdateOptions{Title: ptr("new title")}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := tickets.Show(dir, id)
	if got := result.Body; got != "keep this" {
		t.Errorf("body should be unchanged, got %q", got)
	}
}

func TestUpdate_Body_Clear(t *testing.T) {
	dir := t.TempDir()
	id := setupTicket(t, dir, tickets.AddOptions{Title: "clear body"})
	tickets.Update(dir, id, tickets.UpdateOptions{Body: ptr("some existing content")})

	if err := tickets.Update(dir, id, tickets.UpdateOptions{Body: ptr("")}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := tickets.Show(dir, id)
	if got := result.Body; got != "" {
		t.Errorf("body should be empty after clear, got %q", got)
	}
}

func TestUpdate_Body(t *testing.T) {
	dir := t.TempDir()
	id := setupTicket(t, dir, tickets.AddOptions{Title: "body test"})

	before := time.Now().UTC().Add(-time.Second)
	if err := tickets.Update(dir, id, tickets.UpdateOptions{Body: ptr("## Description\n\nsome details")}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := tickets.Show(dir, id)
	if got := result.Body; got != "## Description\n\nsome details" {
		t.Errorf("body: got %q, want %q", got, "## Description\n\nsome details")
	}
	raw, _ := result.Frontmatter["updated"].(string)
	ts, _ := time.Parse(time.RFC3339, raw)
	if !ts.After(before) {
		t.Errorf("updated not bumped: %v", ts)
	}
}
