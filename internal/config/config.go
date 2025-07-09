package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	CNGTPath string
	DataDir  string
}

func Load() (*Config, error) {
	dataDir, err := getDataDir()
	if err != nil {
		return nil, err
	}

	return &Config{
		CNGTPath: filepath.Join(dataDir, "cngt"),
		DataDir:  dataDir,
	}, nil
}

func getDataDir() (string, error) {
	var dataDir string
	
	switch runtime.GOOS {
	case "windows":
		dataDir = os.Getenv("APPDATA")
		if dataDir == "" {
			dataDir = os.Getenv("USERPROFILE")
		}
	case "darwin":
		dataDir = os.Getenv("HOME")
		if dataDir != "" {
			dataDir = filepath.Join(dataDir, "Library", "Application Support")
		}
	default: // linux and others
		dataDir = os.Getenv("XDG_DATA_HOME")
		if dataDir == "" {
			home := os.Getenv("HOME")
			if home != "" {
				dataDir = filepath.Join(home, ".local", "share")
			}
		}
	}

	if dataDir == "" {
		return "", fmt.Errorf("unable to determine data directory")
	}

	dataDir = filepath.Join(dataDir, "cngt-cli")
	
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create data directory: %w", err)
	}

	return dataDir, nil
}