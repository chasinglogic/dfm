// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package hooks

import (
	"github.com/chasinglogic/dfm/config"
	"github.com/chasinglogic/dfm/dotfiles"
	"github.com/spf13/cobra"
)

// Hooks is a map of "hook_name" to a slice of string shell commands to run
type Hooks map[string][]string

// AddHooks will add before and after hooks to the given command.
func AddHooks(loadHooks func(profile dotfiles.Profile) Hooks, command *cobra.Command) *cobra.Command {
	// Store this for later use
	runFunc := command.Run

	command.Run = func(cmd *cobra.Command, args []string) {
		prof := config.CurrentProfile().Name
		hooks := loadHooks(config.CurrentProfile())

		commands, preHooks := hooks["before_"+command.Use]
		if preHooks {
			RunCommands(commands)
		}

		// Run the real command
		runFunc(cmd, args)

		if prof != config.CurrentProfile().Name {
			// Reload if profile changed
			hooks = loadHooks(config.CurrentProfile())
		}

		commands, postHooks := hooks["after_"+command.Use]
		if postHooks {
			RunCommands(commands)
		}
	}

	return command
}
