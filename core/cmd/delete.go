package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var forceDelete bool
var deleteAll bool

var deleteCmd = &cobra.Command{
	Use:   "delete [commit-id]",
	Short: "Delete a specific commit or all commits",
	Long:  `Removes a commit from history and deletes its folder. Use --force to delete the latest commit, or --all to delete all commits.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(".gvt"); os.IsNotExist(err) {
			fmt.Println("Not a GVT repository.")
			return
		}

		commitsDir := filepath.Join(".gvt", "commits")
		historyFile := filepath.Join(".gvt", "history.json")

		if deleteAll {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Are you sure you want to delete ALL commits? [y/N]: ")
			resp, _ := reader.ReadString('\n')
			if strings.ToLower(strings.TrimSpace(resp)) != "y" {
				fmt.Println("Aborted.")
				return
			}

			os.RemoveAll(commitsDir)
			os.MkdirAll(commitsDir, 0755)
			os.WriteFile(historyFile, []byte("[]"), 0644)
			fmt.Println("All commits deleted successfully.")
			return
		}

		if len(args) != 1 {
			fmt.Println("Please provide a commit ID to delete, or use --all to delete everything.")
			return
		}

		commitID := args[0]
		commitDir := filepath.Join(commitsDir, commitID)
		if _, err := os.Stat(commitDir); os.IsNotExist(err) {
			fmt.Printf("Commit %s does not exist.\n", commitID)
			return
		}

		lastCommit := getLastCommitID()
		if commitID == lastCommit && !forceDelete {
			fmt.Println("Cannot delete the latest commit without --force.")
			return
		}

		// Remove folder
		err := os.RemoveAll(commitDir)
		if err != nil {
			fmt.Printf("Failed to delete commit folder: %v\n", err)
			return
		}

		// Update history.json
		var history []map[string]string
		if data, err := os.ReadFile(historyFile); err == nil {
			json.Unmarshal(data, &history)
		}

		newHistory := []map[string]string{}
		for _, h := range history {
			if h["id"] != commitID {
				newHistory = append(newHistory, h)
			}
		}

		data, _ := json.MarshalIndent(newHistory, "", "  ")
		os.WriteFile(historyFile, data, 0644)

		fmt.Printf("Commit %s deleted successfully.\n", commitID)
	},
}

func init() {
	deleteCmd.Flags().BoolVarP(&forceDelete, "force", "f", false, "Force delete latest commit")
	deleteCmd.Flags().BoolVar(&deleteAll, "all", false, "Delete all commits")
	rootCmd.AddCommand(deleteCmd)
}
