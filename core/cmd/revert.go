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

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

const (
	Green = "\033[32m"
	Red   = "\033[31m"
	Cyan  = "\033[36m"
	Reset = "\033[0m"
)

var revertCmd = &cobra.Command{
	Use:   "revert [commit-id]",
	Short: "Revert working directory to a previous commit",
	Long:  `Replaces working directory files with the state from the specified commit. If no commit ID is given, reverts to the last commit.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(".gvt"); os.IsNotExist(err) {
			fmt.Println(Red + "Not a GVT repository (no .gvt directory found)" + Reset)
			return
		}

		currentBranch := getCurrentBranch()
		var targetCommit string
		if len(args) == 1 {
			targetCommit = args[0]
		} else {
			targetCommit = getLastCommitID(currentBranch)
			if targetCommit == "" {
				fmt.Println(Red + "No commits found to revert to." + Reset)
				return
			}
		}

		commitDir := filepath.Join(".gvt", "commits", currentBranch, targetCommit)
		metaFile := filepath.Join(commitDir, "meta.json")
		if _, err := os.Stat(metaFile); os.IsNotExist(err) {
			fmt.Printf(Red+"Commit %s does not exist.\n"+Reset, targetCommit)
			return
		}

		var meta CommitMeta
		data, _ := os.ReadFile(metaFile)
		json.Unmarshal(data, &meta)

		fmt.Printf(Cyan+"Reverting to commit %s...\n"+Reset, targetCommit)
		bar := progressbar.NewOptions(len(meta.Files),
			progressbar.OptionSetDescription("Restoring files:"),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "=",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}),
		)

		for _, f := range meta.Files {
			zPath := filepath.Join(commitDir, f.Path+".zlib")
			targetPath := filepath.Join(".", f.Path)

			srcFile, err := os.Open(zPath)
			if err != nil {
				fmt.Printf(Red+"Failed to read %s\n"+Reset, zPath)
				bar.Add(1)
				continue
			}

			r, err := zlib.NewReader(srcFile)
			if err != nil {
				fmt.Printf(Red+"Failed to decompress %s: %v\n"+Reset, zPath, err)
				srcFile.Close()
				bar.Add(1)
				continue
			}

			os.MkdirAll(filepath.Dir(targetPath), 0755)
			dstFile, _ := os.Create(targetPath)

			buf := make([]byte, 32*1024)
			for {
				n, err := r.Read(buf)
				if n > 0 {
					dstFile.Write(buf[:n])
				}
				if err != nil {
					if err == io.EOF {
						break
					}
					fmt.Printf(Red+"Error decompressing %s: %v\n"+Reset, f.Path, err)
					break
				}
			}

			r.Close()
			srcFile.Close()
			dstFile.Close()

			bar.Add(1)
		}

		os.WriteFile(filepath.Join(".gvt", "stage.json"), []byte("[]"), 0644)
		fmt.Printf(Green+"\nRevert to commit %s completed.\n"+Reset, targetCommit)
	},
}

func init() {
	rootCmd.AddCommand(revertCmd)
}
