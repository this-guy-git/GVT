/*
Copyright © 2025 this guy Labs <thisguy@thisguylabs.com>

This file is part of GVT (Guy's Versioning Tool).

GVT is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

GVT is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with GVT. If not, see <https://www.gnu.org/licenses/>.

Do not remove or modify this notice.
*/

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
