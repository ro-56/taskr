package tickets

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var ErrAlreadyInitialized = errors.New("already initialized")

type Config struct {
	Prefix string `json:"prefix"`
}

func Init(dir, prefix string) error {
	ticketsDir := filepath.Join(dir, ".tickets")
	archiveDir := filepath.Join(ticketsDir, "archive")

	if _, err := os.Stat(filepath.Join(ticketsDir, "config.json")); err == nil {
		return ErrAlreadyInitialized
	}

	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return err
	}

	if prefix == "" {
		prefix = DerivePrefix(filepath.Base(dir))
	}

	cfg := Config{Prefix: prefix}
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(ticketsDir, "config.json"), data, 0644)
}

var nonAlnum = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func DerivePrefix(dirName string) string {
	slug := nonAlnum.ReplaceAllString(dirName, "")
	return strings.ToUpper(slug)
}
