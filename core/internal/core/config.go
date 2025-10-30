package core

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	UserName  string `json:"user_name,omitempty"`
	UserEmail string `json:"user_email,omitempty"`
}

// default: global
func GetConfigPath(local bool) string {
	if local {
		return filepath.Join(".gvt", "config.json")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".gvt", "config.json") // fallback
	}
	return filepath.Join(home, ".gvt", "config.json")
}

func LoadConfig(local bool) (*Config, error) {
	path := GetConfigPath(local)
	config := &Config{}

	// If file doesn't exist, return empty config (global config will auto-create)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return config, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(data, config)
	return config, nil
}

func SaveConfig(config *Config, local bool) error {
	path := GetConfigPath(local)
	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755) // create folder if missing
	}

	data, _ := json.MarshalIndent(config, "", "  ")
	return os.WriteFile(path, data, 0644)
}
