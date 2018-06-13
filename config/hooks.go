// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package config

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

// Hooks is a map of "hook_name" to a slice of string shell commands to run
type Hooks map[string][]string

// AddHooks will add before and after hooks to the given command.
func AddHooks(loadHooks func(profile string) Hooks, command *cobra.Command) *cobra.Command {
	// Store this for later use
	runFunc := command.Run

	command.Run = func(cmd *cobra.Command, args []string) {
		prof := CurrentProfile()
		hooks := loadHooks(prof)

		commands, preHooks := hooks["before_"+command.Use]
		if preHooks {
			RunCommands(commands)
		}

		// Run the real command
		runFunc(cmd, args)

		if prof != CurrentProfile() {
			// Reload if profile changed
			hooks = loadHooks(CurrentProfile())
		}

		commands, postHooks := hooks["after_"+command.Use]
		if postHooks {
			RunCommands(commands)
		}
	}

	return command
}

// TODO: write a runCommand for windows

// RunCommands will run the given slice of strings each as their own command
func RunCommands(commands []string) {
	for _, cmd := range commands {
		c := exec.Command("bash", "-c", cmd)
		out, err := c.CombinedOutput()
		if err != nil {
			fmt.Println("ERROR Running Command:", cmd, err.Error())
			return
		}
		fmt.Print(string(out))
	}
}
