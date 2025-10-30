/*
Copyright © 2025 this guy Labs <thisguy@thisguylabs.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gvt",
	Short: "GVT is a Version control tool similar to Git.",
	Long: `GVT is a Version control tool similar to Git
With features that make it easier to use for beginners.

Examples:
  gvt init
  gvt config user.name "Your Name"
  gvt config user.email "
  gvt commit
  gvt log
  gvt status
  gvt add <file>
  gvt remove <file>
  gvt push origin main
  gvt pull origin main
  `,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gvt.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
