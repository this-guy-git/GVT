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
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var restoreCmd = &cobra.Command{
	Use:   "restore <file>",
	Short: "Restore a file to its last committed state",
	Long:  `Reverts the specified file to the state it had in the last commit.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(".gvt"); os.IsNotExist(err) {
			fmt.Println("Not a GVT repository (no .gvt directory found)")
			return
		}

		fileToRestore := args[0]
		lastCommit := getLastCommitID()
		if lastCommit == "" {
			fmt.Println("No commits found to restore from.")
			return
		}

		metaFile := filepath.Join(".gvt", "commits", lastCommit, "meta.json")
		if _, err := os.Stat(metaFile); os.IsNotExist(err) {
			fmt.Println("Last commit metadata not found.")
			return
		}

		var meta CommitMeta
		data, _ := os.ReadFile(metaFile)
		json.Unmarshal(data, &meta)

		var targetFile *FileMeta
		for _, f := range meta.Files {
			if f.Path == fileToRestore || f.Path == filepath.ToSlash(fileToRestore) {
				targetFile = &f
				break
			}
		}

		if targetFile == nil {
			fmt.Printf("File %s is not tracked in the last commit.\n", fileToRestore)
			return
		}

		zPath := filepath.Join(".gvt", "commits", lastCommit, targetFile.Path+".zlib")
		srcFile, err := os.Open(zPath)
		if err != nil {
			fmt.Printf("Failed to open %s: %v\n", zPath, err)
			return
		}
		defer srcFile.Close()

		r, err := zlib.NewReader(srcFile)
		if err != nil {
			fmt.Printf("Failed to decompress %s: %v\n", zPath, err)
			return
		}
		defer r.Close()

		os.MkdirAll(filepath.Dir(targetFile.Path), 0755)
		dstFile, _ := os.Create(targetFile.Path)
		defer dstFile.Close()

		buf := make([]byte, 32*1024) // 32 KB chunks
		for {
			n, err := r.Read(buf)
			if n > 0 {
				dstFile.Write(buf[:n])
			}
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Printf("Error restoring %s: %v\n", targetFile.Path, err)
				return
			}
		}

		fmt.Printf("Restored: %s to last committed state.\n", targetFile.Path)
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
}
