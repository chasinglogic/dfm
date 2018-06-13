// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chasinglogic/dfm/config"
	"github.com/chasinglogic/dfm/git"
	"github.com/spf13/cobra"
)

// Init will create a new profile with the given name.
var Init = &cobra.Command{
	Use:   "init",
	Short: "create a new profile with `NAME`",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please specify a profile name.")
			os.Exit(1)
		}

		profileDir := config.ProfileDir()
		profile := filepath.Join(profileDir, args[0])
		err := os.MkdirAll(profile, os.ModePerm)
		if err != nil {
			fmt.Println("ERROR: Unable create profile directory: ", profile, err)
			os.Exit(1)
		}

		err = git.Init(profile)
		if err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(1)
		}
	},
}
