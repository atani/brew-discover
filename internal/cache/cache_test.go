package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCache_SetAndGet(t *testing.T) {
	// Create temp directory for test
	tmpDir, err := os.MkdirTemp("", "cache-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cache := &Cache{dir: tmpDir}

	// Test Set
	testData := []byte(`{"test":"data"}`)
	key := "test.json"

	err = cache.Set(key, testData)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Test Get with valid TTL
	result, ok := cache.Get(key, time.Hour)
	if !ok {
		t.Fatal("Get returned false for existing cache")
	}

	if string(result) != string(testData) {
		t.Errorf("got %s, want %s", result, testData)
	}
}

func TestCache_Get_Expired(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cache-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cache := &Cache{dir: tmpDir}

	testData := []byte(`{"test": "data"}`)
	key := "test.json"

	err = cache.Set(key, testData)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Test Get with expired TTL (0 duration = immediately expired)
	_, ok := cache.Get(key, 0)
	if ok {
		t.Fatal("Get returned true for expired cache")
	}
}

func TestCache_Get_NotFound(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cache-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cache := &Cache{dir: tmpDir}

	_, ok := cache.Get("nonexistent.json", time.Hour)
	if ok {
		t.Fatal("Get returned true for non-existent cache")
	}
}

func TestCache_Clear(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cache-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cache := &Cache{dir: tmpDir}

	// Create some cache files
	_ = cache.Set("test1.json", []byte(`{"test": 1}`))
	_ = cache.Set("test2.json", []byte(`{"test": 2}`))

	err = cache.Clear()
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	// Verify files are deleted
	entries, _ := os.ReadDir(tmpDir)
	if len(entries) != 0 {
		t.Errorf("expected 0 files after clear, got %d", len(entries))
	}
}

func TestCache_Dir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cache-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cache := &Cache{dir: tmpDir}

	if cache.Dir() != tmpDir {
		t.Errorf("Dir() = %s, want %s", cache.Dir(), tmpDir)
	}
}

func TestNew(t *testing.T) {
	cache, err := New()
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Verify cache directory exists
	homeDir, _ := os.UserHomeDir()
	expectedDir := filepath.Join(homeDir, CacheDir)

	if cache.Dir() != expectedDir {
		t.Errorf("Dir() = %s, want %s", cache.Dir(), expectedDir)
	}
}
