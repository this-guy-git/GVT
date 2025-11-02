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
	"github.com/this-guy-git/GVT/core/internal/storage"
)

var mergeCmd = &cobra.Command{
	Use:   "merge <branch>",
	Short: "Merge another branch into the current branch",
	Long:  `Merges the specified branch into the current branch. Only fast-forward merges are supported.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetBranch := args[0]
		currentBranch := GetCurrentBranch()

		if targetBranch == currentBranch {
			fmt.Println("Cannot merge branch into itself.")
			return
		}

		// Check branch exists
		targetRef := filepath.Join(".gvt", "refs", "heads", targetBranch)
		if _, err := os.Stat(targetRef); os.IsNotExist(err) {
			fmt.Printf("Branch '%s' does not exist.\n", targetBranch)
			return
		}

		// Get last commits for both branches
		currentCommit := getLastCommitID(currentBranch)
		targetCommit := getLastCommitID(targetBranch)

		if targetCommit == "" {
			fmt.Printf("Branch '%s' has no commits.\n", targetBranch)
			return
		}

		if currentCommit == "" {
			// No commits on current branch, fast-forward
			fmt.Printf("Current branch '%s' has no commits. Fast-forwarding.\n", currentBranch)
			restoreBranchCommit(targetBranch, targetCommit)
			return
		}

		if currentCommit == targetCommit {
			fmt.Println("Branches are already identical.")
			return
		}

		// Simple fast-forward check: if targetCommit contains all current branch commits
		currentMeta := loadCommitMeta(currentBranch, currentCommit)
		targetMeta := loadCommitMeta(targetBranch, targetCommit)

		conflicts := []string{}
		for _, f := range targetMeta.Files {
			for _, cf := range currentMeta.Files {
				if f.Path == cf.Path && f.Hash != cf.Hash {
					conflicts = append(conflicts, f.Path)
				}
			}
		}

		if len(conflicts) > 0 {
			fmt.Println("Merge aborted due to conflicts in files:")
			for _, c := range conflicts {
				fmt.Println(" -", c)
			}
			fmt.Println("Resolve conflicts manually and commit.")
			return
		}

		// No conflicts, copy files from target branch to current branch
		restoreBranchCommit(targetBranch, targetCommit)
		fmt.Printf("Merged branch '%s' into '%s'.\n", targetBranch, currentBranch)
	},
}

// restoreBranchCommit copies commit files to current branch working dir
func restoreBranchCommit(branch, commitID string) {
	commitDir := filepath.Join(".gvt", "commits", branch, commitID)
	meta := loadCommitMeta(branch, commitID)

	for _, f := range meta.Files {
		src := filepath.Join(commitDir, f.Path+".zlib")
		dst := f.Path
		os.MkdirAll(filepath.Dir(dst), 0755)

		storage.DecompressToFile(src, dst)
	}
}

// loadCommitMeta loads CommitMeta from branch/commit
func loadCommitMeta(branch, commitID string) CommitMeta {
	metaFile := filepath.Join(".gvt", "commits", branch, commitID, "meta.json")
	data, err := os.ReadFile(metaFile)
	if err != nil {
		return CommitMeta{}
	}
	var meta CommitMeta
	json.Unmarshal(data, &meta)
	return meta
}

func init() {
	rootCmd.AddCommand(mergeCmd)
}
