package auth

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/hagelstam/ouractl/internal/config"
)

type configFile struct {
	AccessToken string `json:"access_token"`
}

func SaveToken(token string) error {
	path, err := config.Path()
	if err != nil {
		return err
	}

	data, err := json.Marshal(configFile{AccessToken: token})
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o600)
}

func LoadToken() (string, error) {
	path, err := config.Path()
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", errors.New("not logged in, run `ouractl auth login`")
		}
		return "", err
	}

	var cfg configFile
	if err := json.Unmarshal(data, &cfg); err != nil {
		return "", err
	}
	if cfg.AccessToken == "" {
		return "", errors.New("not logged in, run `ouractl auth login`")
	}

	return cfg.AccessToken, nil
}

func RemoveToken() error {
	path, err := config.Path()
	if err != nil {
		return err
	}

	err = os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func IsLoggedIn() bool {
	_, err := LoadToken()
	return err == nil
}
