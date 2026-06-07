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
	NextID int    `json:"next_id"`
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

	return saveConfig(ticketsDir, Config{Prefix: prefix, NextID: 1})
}

func loadConfig(ticketsDir string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(filepath.Join(ticketsDir, "config.json"))
	if err != nil {
		return cfg, err
	}
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}

func saveConfig(ticketsDir string, cfg Config) error {
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
