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
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

var ignorePatterns []string

func LoadGvtIgnore() {
	ignorePatterns = []string{}
	file, err := os.Open(".gvtignore")
	if err != nil {
		return // no ignore file is fine
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		ignorePatterns = append(ignorePatterns, line)
	}
}

func IsIgnored(path string) bool {
	path = filepath.ToSlash(filepath.Clean(path)) // normalize slashes

	for _, pattern := range ignorePatterns {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" || strings.HasPrefix(pattern, "#") {
			continue
		}
		pattern = filepath.ToSlash(filepath.Clean(pattern))

		// Exact match or subpath for directories
		if strings.HasSuffix(pattern, "/") {
			if strings.HasPrefix(path+"/", pattern) {
				return true
			}
		} else {
			matched, _ := filepath.Match(pattern, filepath.Base(path))
			if matched {
				return true
			}
		}
	}
	return false
}
