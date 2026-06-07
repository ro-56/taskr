package tickets

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Violation struct {
	TicketID string
	Message  string
}

var (
	validStatuses  = map[string]bool{"open": true, "in_progress": true, "closed": true}
	validTypes     = map[string]bool{"bug": true, "feature": true, "task": true, "epic": true, "chore": true}
	validModes     = map[string]bool{"afk": true, "hitl": true}
	requiredFields = []string{"id", "title", "status", "type", "priority", "mode", "created", "updated"}
)

func Check(dir string) ([]Violation, error) {
	ticketsDir := filepath.Join(dir, ".tickets")
	if _, err := os.Stat(filepath.Join(ticketsDir, "config.json")); err != nil {
		return nil, ErrNotInitialized
	}

	cfg, err := loadConfig(ticketsDir)
	if err != nil {
		return nil, err
	}
	// Accepts both the new sequential format (no zero-padding) and the
	// legacy 8-hex-char format used by tickets created before the switch.
	idPattern := regexp.MustCompile(`^` + regexp.QuoteMeta(cfg.Prefix) + `-([1-9][0-9]*|[0-9a-f]{8})$`)

	archiveDir := filepath.Join(ticketsDir, "archive")

	// Load all tickets from both dirs
	type entry struct {
		path    string
		inArchive bool
		fm      map[string]any
	}
	var allEntries []entry
	allFMs := map[string]map[string]any{} // id -> fm

	for _, scanPath := range []struct {
		path      string
		inArchive bool
	}{{ticketsDir, false}, {archiveDir, true}} {
		entries, err := os.ReadDir(scanPath.path)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			raw, err := os.ReadFile(filepath.Join(scanPath.path, e.Name()))
			if err != nil {
				continue
			}
			fm, _, err := parseFrontmatter(raw)
			if err != nil {
				continue
			}
			allEntries = append(allEntries, entry{
				path:      filepath.Join(scanPath.path, e.Name()),
				inArchive: scanPath.inArchive,
				fm:        fm,
			})
			if id, ok := fm["id"].(string); ok && id != "" {
				allFMs[id] = fm
			}
		}
	}

	var violations []Violation
	add := func(id, msg string) {
		violations = append(violations, Violation{TicketID: id, Message: msg})
	}

	for _, e := range allEntries {
		fm := e.fm
		stem := strings.TrimSuffix(filepath.Base(e.path), ".md")
		rawID, _ := fm["id"].(string)
		ticketID := rawID
		if ticketID == "" {
			ticketID = stem
		}

		// Required fields
		for _, field := range requiredFields {
			v, exists := fm[field]
			if !exists || v == nil || fmt.Sprintf("%v", v) == "" {
				add(ticketID, fmt.Sprintf("required field %q is missing or empty", field))
			}
		}

		// Field value checks
		if status, ok := fm["status"].(string); ok && status != "" {
			if !validStatuses[status] {
				add(ticketID, fmt.Sprintf("status %q is not valid; must be one of: open, in_progress, closed", status))
			}
		}
		if typ, ok := fm["type"].(string); ok && typ != "" {
			if !validTypes[typ] {
				add(ticketID, fmt.Sprintf("type %q is not valid; must be one of: bug, feature, task, epic, chore", typ))
			}
		}
		if priority, ok := fm["priority"].(int); ok {
			if priority < 0 || priority > 3 {
				add(ticketID, fmt.Sprintf("priority %d is not valid; must be an integer 0–3", priority))
			}
		}
		if mode, ok := fm["mode"].(string); ok && mode != "" {
			if !validModes[mode] {
				add(ticketID, fmt.Sprintf("mode %q is not valid; must be one of: afk, hitl", mode))
			}
		}

		// ID integrity: frontmatter id must match filename stem
		if rawID != "" && rawID != stem {
			add(ticketID, fmt.Sprintf("id %q does not match filename %q", rawID, stem+".md"))
		}

		// ID format
		if rawID != "" && !idPattern.MatchString(rawID) {
			add(ticketID, fmt.Sprintf("id %q does not match expected format %s-<n>", rawID, cfg.Prefix))
		}

		// File location
		status, _ := fm["status"].(string)
		if status == "closed" && !e.inArchive {
			add(ticketID, "closed ticket must be in .tickets/archive/, not .tickets/")
		}
		if status != "closed" && status != "" && e.inArchive {
			add(ticketID, fmt.Sprintf("ticket with status %q must be in .tickets/, not .tickets/archive/", status))
		}

		// Dependencies
		deps := stringSlice(fm["dependencies"])
		for _, depID := range deps {
			if _, found := allFMs[depID]; !found {
				add(ticketID, fmt.Sprintf("dependency %s not found in .tickets/ or .tickets/archive/", depID))
			}
		}

		// Dependency cycles: for each dep, can it reach back to this ticket?
		if rawID != "" {
			for _, depID := range deps {
				if path := findCycle(depID, rawID, ticketsDir, archiveDir); path != nil {
					add(ticketID, fmt.Sprintf("dependency cycle detected: %s → %s", rawID, joinPath(path)))
					break
				}
			}
		}

		// Links: referenced IDs must exist and must be symmetric
		links := stringSlice(fm["links"])
		for _, linkID := range links {
			if _, found := allFMs[linkID]; !found {
				add(ticketID, fmt.Sprintf("link %s not found in .tickets/ or .tickets/archive/", linkID))
				continue
			}
			// Symmetry check: linkID must list rawID in its links
			otherLinks := stringSlice(allFMs[linkID]["links"])
			symmetric := false
			for _, ol := range otherLinks {
				if ol == rawID {
					symmetric = true
					break
				}
			}
			if !symmetric {
				add(ticketID, fmt.Sprintf("link %s is not symmetric (%s does not link back)", linkID, linkID))
			}
		}
	}

	return violations, nil
}
