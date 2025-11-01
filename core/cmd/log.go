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
	"time"

	"github.com/spf13/cobra"
)

type LogEntry struct {
	ID        string `json:"id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	User      string `json:"user"`
}

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show commit history",
	Long:  `Displays the commit history with commit IDs, authors, timestamps, and messages.`,
	Run: func(cmd *cobra.Command, args []string) {
		historyFile := filepath.Join(".gvt", "history.json")
		if _, err := os.Stat(historyFile); os.IsNotExist(err) {
			fmt.Println("No commits yet.")
			return
		}

		data, err := os.ReadFile(historyFile)
		if err != nil {
			fmt.Printf("Failed to read history: %v\n", err)
			return
		}

		var rawHistory []map[string]string
		json.Unmarshal(data, &rawHistory)

		fmt.Println("GVT Commit History:")
		fmt.Println("------------------")
		for i := 0; i < len(rawHistory); i++ { // oldest first
			entry := rawHistory[i]
			ts := entry["timestamp"]
			t, _ := time.Parse(time.RFC3339, ts)
			author := entry["user"]
			if author == "" {
				author = "unknown"
			}
			fmt.Printf("commit %s\n", entry["id"])
			fmt.Printf("Author: %s\n", author)
			fmt.Printf("Date:   %s\n", t.Format("2006-01-02 15:04:05"))
			fmt.Printf("Message:\n    %s\n\n", entry["message"])
		}
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
