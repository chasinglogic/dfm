package cmd

import (
	"github.com/spf13/cobra"
)

var runHookCmd = &cobra.Command{
	Use:     "run-hook",
	Aliases: []string{"rh"},
	Short:   "Run dfm hooks by name",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		hookName := args[0]
		return currentProfile.RunHook(hookName)
	},
}
