package cmd

import (
	"github.com/chasinglogic/dfm/profiles"
	"github.com/spf13/cobra"
)

var (
	syncCommitMsg   = ""
	syncSkipModules = false
)

var syncCmd = &cobra.Command{
	Use:     "sync",
	Short:   "Sync your dotfiles",
	Aliases: []string{"s"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return currentProfile.Sync(profiles.SyncOptions{
			CommitMessage: syncCommitMsg,
			SkipModules:   syncSkipModules,
		})
	},
}

func init() {
	syncCmd.Flags().BoolVarP(&syncSkipModules, "skip-modules", "s", false, "do not sync modules")
	syncCmd.Flags().StringVarP(&syncCommitMsg, "message", "m", "", "commit message to use when syncing")
}
