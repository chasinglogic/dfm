/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/chasinglogic/dfm/internal/state"
	"github.com/chasinglogic/dfm/internal/utils"
	"github.com/spf13/cobra"
)

var gitCmd = &cobra.Command{
	Use:                "git",
	Short:              "Run the given git command on the current profile",
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: true,
	Aliases:            []string{"g"},
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, err := loadProfile(state.State.CurrentProfile)
		if err != nil {
			return err
		}

		args = append([]string{"git"}, args...)
		return utils.RunIn(profile.GetLocation(), args...)
	},
}

func init() {
	RootCmd.AddCommand(gitCmd)
}
