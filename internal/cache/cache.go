package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	CacheDir          = ".cache/brew-discover"
	FormulaeFile      = "formulae.json"
	CasksFile         = "casks.json"
	FormulaAnalytics  = "formula_analytics_%s.json"
	CaskAnalytics     = "cask_analytics_%s.json"
	FormulaTTL        = 24 * time.Hour
	AnalyticsTTL      = 6 * time.Hour
)

type CacheEntry struct {
	Data      json.RawMessage `json:"data"`
	Timestamp time.Time       `json:"timestamp"`
}

type Cache struct {
	dir string
}

func New() (*Cache, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	cacheDir := filepath.Join(homeDir, CacheDir)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &Cache{dir: cacheDir}, nil
}

func (c *Cache) Get(key string, ttl time.Duration) ([]byte, bool) {
	path := filepath.Join(c.dir, key)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}

	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false
	}

	if time.Since(entry.Timestamp) > ttl {
		return nil, false
	}

	return entry.Data, true
}

func (c *Cache) Set(key string, data []byte) error {
	path := filepath.Join(c.dir, key)

	entry := CacheEntry{
		Data:      data,
		Timestamp: time.Now(),
	}

	entryData, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal cache entry: %w", err)
	}

	if err := os.WriteFile(path, entryData, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}

func (c *Cache) Clear() error {
	entries, err := os.ReadDir(c.dir)
	if err != nil {
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	for _, entry := range entries {
		path := filepath.Join(c.dir, entry.Name())
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("failed to remove cache file %s: %w", entry.Name(), err)
		}
	}

	return nil
}

func (c *Cache) Dir() string {
	return c.dir
}
