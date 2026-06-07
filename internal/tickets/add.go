package tickets

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var ErrNotInitialized = errors.New("not initialized: run taskr init first")

type AddOptions struct {
	Title    string
	Type     string
	Priority *int // nil = use default (2)
	Mode     string
	Tags     []string
	Body     string
}

func Add(dir string, opts AddOptions) (string, error) {
	ticketsDir := filepath.Join(dir, ".tickets")
	if _, err := os.Stat(filepath.Join(ticketsDir, "config.json")); err != nil {
		return "", ErrNotInitialized
	}

	cfg, err := loadConfig(ticketsDir)
	if err != nil {
		return "", err
	}

	id := cfg.Prefix + "-" + strconv.Itoa(cfg.NextID)

	cfg.NextID++
	if err := saveConfig(ticketsDir, cfg); err != nil {
		return "", err
	}

	now := time.Now().UTC()
	ticketType := opts.Type
	if ticketType == "" {
		ticketType = "task"
	}
	priority := 2
	if opts.Priority != nil {
		priority = *opts.Priority
	}
	mode := opts.Mode
	if mode == "" {
		mode = "hitl"
	}
	tags := opts.Tags
	if tags == nil {
		tags = []string{}
	}

	frontmatter := map[string]any{
		"id":           id,
		"title":        opts.Title,
		"status":       "open",
		"type":         ticketType,
		"priority":     priority,
		"mode":         mode,
		"created":      now.Format(time.RFC3339),
		"updated":      now.Format(time.RFC3339),
		"dependencies": []string{},
		"links":        []string{},
		"tags":         tags,
	}

	ticketPath := filepath.Join(ticketsDir, id+".md")
	if err := writeTicket(ticketPath, frontmatter, opts.Body); err != nil {
		return "", err
	}

	return id, nil
}
