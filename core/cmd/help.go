/*
Copyright © 2025 this guy Labs <thisguy@thisguylabs.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// helpCmd represents the help command
var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Show help commands for GVT",
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(helpCmd)
}
