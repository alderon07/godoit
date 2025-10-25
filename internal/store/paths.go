package store

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetDataDir returns the platform-specific data directory for godo
func GetDataDir() (string, error) {
	var baseDir string

	switch runtime.GOOS {
	case "windows":
		// Windows: %APPDATA%/godo
		baseDir = os.Getenv("APPDATA")
		if baseDir == "" {
			baseDir = os.Getenv("USERPROFILE")
		}
	case "darwin":
		// macOS: ~/Library/Application Support/godo
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		baseDir = filepath.Join(homeDir, "Library", "Application Support")
	default:
		// Linux and others: ~/.local/share/godo
		baseDir = os.Getenv("XDG_DATA_HOME")
		if baseDir == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			baseDir = filepath.Join(homeDir, ".local", "share")
		}
	}

	dataDir := filepath.Join(baseDir, "godo")

	// Ensure the directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return "", err
	}

	return dataDir, nil
}

// GetConfigDir returns the platform-specific config directory for godo
func GetConfigDir() (string, error) {
	var baseDir string

	switch runtime.GOOS {
	case "windows":
		// Windows: %APPDATA%/godo
		baseDir = os.Getenv("APPDATA")
		if baseDir == "" {
			baseDir = os.Getenv("USERPROFILE")
		}
	case "darwin":
		// macOS: ~/Library/Application Support/godo
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		baseDir = filepath.Join(homeDir, "Library", "Application Support")
	default:
		// Linux and others: ~/.config/godo
		baseDir = os.Getenv("XDG_CONFIG_HOME")
		if baseDir == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			baseDir = filepath.Join(homeDir, ".config")
		}
	}

	configDir := filepath.Join(baseDir, "godo")

	// Ensure the directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return configDir, nil
}

// GetDataFile returns the full path to the tasks data file
func GetDataFile() (string, error) {
	dataDir, err := GetDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dataDir, "tasks.json"), nil
}

