package tickets

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func List(dir string) ([]TicketSummary, error) {
	ticketsDir := filepath.Join(dir, ".tickets")
	summaries, _, err := loadAll(ticketsDir)
	if err != nil {
		return []TicketSummary{}, nil
	}
	if summaries == nil {
		return []TicketSummary{}, nil
	}
	return summaries, nil
}

func ListByStatus(dir, status string) ([]TicketSummary, error) {
	ticketsDir := filepath.Join(dir, ".tickets")
	archiveDir := filepath.Join(ticketsDir, "archive")

	dirs := []string{ticketsDir}
	if status == "closed" {
		dirs = append(dirs, archiveDir)
	}

	var out []TicketSummary
	for _, d := range dirs {
		summaries, _, err := loadAll(d)
		if err != nil {
			continue
		}
		for _, s := range summaries {
			if s.Status == status {
				out = append(out, s)
			}
		}
	}
	if out == nil {
		return []TicketSummary{}, nil
	}
	return out, nil
}

func ListTags(dir string) ([]string, error) {
	ticketsDir := filepath.Join(dir, ".tickets")
	entries, err := os.ReadDir(ticketsDir)
	if err != nil {
		return []string{}, nil
	}

	seen := make(map[string]struct{})
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
		for _, tag := range stringSlice(fm["tags"]) {
			seen[tag] = struct{}{}
		}
	}

	tags := make([]string, 0, len(seen))
	for t := range seen {
		tags = append(tags, t)
	}
	sort.Strings(tags)
	return tags, nil
}

func ListCounts(dir string) (map[string]int, error) {
	ticketsDir := filepath.Join(dir, ".tickets")
	archiveDir := filepath.Join(ticketsDir, "archive")

	counts := map[string]int{"open": 0, "in_progress": 0, "closed": 0}

	for _, d := range []string{ticketsDir, archiveDir} {
		summaries, _, err := loadAll(d)
		if err != nil {
			continue
		}
		for _, s := range summaries {
			if _, ok := counts[s.Status]; ok {
				counts[s.Status]++
			}
		}
	}
	return counts, nil
}
