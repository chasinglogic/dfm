// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package commands

import (
	"fmt"
	"os"

	"github.com/chasinglogic/dfm/config"
	"github.com/spf13/cobra"
)

// Sync will sync the current profile with the configured backend
var Sync = &cobra.Command{
	Use:   "sync",
	Short: "sync the current profile with the configured backend",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		profile := config.CurrentProfile()
		if err := profile.Sync(); err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(1)
		}
	},
}
