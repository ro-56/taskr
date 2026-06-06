package tickets

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

var ErrWrongStatus = errors.New("wrong status for this transition")

func Start(dir, partialID string) error {
	ticketsDir := filepath.Join(dir, ".tickets")
	archiveDir := filepath.Join(ticketsDir, "archive")

	path, err := resolveID(partialID, ticketsDir, archiveDir)
	if err != nil {
		return err
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	fm, body, err := parseFrontmatter(raw)
	if err != nil {
		return err
	}

	if fm["status"] == "closed" {
		return ErrWrongStatus
	}

	fm["status"] = "in_progress"
	fm["updated"] = time.Now().UTC().Format(time.RFC3339)

	return writeTicket(path, fm, body)
}

func Close(dir, partialID, summary string) error {
	ticketsDir := filepath.Join(dir, ".tickets")
	archiveDir := filepath.Join(ticketsDir, "archive")

	path, err := resolveID(partialID, ticketsDir, archiveDir)
	if err != nil {
		return err
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	fm, body, err := parseFrontmatter(raw)
	if err != nil {
		return err
	}

	if fm["status"] == "closed" {
		return ErrWrongStatus
	}

	fm["status"] = "closed"
	fm["updated"] = time.Now().UTC().Format(time.RFC3339)

	if summary != "" {
		body = appendNotes(body, summary)
	}

	id, _ := fm["id"].(string)
	archivePath := filepath.Join(archiveDir, id+".md")

	if err := writeTicket(archivePath, fm, body); err != nil {
		return err
	}
	return os.Remove(path)
}

func appendNotes(body, summary string) string {
	note := "## Notes\n\n" + summary + "\n"
	if body == "" {
		return note
	}
	return body + "\n" + note
}

func writeTicket(path string, fm map[string]any, body string) error {
	out, err := yaml.Marshal(fm)
	if err != nil {
		return err
	}
	content := "---\n" + string(out) + "---\n"
	if body != "" {
		content += "\n" + body
	}
	return os.WriteFile(path, []byte(content), 0644)
}
