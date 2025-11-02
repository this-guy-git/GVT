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
	"bytes"
	"compress/zlib"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/this-guy-git/GVT/core/internal/utils"
)

var commitMsg string

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Record the current staged changes as a new commit",
	Long:  `Saves the currently staged files as a new commit. If no message is given, one will be generated automatically based on the changes.`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(".gvt"); os.IsNotExist(err) {
			fmt.Println("Not a GVT repository (no .gvt directory found)")
			return
		}

		currentBranch := getCurrentBranch()
		if currentBranch == "" {
			fmt.Println("No branch found. (HEAD is detached or invalid)")
			return
		}

		stageFile := filepath.Join(".gvt", "stage.json")
		if _, err := os.Stat(stageFile); os.IsNotExist(err) {
			fmt.Println("Nothing staged to commit.")
			return
		}

		// Load staged files
		var staged []string
		if data, err := os.ReadFile(stageFile); err == nil {
			json.Unmarshal(data, &staged)
		}
		if len(staged) == 0 {
			fmt.Println("Nothing staged to commit.")
			return
		}

		lastCommitID := getLastCommitID(currentBranch)
		lastHashes := map[string]string{}
		if lastCommitID != "" {
			metaFile := filepath.Join(".gvt", "commits", currentBranch, lastCommitID, "meta.json")
			if data, err := os.ReadFile(metaFile); err == nil {
				var meta CommitMeta
				json.Unmarshal(data, &meta)
				for _, f := range meta.Files {
					lastHashes[f.Path] = f.Hash
				}
			}
		}

		var filesToCommit []FileMeta
		for _, f := range staged {
			info, err := os.Stat(f)
			if err != nil {
				continue
			}
			if info.IsDir() {
				filepath.WalkDir(f, func(path string, d fs.DirEntry, err error) error {
					if err != nil || d.IsDir() {
						return nil
					}
					if shouldSkip(path) {
						return nil
					}
					if fm := prepareFile(path, lastHashes); fm.Path != "" {
						filesToCommit = append(filesToCommit, fm)
					}
					return nil
				})
			} else {
				if fm := prepareFile(f, lastHashes); fm.Path != "" {
					filesToCommit = append(filesToCommit, fm)
				}
			}
		}

		if len(filesToCommit) == 0 {
			fmt.Println("No changes to commit.")
			return
		}

		if commitMsg == "" {
			if len(filesToCommit) == 1 {
				commitMsg = fmt.Sprintf("Changed 1 file: %s", filesToCommit[0].Path)
			} else {
				commitMsg = fmt.Sprintf("Changed %d files", len(filesToCommit))
			}
		}

		commitID := time.Now().Format("20060102-150405")
		commitDir := filepath.Join(".gvt", "commits", currentBranch, commitID)
		os.MkdirAll(commitDir, 0755)

		historyFile := filepath.Join(".gvt", "history.json")

		var history []map[string]string
		if data, err := os.ReadFile(historyFile); err == nil {
			json.Unmarshal(data, &history)
		}

		user := utils.GetCommitUser(".")

		history = append(history, map[string]string{
			"id":        commitID,
			"branch":    currentBranch,
			"message":   commitMsg,
			"timestamp": time.Now().Format(time.RFC3339),
			"user":      user,
		})

		data, _ := json.MarshalIndent(history, "", "  ")
		os.WriteFile(historyFile, data, 0644)

		for _, f := range filesToCommit {
			dest := filepath.Join(commitDir, f.Path+".zlib")
			os.MkdirAll(filepath.Dir(dest), 0755)
			os.WriteFile(dest, f.Data, 0644)
		}

		meta := CommitMeta{
			ID:        commitID,
			Message:   commitMsg,
			Timestamp: time.Now().Format(time.RFC3339),
			Files:     filesToCommit,
			Parent:    lastCommitID,
			Branch:    currentBranch,
		}
		data, _ = json.MarshalIndent(meta, "", "  ")
		os.WriteFile(filepath.Join(commitDir, "meta.json"), data, 0644)

		// Update branch ref
		refPath := filepath.Join(".gvt", "refs", "heads", currentBranch)
		os.WriteFile(refPath, []byte(commitID), 0644)

		// Clear stage
		os.WriteFile(stageFile, []byte("[]"), 0644)

		fmt.Printf("[%s] committed %d file(s) as %s\nMessage: %s\n", currentBranch, len(filesToCommit), commitID, commitMsg)
	},
}

type CommitMeta struct {
	ID        string     `json:"id"`
	Message   string     `json:"message"`
	Timestamp string     `json:"timestamp"`
	Files     []FileMeta `json:"files"`
	Parent    string     `json:"parent,omitempty"`
	Branch    string     `json:"branch"`
}

type FileMeta struct {
	Path string `json:"path"`
	Hash string `json:"hash"`
	Data []byte `json:"-"`
}

func getLastCommitID(branch string) string {
	commitsDir := filepath.Join(".gvt", "commits", branch)
	entries, err := os.ReadDir(commitsDir)
	if err != nil || len(entries) == 0 {
		return ""
	}
	var last string
	for _, e := range entries {
		if e.IsDir() && e.Name() > last {
			last = e.Name()
		}
	}
	return last
}

func prepareFile(path string, lastHashes map[string]string) FileMeta {
	file, err := os.Open(path)
	if err != nil {
		return FileMeta{}
	}
	defer file.Close()

	h := md5.New()
	var compressedBuf bytes.Buffer
	w := zlib.NewWriter(&compressedBuf)
	buf := make([]byte, 32*1024)

	for {
		n, err := file.Read(buf)
		if n > 0 {
			h.Write(buf[:n])
			w.Write(buf[:n])
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return FileMeta{}
		}
	}
	w.Close()

	hashStr := hex.EncodeToString(h.Sum(nil))
	relPath, _ := filepath.Rel(".", path)
	relPath = filepath.ToSlash(relPath)

	if lastHashes[relPath] == hashStr {
		return FileMeta{}
	}

	return FileMeta{
		Path: relPath,
		Hash: hashStr,
		Data: compressedBuf.Bytes(),
	}
}

func shouldSkip(path string) bool {
	return strings.Contains(path, string(os.PathSeparator)+".gvt")
}

func init() {
	commitCmd.Flags().StringVarP(&commitMsg, "message", "m", "", "Commit message")
	rootCmd.AddCommand(commitCmd)
}
