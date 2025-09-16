/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/chasinglogic/dfm/internal/state"
	"github.com/spf13/cobra"
)

var runHookCmd = &cobra.Command{
	Use:     "run-hook <HOOK_NAME>",
	Short:   "Runs the given hook without invoking the associated event",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"rh"},
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, err := loadProfile(state.State.CurrentProfile)
		if err != nil {
			return err
		}

		return profile.RunHook(args[0])
	},
}

func init() {
	RootCmd.AddCommand(runHookCmd)
}
