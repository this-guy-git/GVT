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
