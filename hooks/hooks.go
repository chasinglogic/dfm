// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package hooks

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/chasinglogic/dfm/config"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

// Hooks is a map of "hook_name" to a slice of string shell commands to run
type Hooks map[string][]string

// DFMYml is used for extending DFM. It is the .dfm.yml file found
// in the root of a profile.
type DFMYml struct {
	Hooks Hooks `yaml:"hooks"`
}

// Load will load the hooks file for the current Profile
func Load() Hooks {
	userDir := filepath.Join(config.ProfileDir(), config.CurrentProfile)

	dfmyml, err := ioutil.ReadFile(filepath.Join(userDir, ".dfm.yml"))
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("ERROR loading hooks:", err.Error())
		}

		return nil
	}

	var yml DFMYml

	err = yaml.Unmarshal(dfmyml, &yml)
	if err != nil {
		fmt.Println("ERROR loading hooks:", err.Error())
		return yml.Hooks
	}

	return yml.Hooks
}

// AddHooks will add before and after hooks to the given command.
func AddHooks(command *cobra.Command) *cobra.Command {
	// Store this for later use
	runFunc := command.Run

	command.Run = func(cmd *cobra.Command, args []string) {
		hooks := Load()
		prof := config.CurrentProfile

		commands, preHooks := hooks["before_"+command.Use]
		if preHooks {
			RunCommands(commands)
		}

		// Run the real command
		runFunc(cmd, args)

		if prof != config.CurrentProfile {
			// Reload if profile changed
			hooks = Load()
		}

		commands, postHooks := hooks["after_"+command.Use]
		if postHooks {
			RunCommands(commands)
		}
	}

	return command
}
