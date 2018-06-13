// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package commands

import (
	"fmt"
	"os"

	"github.com/chasinglogic/dfm/config"
	"github.com/chasinglogic/dfm/filemap"
	"github.com/chasinglogic/dfm/git"
	"github.com/chasinglogic/dfm/linking"
	"github.com/spf13/cobra"
)

func init() {
	Link.Flags().BoolVarP(&overwrite, "overwrite", "o", false,
		"if provided dfm will remove files that exist where a link should go")
}

func linkModule(module config.Module, profile string, mappings filemap.Mappings) {
	location := module.Location(config.ModuleDir(profile))
	if _, err := os.Stat(location); os.IsNotExist(err) {
		err = git.RunGitCMD(profile, "clone", module.Repo, location)
		if err != nil {
			fmt.Println("ERROR: Unable to clone module:", err)
			return
		}
	}

	moduleMappings := append(mappings, module.Mappings...)

	err := linking.CreateSymlinks(
		location,
		os.Getenv("HOME"),
		linking.Config{
			DryRun:    dryRun,
			Overwrite: overwrite,
		},
		moduleMappings,
	)

	if err != nil {
		fmt.Println("ERROR:", err.Error())
		os.Exit(1)
	}
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

		fmt.Println("Linking profile", profile)

		mappings := filemap.DefaultMappings()
		dfmyml := config.LoadDotDFM(profile)

		for _, module := range dfmyml.PreLinkModules() {
			linkModule(module, profile, mappings)
		}

		parentMappings := append(mappings, dfmyml.Mappings...)
		err := linking.CreateSymlinks(
			profile,
			os.Getenv("HOME"),
			linking.Config{
				DryRun:    dryRun,
				Overwrite: overwrite,
			},
			parentMappings,
		)

		if err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(1)
		}

		for _, module := range dfmyml.PostLinkModules() {
			linkModule(module, profile, mappings)
		}
	},
}
