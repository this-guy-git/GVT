/*
Copyright © 2025 this guy Labs <thisguy@thisguylabs.com>
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
