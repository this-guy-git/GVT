package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [files or directories]",
	Short: "Unstage files so they are no longer tracked",
	Long:  `Removes files or directories from the staging area (but does not delete them).`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("No files specified. Usage: gvt remove <file1> <file2> ... or gvt remove .")
			return
		}

		// Ensure repo exists
		if _, err := os.Stat(".gvt"); os.IsNotExist(err) {
			fmt.Println("Not a GVT repository (no .gvt directory found)")
			return
		}

		stageFile := filepath.Join(".gvt", "stage.json")

		var staged []string
		if data, err := os.ReadFile(stageFile); err == nil {
			json.Unmarshal(data, &staged)
		} else {
			fmt.Println("Nothing is staged.")
			return
		}

		// Expand any directories passed
		var targets []string
		for _, arg := range args {
			info, err := os.Stat(arg)
			if err != nil {
				fmt.Printf("File not found: %s\n", arg)
				continue
			}

			if info.IsDir() {
				filepath.Walk(arg, func(p string, info os.FileInfo, err error) error {
					if err == nil && !info.IsDir() {
						targets = append(targets, p)
					}
					return nil
				})
			} else {
				targets = append(targets, arg)
			}
		}

		// Normalize to absolute paths for matching
		targetMap := make(map[string]bool)
		for _, t := range targets {
			abs, err := filepath.Abs(t)
			if err == nil {
				targetMap[abs] = true
			}
		}

		// Filter staged list
		var newStage []string
		removedCount := 0
		for _, s := range staged {
			if !targetMap[s] {
				newStage = append(newStage, s)
			} else {
				removedCount++
				fmt.Printf("Removed %s\n", relPath(s))
			}
		}

		if removedCount == 0 {
			fmt.Println("No matching files were staged.")
			return
		}

		data, _ := json.MarshalIndent(newStage, "", "  ")
		if err := os.WriteFile(stageFile, data, 0644); err != nil {
			fmt.Printf("Error saving stage file: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
