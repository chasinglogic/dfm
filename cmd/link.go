/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"path/filepath"

	"github.com/chasinglogic/dfm/internal/profiles"
	"github.com/chasinglogic/dfm/internal/state"
	"github.com/spf13/cobra"
)

// linkCmd represents the link command
var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "A brief description of your command",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var profileName string
		if len(args) > 0 {
			profileName = args[0]
		} else {
			profileName = state.State.CurrentProfile
		}

		profileDir, err := state.ProfilesDir()
		if err != nil {
			return err
		}

		profile, err := profiles.Load(filepath.Join(profileDir, profileName))
		if err != nil {
			return err
		}

		return profile.Link(overwrite)
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)
	linkCmd.Flags().BoolVarP(
		&overwrite,
		"overwrite",
		"o",
		false,
		"Delete existing files if they conflict with a link target",
	)
}
