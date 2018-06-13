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
	"github.com/chasinglogic/dfm/filemap"
	"github.com/chasinglogic/dfm/linking"
	"github.com/spf13/cobra"
)

// Remove will remove the specified profile
var Remove = &cobra.Command{
	Use:   "remove",
	Short: "remove the profile with `NAME`",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		profile := args[0]
		userDir := filepath.Join(config.ProfileDir(), profile)

		dfmyml := config.LoadDotDFM(userDir)
		mappings := append(dfmyml.Mappings, filemap.DefaultMappings()...)

		links, err := linking.GenerateSymlinks(userDir, os.Getenv("HOME"), mappings)
		if err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(1)
		}

		rmerr := os.RemoveAll(userDir)
		if rmerr != nil {
			fmt.Println("ERROR:", rmerr.Error())
			os.Exit(1)
		}

		if verbose {
			fmt.Println("Removed profile directory:", userDir)
		}

		if err := RemoveSymlinks(links, args[0]); err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(1)
		}
	},
}

// RemoveSymlinks will remove all of the symlinks after removing a profile it
// will first Check if the link is still valid after removing the profile, and
// if so just verify that it doesn't contain the username of the profile
// we're deleting. If the profile we're removing is the one that was currently
// in use then both conditions should be true.
func RemoveSymlinks(l []linking.LinkInfo, username string) error {
	for _, link := range l {
		// Check if the link is still valid after removing the profile, and if
		if path, err := os.Readlink(link.Dest); err != nil ||
			strings.Contains(path, username) {

			if dryRun || verbose {
				fmt.Printf("Removing symlink %s\n", link.Dest)
			}

			if !dryRun {
				if err := os.Remove(link.Dest); err != nil {
					return err
				}
			}
		}

	}

	return nil
}
