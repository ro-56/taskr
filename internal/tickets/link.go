package tickets

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrSelfLink      = errors.New("a ticket cannot depend on itself")
	ErrAlreadyLinked = errors.New("dependency already exists")
)

type ErrCycle struct {
	Path []string // cycle path, e.g. [B, A, B]
}

func (e *ErrCycle) Error() string {
	return "dependency cycle detected: " + joinPath(e.Path)
}

func joinPath(p []string) string {
	out := ""
	for i, s := range p {
		if i > 0 {
			out += " → "
		}
		out += s
	}
	return out
}

func Link(dir, dependentPartialID, dependsOnPartialID string) error {
	ticketsDir := filepath.Join(dir, ".tickets")
	archiveDir := filepath.Join(ticketsDir, "archive")

	depPath, err := resolveID(dependentPartialID, ticketsDir, archiveDir)
	if err != nil {
		return err
	}
	onPath, err := resolveID(dependsOnPartialID, ticketsDir, archiveDir)
	if err != nil {
		return err
	}

	depRaw, err := os.ReadFile(depPath)
	if err != nil {
		return err
	}
	onRaw, err := os.ReadFile(onPath)
	if err != nil {
		return err
	}

	depFM, depBody, err := parseFrontmatter(depRaw)
	if err != nil {
		return err
	}
	onFM, _, err := parseFrontmatter(onRaw)
	if err != nil {
		return err
	}

	dependentID, _ := depFM["id"].(string)
	dependsOnID, _ := onFM["id"].(string)

	if dependentID == dependsOnID {
		return ErrSelfLink
	}

	deps := stringSlice(depFM["dependencies"])
	for _, d := range deps {
		if d == dependsOnID {
			return ErrAlreadyLinked
		}
	}

	// Cycle check: can we reach dependentID starting from dependsOnID?
	if path := findCycle(dependsOnID, dependentID, ticketsDir, archiveDir); path != nil {
		return &ErrCycle{Path: append([]string{dependentID}, path...)}
	}

	deps = append(deps, dependsOnID)
	depFM["dependencies"] = deps
	depFM["updated"] = time.Now().UTC().Format(time.RFC3339)
	return writeTicket(depPath, depFM, depBody)
}

func Prune(dir, partialID string) error {
	ticketsDir := filepath.Join(dir, ".tickets")
	archiveDir := filepath.Join(ticketsDir, "archive")

	targetPath, err := resolveID(partialID, ticketsDir, archiveDir)
	if err != nil {
		return err
	}

	raw, err := os.ReadFile(targetPath)
	if err != nil {
		return err
	}
	fm, body, err := parseFrontmatter(raw)
	if err != nil {
		return err
	}
	targetID, _ := fm["id"].(string)

	// Clear own dependencies
	fm["dependencies"] = []string{}
	fm["updated"] = time.Now().UTC().Format(time.RFC3339)
	if err := writeTicket(targetPath, fm, body); err != nil {
		return err
	}

	// Remove targetID from all other tickets' dependencies
	for _, scanDir := range []string{ticketsDir, archiveDir} {
		entries, err := os.ReadDir(scanDir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			p := filepath.Join(scanDir, e.Name())
			eRaw, err := os.ReadFile(p)
			if err != nil {
				continue
			}
			eFM, eBody, err := parseFrontmatter(eRaw)
			if err != nil {
				continue
			}
			deps := stringSlice(eFM["dependencies"])
			filtered := deps[:0]
			changed := false
			for _, d := range deps {
				if d == targetID {
					changed = true
					continue
				}
				filtered = append(filtered, d)
			}
			if changed {
				eFM["dependencies"] = filtered
				eFM["updated"] = time.Now().UTC().Format(time.RFC3339)
				writeTicket(p, eFM, eBody)
			}
		}
	}
	return nil
}

func Unlink(dir, dependentPartialID, dependsOnPartialID string) error {
	ticketsDir := filepath.Join(dir, ".tickets")
	archiveDir := filepath.Join(ticketsDir, "archive")

	depPath, err := resolveID(dependentPartialID, ticketsDir, archiveDir)
	if err != nil {
		return err
	}
	onPath, err := resolveID(dependsOnPartialID, ticketsDir, archiveDir)
	if err != nil {
		return err
	}

	depRaw, err := os.ReadFile(depPath)
	if err != nil {
		return err
	}
	onRaw, err := os.ReadFile(onPath)
	if err != nil {
		return err
	}

	depFM, depBody, err := parseFrontmatter(depRaw)
	if err != nil {
		return err
	}
	onFM, _, err := parseFrontmatter(onRaw)
	if err != nil {
		return err
	}

	dependsOnID, _ := onFM["id"].(string)

	deps := stringSlice(depFM["dependencies"])
	filtered := deps[:0]
	for _, d := range deps {
		if d != dependsOnID {
			filtered = append(filtered, d)
		}
	}
	depFM["dependencies"] = filtered
	depFM["updated"] = time.Now().UTC().Format(time.RFC3339)
	return writeTicket(depPath, depFM, depBody)
}

// findCycle does a BFS from startID looking for targetID in the dependency graph.
// Returns the path from startID to targetID (inclusive) if found, nil otherwise.
func findCycle(startID, targetID, ticketsDir, archiveDir string) []string {
	type node struct {
		id   string
		path []string
	}
	visited := map[string]bool{}
	queue := []node{{id: startID, path: []string{startID}}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if visited[cur.id] {
			continue
		}
		visited[cur.id] = true

		if cur.id == targetID {
			return cur.path
		}

		deps := depsOf(cur.id, ticketsDir, archiveDir)
		for _, d := range deps {
			if !visited[d] {
				newPath := make([]string, len(cur.path)+1)
				copy(newPath, cur.path)
				newPath[len(cur.path)] = d
				queue = append(queue, node{id: d, path: newPath})
			}
		}
	}
	return nil
}

// depsOf returns the dependencies list of a ticket by full ID.
func depsOf(id, ticketsDir, archiveDir string) []string {
	for _, dir := range []string{ticketsDir, archiveDir} {
		raw, err := os.ReadFile(filepath.Join(dir, id+".md"))
		if err != nil {
			continue
		}
		fm, _, err := parseFrontmatter(raw)
		if err != nil {
			continue
		}
		return stringSlice(fm["dependencies"])
	}
	return nil
}

// stringSlice coerces a frontmatter []any to []string.
func stringSlice(v any) []string {
	raw, _ := v.([]any)
	out := make([]string, 0, len(raw))
	for _, r := range raw {
		if s, ok := r.(string); ok {
			out = append(out, s)
		}
	}
	return out
}
