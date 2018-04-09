// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package commands

import (
	"github.com/chasinglogic/dfm/config"
	"github.com/chasinglogic/dfm/dotdfm"
	"github.com/chasinglogic/dfm/dotfiles"
	"github.com/chasinglogic/dfm/hooks"
	"github.com/spf13/cobra"
)

// Global variables to represent root flags available to sub commands
var (
	verbose bool
	dryRun  bool

	// Whether or not to overwrite existing files when linking
	overwrite bool
)

func loadHooks(profile dotfiles.Profile) hooks.Hooks {
	var hooks hooks.Hooks

	for _, location := range profile.Locations {
		dotdfm := dotdfm.LoadDotDFM(location)

		if hooks == nil {
			hooks = dotdfm.Hooks
			continue
		}

		for hookName := range hooks {
			hooks[hookName] = append(hooks[hookName], dotdfm.Hooks[hookName]...)
		}
	}

	return hooks
}

func init() {
	Root.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	Root.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "don't make changes just print what would happen")

	Root.AddCommand(hooks.AddHooks(loadHooks, Init))
	Root.AddCommand(hooks.AddHooks(loadHooks, Add))
	Root.AddCommand(hooks.AddHooks(loadHooks, Link))
	Root.AddCommand(hooks.AddHooks(loadHooks, List))
	Root.AddCommand(hooks.AddHooks(loadHooks, Remove))
	Root.AddCommand(hooks.AddHooks(loadHooks, Where))
	Root.AddCommand(hooks.AddHooks(loadHooks, Sync))
	Root.AddCommand(hooks.AddHooks(loadHooks, Clean))
	Root.AddCommand(RunHook)
}

// Root is the root dfm command.
var Root = &cobra.Command{
	Use:   "dfm",
	Short: "Manage dotfiles.",
	Long: `Dotfile management written for pair programmers. Examples on getting
started with dfm are avialable at https://github.com/chasinglogic/dfm`,
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		_ = config.SaveConfig()
	},
}

// Execute aliases to running Execute on the root command
func Execute() error {
	return Root.Execute()
}
