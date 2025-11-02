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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var branchCmd = &cobra.Command{
	Use:   "branch [name]",
	Short: "List or create branches",
	Long: `Manage branches in your GVT repository.
Running without arguments lists all branches.
Running with a name creates a new branch at the current commit.`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(".gvt"); os.IsNotExist(err) {
			fmt.Println("Not a GVT repository (no .gvt directory found)")
			return
		}

		refsDir := filepath.Join(".gvt", "refs", "heads")
		os.MkdirAll(refsDir, 0755)

		currentBranch := getCurrentBranch()
		headCommit := getBranchHead(currentBranch)

		// If no args → list branches
		if len(args) == 0 {
			files, err := os.ReadDir(refsDir)
			if err != nil {
				fmt.Println("Error reading branches:", err)
				return
			}

			if len(files) == 0 {
				fmt.Println("No branches found.")
				return
			}

			fmt.Println("Branches:")
			for _, f := range files {
				name := f.Name()
				prefix := "  "
				if name == currentBranch {
					prefix = "* "
				}
				fmt.Printf("%s%s\n", prefix, name)
			}
			return
		}

		// Creating a branch
		newBranch := args[0]
		refPath := filepath.Join(refsDir, newBranch)

		if _, err := os.Stat(refPath); err == nil {
			fmt.Printf("Branch '%s' already exists.\n", newBranch)
			return
		}

		if headCommit == "" {
			fmt.Println("No commits yet to branch from.")
			return
		}

		err := os.WriteFile(refPath, []byte(headCommit), 0644)
		if err != nil {
			fmt.Println("Failed to create branch:", err)
			return
		}

		fmt.Printf("Created branch '%s' at commit %s\n", newBranch, headCommit)
	},
}

func init() {
	rootCmd.AddCommand(branchCmd)
}

// getCurrentBranch and getBranchHead reused from commit.go
func getCurrentBranch() string {
	data, err := os.ReadFile(".gvt/HEAD")
	if err != nil {
		return ""
	}
	line := strings.TrimSpace(string(data))
	if strings.HasPrefix(line, "ref: ") {
		return filepath.Base(line[5:])
	}
	return ""
}

func getBranchHead(branch string) string {
	refPath := filepath.Join(".gvt", "refs", "heads", branch)
	data, err := os.ReadFile(refPath)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}
