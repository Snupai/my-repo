package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.CNGTPath == "" {
		t.Error("CNGTPath should not be empty")
	}

	if cfg.DataDir == "" {
		t.Error("DataDir should not be empty")
	}

	// Check that paths are absolute
	if !filepath.IsAbs(cfg.CNGTPath) {
		t.Error("CNGTPath should be absolute")
	}

	if !filepath.IsAbs(cfg.DataDir) {
		t.Error("DataDir should be absolute")
	}

	// Check that data directory exists
	if _, err := os.Stat(cfg.DataDir); os.IsNotExist(err) {
		t.Error("DataDir should exist after Load()")
	}
}