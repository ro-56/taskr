package tickets

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var ErrNotFound = errors.New("ticket not found")

type ErrAmbiguous struct {
	Matches []string // full file paths
}

func (e *ErrAmbiguous) Error() string {
	return "ambiguous partial ID"
}

func newErrAmbiguous(paths []string) error {
	ids := make([]string, len(paths))
	for i, p := range paths {
		ids[i] = strings.TrimSuffix(filepath.Base(p), ".md")
	}
	return &ErrAmbiguous{Matches: ids}
}

type DepRef struct {
	ID     string
	Status string
}

type ShowResult struct {
	ID          string
	Frontmatter map[string]any
	Body        string
	DependsOn   []DepRef
	RequiredBy  []DepRef
}

func Show(dir, partialID string) (*ShowResult, error) {
	ticketsDir := filepath.Join(dir, ".tickets")
	archiveDir := filepath.Join(ticketsDir, "archive")

	path, err := resolveID(partialID, ticketsDir, archiveDir)
	if err != nil {
		return nil, err
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	fm, body, err := parseFrontmatter(raw)
	if err != nil {
		return nil, err
	}

	id, _ := fm["id"].(string)

	dependsOn, err := resolveDependsOn(fm, ticketsDir, archiveDir)
	if err != nil {
		return nil, err
	}

	requiredBy, err := resolveRequiredBy(id, ticketsDir, archiveDir)
	if err != nil {
		return nil, err
	}

	return &ShowResult{
		ID:          id,
		Frontmatter: fm,
		Body:        body,
		DependsOn:   dependsOn,
		RequiredBy:  requiredBy,
	}, nil
}

// resolveRequiredBy scans all tickets and returns those that list id in their dependencies.
func resolveRequiredBy(id, ticketsDir, archiveDir string) ([]DepRef, error) {
	var refs []DepRef
	for _, dir := range []string{ticketsDir, archiveDir} {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
			if err != nil {
				continue
			}
			fm, _, err := parseFrontmatter(raw)
			if err != nil {
				continue
			}
			if listsDep(fm, id) {
				candidateID, _ := fm["id"].(string)
				status, _ := fm["status"].(string)
				refs = append(refs, DepRef{ID: candidateID, Status: status})
			}
		}
	}
	if refs == nil {
		refs = []DepRef{}
	}
	return refs, nil
}

// listsDep reports whether the ticket's dependencies field contains depID.
func listsDep(fm map[string]any, depID string) bool {
	rawDeps, _ := fm["dependencies"].([]any)
	for _, d := range rawDeps {
		if s, _ := d.(string); s == depID {
			return true
		}
	}
	return false
}

// resolveDependsOn reads the dependencies list from frontmatter and looks up each ticket's status.
func resolveDependsOn(fm map[string]any, ticketsDir, archiveDir string) ([]DepRef, error) {
	rawDeps, _ := fm["dependencies"].([]any)
	if len(rawDeps) == 0 {
		return []DepRef{}, nil
	}

	refs := make([]DepRef, 0, len(rawDeps))
	for _, d := range rawDeps {
		depID, _ := d.(string)
		if depID == "" {
			continue
		}
		status := statusOf(depID, ticketsDir, archiveDir)
		refs = append(refs, DepRef{ID: depID, Status: status})
	}
	return refs, nil
}

// statusOf looks up the status field of a ticket by full ID.
func statusOf(id, ticketsDir, archiveDir string) string {
	for _, dir := range []string{ticketsDir, archiveDir} {
		raw, err := os.ReadFile(filepath.Join(dir, id+".md"))
		if err != nil {
			continue
		}
		fm, _, err := parseFrontmatter(raw)
		if err != nil {
			continue
		}
		if s, ok := fm["status"].(string); ok {
			return s
		}
	}
	return "unknown"
}

var bareNumber = regexp.MustCompile(`^[0-9]+$`)

// resolveID finds the full path for a ticket ID. A bare number (e.g. "42")
// is expanded to "<prefix>-42" before lookup. Matching is exact — no
// prefix-matching against partial IDs.
func resolveID(id, ticketsDir, archiveDir string) (string, error) {
	if bareNumber.MatchString(id) {
		if cfg, err := loadConfig(ticketsDir); err == nil {
			id = cfg.Prefix + "-" + id
		}
	}

	var matches []string

	for _, dir := range []string{ticketsDir, archiveDir} {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			name := strings.TrimSuffix(e.Name(), ".md")
			if name == id {
				matches = append(matches, filepath.Join(dir, e.Name()))
			}
		}
	}

	switch len(matches) {
	case 0:
		return "", ErrNotFound
	case 1:
		return matches[0], nil
	default:
		return "", newErrAmbiguous(matches)
	}
}

// parseFrontmatter splits a markdown file into YAML frontmatter and body.
func parseFrontmatter(raw []byte) (map[string]any, string, error) {
	s := string(raw)
	if !strings.HasPrefix(s, "---\n") {
		return map[string]any{}, s, nil
	}
	end := strings.Index(s[4:], "\n---\n")
	if end < 0 {
		return map[string]any{}, s, nil
	}
	fmRaw := s[4 : end+4]
	body := strings.TrimPrefix(s[end+9:], "\n")

	var fm map[string]any
	if err := yaml.Unmarshal([]byte(fmRaw), &fm); err != nil {
		return nil, "", err
	}
	return fm, body, nil
}
