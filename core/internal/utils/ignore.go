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
