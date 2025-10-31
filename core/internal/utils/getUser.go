package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// GetCommitUser returns the user.name from local config first (if set), then global, else "unknown"
func GetCommitUser(repoRoot string) string {
	// 1. Local config
	localCfg := filepath.Join(repoRoot, ".gvt", "config.json")
	if data, err := os.ReadFile(localCfg); err == nil {
		var cfg map[string]string
		if json.Unmarshal(data, &cfg) == nil {
			if u, ok := cfg["user_name"]; ok && u != "" {
				return u
			}
		}
	}

	// 2. Global config
	home, err := os.UserHomeDir()
	if err == nil {
		globalCfg := filepath.Join(home, ".gvt", "config.json")
		if data, err := os.ReadFile(globalCfg); err == nil {
			var cfg map[string]string
			if json.Unmarshal(data, &cfg) == nil {
				if u, ok := cfg["user_name"]; ok && u != "" {
					return u
				}
			}
		}
	}

	return "unknown"
}
