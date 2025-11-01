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

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone <source> [directory]",
	Short: "Clone a GVT repository",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		source := args[0]
		target := ""
		if len(args) >= 2 {
			target = args[1]
		} else {
			target = filepath.Base(source)
		}

		// Ensure target dir does not exist
		if _, err := os.Stat(target); !os.IsNotExist(err) {
			fmt.Printf("Directory %s already exists.\n", target)
			return
		}

		// Load ignore patterns
		ignoreList := loadGvtIgnore(filepath.Join(source, ".gvtignore"))

		// Collect all files to copy first (for progress bar)
		var files []string
		filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}
			rel, _ := filepath.Rel(source, path)
			for _, pattern := range ignoreList {
				match, _ := filepath.Match(pattern, rel)
				if match || strings.HasPrefix(rel, pattern) {
					return nil
				}
			}
			files = append(files, path)
			return nil
		})

		bar := progressbar.NewOptions(len(files),
			progressbar.OptionSetDescription("Cloning..."),
			progressbar.OptionShowCount(),
			progressbar.OptionShowIts(),
			progressbar.OptionSetPredictTime(false),
			progressbar.OptionSetWidth(40),
			progressbar.OptionClearOnFinish(),
		)

		// Copy files with progress
		for _, file := range files {
			bar.Add(1)
			rel, _ := filepath.Rel(source, file)
			targetPath := filepath.Join(target, rel)

			os.MkdirAll(filepath.Dir(targetPath), 0755)

			data, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			os.WriteFile(targetPath, data, 0644)
		}

		bar.Finish()
		fmt.Printf("Cloned %s to %s\n", source, target)
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
}

func loadGvtIgnore(path string) []string {
	var ignoreList []string

	file, err := os.Open(path)
	if err != nil {
		return ignoreList // no .gvtignore, ignore nothing
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		ignoreList = append(ignoreList, line)
	}
	return ignoreList
}
