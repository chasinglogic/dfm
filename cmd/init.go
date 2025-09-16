/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"path/filepath"

	"github.com/chasinglogic/dfm/internal/state"
	"github.com/chasinglogic/dfm/internal/utils"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init <PROFILE_NAME>",
	Short: "Create a new profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profilesDir, err := state.ProfilesDir()
		if err != nil {
			return err
		}

		profilePath := filepath.Join(profilesDir, args[0])
		if err := os.MkdirAll(profilePath, 0744); err != nil {
			return err
		}

		return utils.RunIn(profilePath, "git", "init")
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
