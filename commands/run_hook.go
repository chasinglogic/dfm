// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package commands

import (
	"fmt"
	"os"

	"github.com/chasinglogic/dfm/config"
	"github.com/chasinglogic/dfm/dotdfm"
	"github.com/chasinglogic/dfm/hooks"
	"github.com/spf13/cobra"
)

// RunHook will add the specified profile to the current profile, linking it as
// necessary.
var RunHook = &cobra.Command{
	Use:   "run-hook",
	Short: "Run the hook specified by `HOOK`",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("ERROR: Must specify at least one hook to run")
			os.Exit(128)
		}

		profile := config.CurrentProfile()

		allHooks := hooks.Hooks{}
		for _, location := range profile.Locations {
			yml := dotdfm.LoadDotDFM(location)
			for hookName, hookCommands := range yml.Hooks {
				allHooks[hookName] = append(allHooks[hookName], hookCommands...)
			}
		}

		hooks.RunCommands(allHooks[args[0]])
	},
}
