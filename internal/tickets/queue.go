package tickets

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type TicketSummary struct {
	ID       string
	Title    string
	Priority int
	Status   string
	Mode     string
	Updated  time.Time
}

func Ready(dir, mode string) ([]TicketSummary, error) {
	ticketsDir := filepath.Join(dir, ".tickets")
	archiveDir := filepath.Join(ticketsDir, "archive")

	summaries, fms, err := loadAll(ticketsDir)
	if err != nil {
		return []TicketSummary{}, nil
	}

	var result []TicketSummary
	for i, s := range summaries {
		if s.Status == "closed" {
			continue
		}
		if mode != "" && s.Mode != mode {
			continue
		}
		if isBlocked(fms[i], ticketsDir, archiveDir) {
			continue
		}
		result = append(result, s)
	}

	sortSummaries(result)
	if result == nil {
		return []TicketSummary{}, nil
	}
	return result, nil
}

func Blocked(dir string) ([]TicketSummary, error) {
	ticketsDir := filepath.Join(dir, ".tickets")
	archiveDir := filepath.Join(ticketsDir, "archive")

	summaries, fms, err := loadAll(ticketsDir)
	if err != nil {
		return []TicketSummary{}, nil
	}

	var result []TicketSummary
	for i, s := range summaries {
		if s.Status == "closed" {
			continue
		}
		if isBlocked(fms[i], ticketsDir, archiveDir) {
			result = append(result, s)
		}
	}

	sortSummaries(result)
	if result == nil {
		return []TicketSummary{}, nil
	}
	return result, nil
}

func loadAll(ticketsDir string) ([]TicketSummary, []map[string]any, error) {
	var summaries []TicketSummary
	var fms []map[string]any

	entries, err := os.ReadDir(ticketsDir)
	if err != nil {
		return nil, nil, err
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(ticketsDir, e.Name()))
		if err != nil {
			continue
		}
		fm, _, err := parseFrontmatter(raw)
		if err != nil {
			continue
		}
		s := summaryFrom(fm)
		summaries = append(summaries, s)
		fms = append(fms, fm)
	}
	return summaries, fms, nil
}

func summaryFrom(fm map[string]any) TicketSummary {
	id, _ := fm["id"].(string)
	title, _ := fm["title"].(string)
	status, _ := fm["status"].(string)
	mode, _ := fm["mode"].(string)

	priority := 2
	if p, ok := fm["priority"].(int); ok {
		priority = p
	}

	var updated time.Time
	if u, ok := fm["updated"].(string); ok {
		updated, _ = time.Parse(time.RFC3339, u)
	}

	return TicketSummary{
		ID:       id,
		Title:    title,
		Priority: priority,
		Status:   status,
		Mode:     mode,
		Updated:  updated,
	}
}

func isBlocked(fm map[string]any, ticketsDir, archiveDir string) bool {
	for _, depID := range stringSlice(fm["dependencies"]) {
		s := statusOf(depID, ticketsDir, archiveDir)
		if s == "open" || s == "in_progress" {
			return true
		}
	}
	return false
}

func sortSummaries(ss []TicketSummary) {
	sort.Slice(ss, func(i, j int) bool {
		if ss[i].Priority != ss[j].Priority {
			return ss[i].Priority < ss[j].Priority
		}
		return ss[i].Updated.After(ss[j].Updated)
	})
}
