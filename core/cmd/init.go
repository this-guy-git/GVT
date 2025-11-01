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

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new GVT repository",
	Run: func(cmd *cobra.Command, args []string) {
		repoDir := ".gvt"
		dirs := []string{
			repoDir,
			filepath.Join(repoDir, "objects"),
			filepath.Join(repoDir, "refs"),
			filepath.Join(repoDir, "refs", "heads"),
		}

		if _, err := os.Stat(repoDir); err == nil {
			fmt.Println("GVT repository already exists.")
			return
		}

		for _, dir := range dirs {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Println("Error creating directory:", err)
				return
			}
		}

		// initialize HEAD
		headPath := filepath.Join(repoDir, "HEAD")
		if err := os.WriteFile(headPath, []byte("ref: refs/heads/main\n"), 0644); err != nil {
			fmt.Println("Error initializing HEAD:", err)
			return
		}

		configPath := filepath.Join(repoDir, "config.json")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			if err := os.WriteFile(configPath, []byte("{}"), 0644); err != nil {
				fmt.Println("Error creating config file:", err)
				return
			}
		}

		fmt.Println("Initialized empty GVT repository in", repoDir)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
