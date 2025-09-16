/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"path/filepath"

	"github.com/chasinglogic/dfm/internal/state"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove <PROFILE_NAME>",
	Short:   "Remove a profile",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"rm"},
	RunE: func(cmd *cobra.Command, args []string) error {
		profilesDir, err := state.ProfilesDir()
		if err != nil {
			return err
		}

		profilePath := filepath.Join(profilesDir, args[0])
		return os.RemoveAll(profilePath)
	},
}

func init() {
	RootCmd.AddCommand(removeCmd)
}
