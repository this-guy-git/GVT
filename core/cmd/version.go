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
