/*
Copyright © 2025 this guy Labs <thisguy@thisguylabs.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const GVTVersion = "v0.1.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Return the current GVT version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(`
	     /$$$$$$  /$$    /$$ /$$$$$$$$
	    /$$__  $$| $$   | $$|__  $$__/
	   | $$  \__/| $$   | $$   | $$   
	   | $$ /$$$$|  $$ / $$/   | $$   
	   | $$|_  $$ \  $$ $$/    | $$   
	   | $$  \ $$  \  $$$/     | $$   
	   |  $$$$$$/   \  $/      | $$   
	    \______/     \_/       |__/   
		`)
		fmt.Println("GVT Version:", GVTVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
