package config

import (
	"os"
	"path/filepath"
)

func Dir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(configDir, "ouractl")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}

	return dir, nil
}

func Path() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}
