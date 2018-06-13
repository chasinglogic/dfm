// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package commands

import (
	"github.com/chasinglogic/dfm/config"
	"github.com/spf13/cobra"
)

// Global variables to represent root flags available to sub commands
var (
	verbose bool
	dryRun  bool

	// Whether or not to overwrite existing files when linking
	overwrite bool
)

func loadHooks(profile string) config.Hooks {
	profileCfg := config.LoadDotDFM(profile)
	return profileCfg.Hooks
}

func init() {
	Root.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	Root.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "don't make changes just print what would happen")

	Root.AddCommand(config.AddHooks(loadHooks, Init))
	Root.AddCommand(config.AddHooks(loadHooks, Add))
	Root.AddCommand(config.AddHooks(loadHooks, Link))
	Root.AddCommand(config.AddHooks(loadHooks, List))
	Root.AddCommand(config.AddHooks(loadHooks, Remove))
	Root.AddCommand(config.AddHooks(loadHooks, Where))
	Root.AddCommand(config.AddHooks(loadHooks, Sync))
	Root.AddCommand(config.AddHooks(loadHooks, Clean))
	Root.AddCommand(RunHook)
	Root.AddCommand(Git)
	Root.AddCommand(Clone)
}

// Root is the root dfm command.
var Root = &cobra.Command{
	Use:   "dfm",
	Short: "Manage dotfiles.",
	Long: `Dotfile management written for pair programmers. Examples on getting
started with dfm are avialable at https://github.com/chasinglogic/dfm`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.Init()
	},
}

// Execute aliases to running Execute on the root command
func Execute() error {
	return Root.Execute()
}
