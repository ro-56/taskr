package tickets

import (
	"errors"
	"os"
	"path/filepath"
	"time"
)

var ErrInvalidPriority = errors.New("priority must be 0-3")
var ErrInvalidMode = errors.New("mode must be afk or hitl")

type UpdateOptions struct {
	Title    *string
	Priority *int
	Mode     *string
	Tags     []string
	Body     *string
}

func Update(dir, partialID string, opts UpdateOptions) error {
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

	if opts.Title != nil {
		fm["title"] = *opts.Title
	}

	if opts.Priority != nil {
		if *opts.Priority < 0 || *opts.Priority > 3 {
			return ErrInvalidPriority
		}
		fm["priority"] = *opts.Priority
	}

	if opts.Mode != nil {
		if *opts.Mode != "afk" && *opts.Mode != "hitl" {
			return ErrInvalidMode
		}
		fm["mode"] = *opts.Mode
	}

	if opts.Tags != nil {
		fm["tags"] = opts.Tags
	}

	if opts.Body != nil {
		body = *opts.Body
	}

	fm["updated"] = time.Now().UTC().Format(time.RFC3339)

	return writeTicket(path, fm, body)
}
