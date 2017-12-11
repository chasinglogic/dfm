// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package commands

import (
	"fmt"

	"github.com/chasinglogic/dfm/backend"
	"github.com/chasinglogic/dfm/backend/dropbox"
	"github.com/chasinglogic/dfm/backend/git"
	"github.com/chasinglogic/dfm/config"
	"github.com/chasinglogic/dfm/hooks"
	"github.com/spf13/cobra"
)

// Global variables to represent root flags available to sub commands
var (
	Verbose bool
	DryRun  bool

	// Whether or not to overwrite existing files when linking
	overwrite bool

	Backend = loadBackend(config.Backend)
)

func init() {
	Root.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	Root.PersistentFlags().BoolVarP(&DryRun, "dry-run", "d", false, "don't make changes just print what would happen")

	Root.AddCommand(hooks.AddHooks(Init))
	Root.AddCommand(hooks.AddHooks(Add))
	Root.AddCommand(hooks.AddHooks(Link))
	Root.AddCommand(hooks.AddHooks(List))
	Root.AddCommand(hooks.AddHooks(Remove))
	Root.AddCommand(hooks.AddHooks(Where))
	Root.AddCommand(hooks.AddHooks(Sync))
	Root.AddCommand(hooks.AddHooks(Clean))
	Root.AddCommand(RunHook)

	for _, c := range Backend.Commands() {
		Root.AddCommand(hooks.AddHooks(c))
	}
}

func loadBackend(backendName string) backend.Backend {
	switch backendName {
	case "git":
		return git.Backend{}
	case "dropbox":
		return dropbox.Backend{}
	default:
		fmt.Printf("Backend \"%s\" not found defaulting to git\n.", backendName)
		return git.Backend{}
	}
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
