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

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <FILES>...",
	Short: "Add files to the current dotfile profile",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, err := loadProfile(state.State.CurrentProfile)
		if err != nil {
			return err
		}

		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		for _, file := range args {
			absPath, err := filepath.Abs(file)
			if err != nil {
				return err
			}

			relPath, err := filepath.Rel(home, absPath)
			if err != nil {
				return err
			}

			profilePath := filepath.Join(profile.GetDotfileDirectory(), relPath)

			if err := os.MkdirAll(filepath.Dir(profilePath), 0744); err != nil {
				return err
			}

			if err := os.Rename(absPath, profilePath); err != nil {
				return err
			}
		}

		return profile.Link(false)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
