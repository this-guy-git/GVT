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

var createBranch bool

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch <branch>",
	Short: "Switch to another branch",
	Long:  `Switches the current working branch by updating HEAD. Use -c to create a new branch if it doesn't exist.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		branch := args[0]
		headFile := filepath.Join(".gvt", "HEAD")
		branchPath := filepath.Join(".gvt", "refs", "heads", branch)

		if _, err := os.Stat(".gvt"); os.IsNotExist(err) {
			fmt.Println("Not a GVT repository (no .gvt directory found).")
			return
		}

		// Create branch if it doesn't exist and --create is passed
		if _, err := os.Stat(branchPath); os.IsNotExist(err) {
			if createBranch {
				// Create refs and commits dir for branch
				os.MkdirAll(filepath.Dir(branchPath), 0755)
				os.WriteFile(branchPath, []byte(""), 0644)
				os.MkdirAll(filepath.Join(".gvt", "commits", branch), 0755)
				fmt.Printf("Created new branch '%s'.\n", branch)
			} else {
				fmt.Printf("Branch '%s' does not exist. Use -c to create it.\n", branch)
				return
			}
		}

		// Update HEAD
		os.WriteFile(headFile, []byte(fmt.Sprintf("ref: refs/heads/%s\n", branch)), 0644)
		fmt.Printf("Switched to branch '%s'.\n", branch)
	},
}

// GetCurrentBranch reads HEAD and returns the active branch
func GetCurrentBranch() string {
	headFile := filepath.Join(".gvt", "HEAD")
	data, err := os.ReadFile(headFile)
	if err != nil {
		return "main"
	}
	text := strings.TrimSpace(string(data))
	if strings.HasPrefix(text, "ref:") {
		return filepath.Base(strings.TrimPrefix(text, "ref: refs/heads/"))
	}
	return "main"
}

func init() {
	switchCmd.Flags().BoolVarP(&createBranch, "create", "c", false, "Create branch if it doesn't exist")
	rootCmd.AddCommand(switchCmd)
}
