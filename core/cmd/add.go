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

	"github.com/spf13/cobra"
	"github.com/this-guy-git/GVT/core/internal/utils"
)

var addCmd = &cobra.Command{
	Use:   "add [files or directories]",
	Short: "Stage files for the next commit",
	Long:  `Stages one or more files or directories to be committed in the next snapshot.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("No files specified. Usage: gvt add <file1> <file2> ... or gvt add .")
			return
		}

		// Ensure repo exists
		if _, err := os.Stat(".gvt"); os.IsNotExist(err) {
			fmt.Println("Not a GVT repository (no .gvt directory found)")
			return
		}

		stageFile := filepath.Join(".gvt", "stage.json")

		// Load existing staged files
		var staged []string
		if data, err := os.ReadFile(stageFile); err == nil {
			json.Unmarshal(data, &staged)
		}
		stagedSet := make(map[string]bool)
		for _, s := range staged {
			stagedSet[s] = true
		}

		// Load ignore patterns
		utils.LoadGvtIgnore()

		var filesToAdd []string

		for _, path := range args {
			info, err := os.Stat(path)
			if err != nil {
				fmt.Printf("File not found: %s\n", path)
				continue
			}

			if info.IsDir() {
				filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
					if err != nil {
						return nil
					}
					if info.IsDir() {
						if info.Name() == ".gvt" {
							return filepath.SkipDir
						}
						relDir, _ := filepath.Rel(".", p)
						if utils.IsIgnored(relDir + "/") {
							return filepath.SkipDir
						}
						return nil
					}

					relFile, _ := filepath.Rel(".", p)
					if utils.IsIgnored(relFile) {
						return nil
					}

					filesToAdd = append(filesToAdd, relFile)
					return nil
				})
			} else {
				relFile, _ := filepath.Rel(".", path)
				if !utils.IsIgnored(relFile) {
					filesToAdd = append(filesToAdd, relFile)
				}
			}
		}

		for _, rel := range filesToAdd {
			if !stagedSet[rel] {
				staged = append(staged, rel)
				stagedSet[rel] = true
				fmt.Printf("Added %s\n", rel)
			}
		}

		data, _ := json.MarshalIndent(staged, "", "  ")
		if err := os.WriteFile(stageFile, data, 0644); err != nil {
			fmt.Printf("Error saving stage file: %v\n", err)
			return
		}

		if len(filesToAdd) == 0 {
			fmt.Println("No files were added.")
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
