package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func getAppDataDir(appName string) (string, error) {
	var baseDir string

	switch runtime.GOOS {
	case "windows":
		baseDir = os.Getenv("APPDATA")
		if baseDir == "" {
			baseDir = os.Getenv("USERPROFILE") + "\\AppData\\Roaming"
		}

	case "darwin":
		baseDir = os.Getenv("HOME")
		if baseDir != "" {
			baseDir = filepath.Join(baseDir, "Library", "Application Support")
		}

	default:
		baseDir = os.Getenv("XDG_DATA_HOME")
		if baseDir == "" {
			home := os.Getenv("HOME")
			if home != "" {
				baseDir = filepath.Join(home, ".local", "share")
			}
		}
	}

	if baseDir == "" {
		return "", fmt.Errorf("could not determine application data directory")
	}

	appDir := filepath.Join(baseDir, appName)

	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}

	return appDir, nil
}
