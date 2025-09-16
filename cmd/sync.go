/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/chasinglogic/dfm/internal/state"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:     "sync",
	Short:   "Sync your dotfiles with git",
	Aliases: []string{"s"},
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, err := loadProfile(state.State.CurrentProfile)
		if err != nil {
			return err
		}

		return profile.Sync()
	},
}

func init() {
	RootCmd.AddCommand(syncCmd)
}
