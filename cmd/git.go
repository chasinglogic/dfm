package cmd

import (
	"github.com/spf13/cobra"
)

var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Run git commands in the currently active dotfile profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		return currentProfile.Git(args...)
	},
}
