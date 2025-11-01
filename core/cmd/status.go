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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/this-guy-git/GVT/core/internal/utils"
)

func init() {
	utils.LoadGvtIgnore()
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the current status of your working directory",
	Long:  `Displays which files are staged, modified, or untracked in the current GVT repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(".gvt"); os.IsNotExist(err) {
			fmt.Println("Not a GVT repository (no .gvt directory found)")
			return
		}

		stageFile := filepath.Join(".gvt", "stage.json")
		var staged []string
		if data, err := os.ReadFile(stageFile); err == nil {
			json.Unmarshal(data, &staged)
		}

		stagedSet := make(map[string]bool)
		for _, s := range staged {
			stagedSet[s] = true
		}

		// Gather all files recursively
		var allFiles []string
		filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				// skip .gvt directory
				if strings.HasPrefix(path, ".gvt") {
					return filepath.SkipDir
				}
				return nil
			}
			allFiles = append(allFiles, path)
			return nil
		})

		fmt.Println("On branch: main")
		fmt.Println()

		if len(staged) > 0 {
			fmt.Println("Staged files:")
			for _, f := range staged {
				fmt.Printf("  %s\n", relPath(f))
			}
			fmt.Println()
		}

		// check modified/untracked
		fmt.Println("Changes not staged for commit:")
		for _, f := range allFiles {
			abs, _ := filepath.Abs(f)
			if stagedSet[abs] {
				// Compare mod times to see if changed since staged
				info, _ := os.Stat(f)
				if time.Since(info.ModTime()) < 2*time.Second {
					continue // freshly staged probably
				}
				fmt.Printf("  modified: %s\n", f)
			}
		}

		fmt.Println()
		fmt.Println("Untracked files:")
		for _, f := range allFiles {
			abs, _ := filepath.Abs(f)
			if !stagedSet[abs] {
				fmt.Printf("  %s\n", f)
			}
		}

		fmt.Println()
	},
}

func relPath(abs string) string {
	wd, _ := os.Getwd()
	rel, err := filepath.Rel(wd, abs)
	if err != nil {
		return abs
	}
	return rel
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
