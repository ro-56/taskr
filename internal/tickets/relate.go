package tickets

import (
	"errors"
	"os"
	"path/filepath"
	"time"
)

var (
	ErrSelfRelate     = errors.New("a ticket cannot be related to itself")
	ErrAlreadyRelated = errors.New("tickets are already related")
)

func Relate(dir string, partialIDs ...string) error {
	ticketsDir := filepath.Join(dir, ".tickets")
	archiveDir := filepath.Join(ticketsDir, "archive")

	// Resolve all partial IDs upfront.
	paths := make([]string, len(partialIDs))
	for i, pid := range partialIDs {
		p, err := resolveID(pid, ticketsDir, archiveDir)
		if err != nil {
			return err
		}
		paths[i] = p
	}

	// Read all tickets.
	type entry struct {
		path string
		fm   map[string]any
		body string
	}
	entries := make([]entry, len(paths))
	ids := make([]string, len(paths))
	for i, p := range paths {
		raw, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		fm, body, err := parseFrontmatter(raw)
		if err != nil {
			return err
		}
		entries[i] = entry{path: p, fm: fm, body: body}
		ids[i], _ = fm["id"].(string)
	}

	// Reject duplicate IDs (self-relate).
	seen := map[string]bool{}
	for _, id := range ids {
		if seen[id] {
			return ErrSelfRelate
		}
		seen[id] = true
	}

	// Track current links per entry in a []string so repeated appends within
	// the same call are visible without re-reading fm (which holds []any).
	currentLinks := make([][]string, len(entries))
	for i := range entries {
		currentLinks[i] = stringSlice(entries[i].fm["links"])
	}

	now := time.Now().UTC().Format(time.RFC3339)
	newLinks := 0
	modified := make([]bool, len(entries))

	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			idI, idJ := ids[i], ids[j]
			alreadyLinked := false
			for _, l := range currentLinks[i] {
				if l == idJ {
					alreadyLinked = true
					break
				}
			}
			if alreadyLinked {
				continue
			}

			newLinks++
			currentLinks[i] = append(currentLinks[i], idJ)
			entries[i].fm["updated"] = now
			modified[i] = true

			currentLinks[j] = append(currentLinks[j], idI)
			entries[j].fm["updated"] = now
			modified[j] = true
		}
	}

	if newLinks == 0 {
		return ErrAlreadyRelated
	}

	for i, e := range entries {
		if modified[i] {
			e.fm["links"] = currentLinks[i]
			if err := writeTicket(e.path, e.fm, e.body); err != nil {
				return err
			}
		}
	}

	return nil
}
