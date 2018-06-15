// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package commands

import (
	"fmt"
	"os"

	"github.com/chasinglogic/dfm/config"
	"github.com/chasinglogic/dfm/git"
	"github.com/spf13/cobra"
)

var syncModules bool

func init() {
	Sync.Flags().BoolVarP(&syncModules, "modules", "m", false,
		"if provided dfm will sync modules as well as the primary profile")
}

// Sync will sync the current profile with the configured backend
var Sync = &cobra.Command{
	Use:   "sync",
	Short: "sync the current profile with the configured backend",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		profile := config.CurrentProfile()
		if err := git.Sync(profile); err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(1)
		}

		yml := config.LoadDotDFM(profile)
		if !(yml.SyncModules || syncModules) {
			return
		}

		moduleDir := config.ModuleDir()
		for _, module := range yml.Modules {
			location := module.Location()
			if module.PullOnly {
				err := git.Pull(location)
				if err != nil {
					fmt.Println("ERROR: Unable to update module:", err)
				}
			} else {
				err := git.Sync(location)
				if err != nil {
					fmt.Println("ERROR: Unable to sync module:", err)
				}
			}
		}
	},
}
