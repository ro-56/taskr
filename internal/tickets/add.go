package tickets

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

var ErrNotInitialized = errors.New("not initialized: run taskr init first")

type AddOptions struct {
	Title    string
	Type     string
	Priority *int // nil = use default (2)
	Mode     string
	Tags     []string
}

func Add(dir string, opts AddOptions) (string, error) {
	ticketsDir := filepath.Join(dir, ".tickets")
	if _, err := os.Stat(filepath.Join(ticketsDir, "config.json")); err != nil {
		return "", ErrNotInitialized
	}

	data, err := os.ReadFile(filepath.Join(ticketsDir, "config.json"))
	if err != nil {
		return "", err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return "", err
	}

	var b [4]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	id := cfg.Prefix + "-" + hex.EncodeToString(b[:])

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

	fm, err := yaml.Marshal(frontmatter)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "---\n%s---\n", fm)

	ticketPath := filepath.Join(ticketsDir, id+".md")
	if err := os.WriteFile(ticketPath, buf.Bytes(), 0644); err != nil {
		return "", err
	}

	return id, nil
}
