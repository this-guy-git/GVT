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

	"github.com/this-guy-git/GVT/core/internal/core"
)

var local bool
var reason string

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Set or get GVT configuration",
	Long: `Examples:
    gvt config user.name "Your Name"
    gvt config user.email "you@example.com"
    gvt config --get user.name
    gvt config --get user.email
    `,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Usage: gvt config <key> [value]")
			return
		}

		key := args[0]
		config, _ := core.LoadConfig(local)

		if len(args) == 1 { // get
			switch key {
			case "user.name":
				reason = "Username:"
				fmt.Println(reason, config.UserName)

			case "user.email":
				reason = "Email:"
				fmt.Println(reason, config.UserEmail)

			default:
				fmt.Println("Unknown key:", key)
			}
			return
		}

		value := args[1] // set
		switch key {
		case "user.name":
			config.UserName = value
		case "user.email":
			config.UserEmail = value
		default:
			fmt.Println("Unknown key:", key)
			return
		}

		if err := core.SaveConfig(config, local); err != nil {
			fmt.Println("Error saving config:", err)
		} else {
			scope := "global"
			if local {
				scope = "local"
			}
			fmt.Printf("Set %s to %s (%s)\n", key, value, scope)
		}
	},
}

func init() {
	configCmd.Flags().BoolVarP(&local, "local", "l", false, "Use local repo config instead of global")
	configCmd.Flags().BoolP("get", "g", false, "Get the value of a config key")
	rootCmd.AddCommand(configCmd)
}
