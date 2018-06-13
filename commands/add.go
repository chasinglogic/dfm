// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chasinglogic/dfm/config"
	"github.com/chasinglogic/dfm/git"
	"github.com/chasinglogic/dfm/linking"
	"github.com/spf13/cobra"
)

func renameAndLink(userDir, file string) error {
	s := strings.Split(file, string(filepath.Separator))
	newFile := s[len(s)-1]
	newFile = strings.TrimPrefix(newFile, ".")

	// Check if file is in XDG_config.CONFIG_HOME
	xdgConfigHome, _ := filepath.Abs(os.Getenv("XDG_CONFIG_HOME"))
	if s[len(s)-2] == ".config" || s[len(s)-2] == xdgConfigHome {
		newFile = filepath.Join("config", s[len(s)-1])
	}

	newFile = filepath.Join(userDir, newFile)

	err := os.Rename(file, newFile)
	if err != nil {
		fmt.Println("Encountered error:", err)
		fmt.Println("Trying to create intermediate directories...")

		err = os.MkdirAll(filepath.Dir(newFile), 0700)
		if err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(1)
		}

		err = os.Rename(file, newFile)
		if err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(1)
		}
	}

	linking.CreateSymlinks(userDir, os.Getenv("HOME"), linking.Config{false, false}, nil)
	return nil
}

// Add will add the specified profile to the current profile, linking it as
// necessary.
var Add = &cobra.Command{
	Use:   "add",
	Short: "Add a file to the current profile.",
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			fmt.Println("Adding files:", args)
		}

		profile := config.CurrentProfile()

		for _, f := range args {
			addFileToProfile(f, profile)
		}

		err := git.Sync(profile)
		if err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(1)
		}
	},
}

func addFileToProfile(f string, profile string) {
	file, err := filepath.Abs(f)
	if err != nil {
		fmt.Println("ERROR:", err.Error())
		os.Exit(1)
	}

	if verbose {
		fmt.Println("Absolute path:", file)
	}

	err = renameAndLink(profile, file)
	if err != nil {
		fmt.Println("ERROR:", err.Error())
		os.Exit(1)
	}
}
