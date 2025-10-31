/*
Copyright © 2025 this guy Labs <thisguy@thisguylabs.com>
*/
package cmd

import (
	"bufio"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/this-guy-git/GVT/core/internal/utils"
)

var diffAll bool

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show a Git-like diff of the last two commits",
	Long:  "Displays a hunk-style diff with context lines and colored + / - for the last two commits.",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(".gvt"); os.IsNotExist(err) {
			fmt.Println("Not a GVT repository")
			return
		}

		commits := getLastTwoCommits()
		if len(commits) < 2 {
			fmt.Println("Not enough commits to show a diff.")
			return
		}

		oldCommitDir := filepath.Join(".gvt", "commits", commits[1])
		newCommitDir := filepath.Join(".gvt", "commits", commits[0])

		oldFiles := loadCommitFiles(oldCommitDir)
		newFiles := loadCommitFiles(newCommitDir)

		var files []string
		if diffAll {
			fileSet := map[string]bool{}
			for f := range oldFiles {
				fileSet[f] = true
			}
			for f := range newFiles {
				fileSet[f] = true
			}
			for f := range fileSet {
				files = append(files, f)
			}
		} else {
			for f := range newFiles {
				if !compareLines(oldFiles[f], newFiles[f]) {
					files = append(files, f)
				}
			}
		}

		if len(files) == 0 {
			fmt.Println("No changes detected between the last two commits.")
			return
		}

		for _, f := range files {
			oldLines := oldFiles[f]
			newLines := newFiles[f]

			fmt.Printf("\nFile: %s\n", f)
			printHunkDiff(oldLines, newLines)
		}
	},
}

func getLastTwoCommits() []string {
	commitsDir := filepath.Join(".gvt", "commits")
	entries, err := os.ReadDir(commitsDir)
	if err != nil || len(entries) == 0 {
		return []string{}
	}
	var commits []string
	for _, e := range entries {
		if e.IsDir() {
			commits = append(commits, e.Name())
		}
	}
	// sort ascending (oldest → newest)
	for i := 0; i < len(commits)-1; i++ {
		for j := i + 1; j < len(commits); j++ {
			if commits[i] > commits[j] {
				commits[i], commits[j] = commits[j], commits[i]
			}
		}
	}
	if len(commits) >= 2 {
		return commits[len(commits)-2:]
	}
	return commits
}

func loadCommitFiles(commitDir string) map[string][]string {
	files := map[string][]string{}
	metaFile := filepath.Join(commitDir, "meta.json")
	data, err := os.ReadFile(metaFile)
	if err != nil {
		return files
	}
	var meta CommitMeta
	json.Unmarshal(data, &meta)
	for _, f := range meta.Files {
		if utils.IsIgnored(f.Path) || shouldSkip(f.Path) {
			continue
		}
		zPath := filepath.Join(commitDir, f.Path+".zlib")
		lines := readZlibLines(zPath)
		if len(lines) > 0 {
			files[f.Path] = lines
		}
	}
	return files
}

func readZlibLines(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		return []string{}
	}
	defer file.Close()
	r, _ := zlib.NewReader(file)
	defer r.Close()
	scanner := bufio.NewScanner(r)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if isBinaryLine(line) {
			return []string{}
		}
		lines = append(lines, line)
	}
	return lines
}

func isBinaryLine(line string) bool {
	for i := 0; i < len(line); i++ {
		if line[i] == 0 {
			return true
		}
	}
	return false
}

// printHunkDiff prints a Git-like colored hunk diff
func printHunkDiff(oldLines, newLines []string) {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	// Simple line-by-line diff
	i, j := 0, 0
	for i < len(oldLines) || j < len(newLines) {
		if i < len(oldLines) && j < len(newLines) && oldLines[i] == newLines[j] {
			// Context line
			fmt.Printf("  %s\n", oldLines[i])
			i++
			j++
		} else {
			startOld, startNew := i+1, j+1
			hunkOld, hunkNew := []string{}, []string{}

			// Collect removed lines
			for i < len(oldLines) && (j >= len(newLines) || oldLines[i] != newLines[j]) {
				hunkOld = append(hunkOld, oldLines[i])
				i++
			}

			// Collect added lines
			for j < len(newLines) && (i >= len(oldLines) || oldLines[i-1] != newLines[j]) {
				hunkNew = append(hunkNew, newLines[j])
				j++
			}

			fmt.Printf("@@ -%d,%d +%d,%d @@\n", startOld, len(hunkOld), startNew, len(hunkNew))
			for _, l := range hunkOld {
				fmt.Printf("%s %s\n", red("-"), l)
			}
			for _, l := range hunkNew {
				fmt.Printf("%s %s\n", green("+"), l)
			}
		}
	}
}

func compareLines(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func init() {
	diffCmd.Flags().BoolVar(&diffAll, "all", false, "Show diff for all files in repository")
	rootCmd.AddCommand(diffCmd)
}
