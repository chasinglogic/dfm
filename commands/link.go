// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package commands

import (
	"fmt"
	"os"

	"github.com/chasinglogic/dfm/config"
	"github.com/chasinglogic/dfm/dotdfm"
	"github.com/chasinglogic/dfm/filemap"
	"github.com/chasinglogic/dfm/linking"
	"github.com/spf13/cobra"
)

func init() {
	Link.Flags().BoolVarP(&overwrite, "overwrite", "o", false,
		"whether dfm should remove files that exist where a link should go")
}

// Link will generate and create the symlinks to the dotfiles in the repo.
var Link = &cobra.Command{
	Use:   "link",
	Short: "link the profile with `NAME`",
	Long:  "will generate and create the symlinks to the dotfiles in the profile",
	Run: func(cmd *cobra.Command, args []string) {
		profile := config.CurrentProfile()
		if len(args) >= 1 {
			profile = config.GetProfileByName(args[0])
		}

		fmt.Println("Linking profile", profile.Name)

		mappings := filemap.DefaultMappings()
		for _, location := range profile.Locations {
			dfmyml := dotdfm.LoadDotDFM(location)
			mappings = append(mappings, dfmyml.Mappings...)
		}

		err := profile.Link(
			os.Getenv("HOME"),
			mappings,
			linking.Config{
				dryRun,
				overwrite,
			},
		)

		if err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(1)
		}

		config.SetCurrentProfile(profile)
	},
}
