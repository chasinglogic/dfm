/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chasinglogic/dfm/internal/profiles"
	"github.com/chasinglogic/dfm/internal/state"
	"github.com/spf13/cobra"
	"github.com/yarlson/pin"
)

func loadProfile(profileName string) (*profiles.Profile, error) {
	if profileName == "" {
		return nil, errors.New("no current profile is set and no profile name provided")
	}

	profileDir, err := state.ProfilesDir()
	if err != nil {
		return nil, err
	}

	return profiles.Load(filepath.Join(profileDir, profileName))
}

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

		profile, err := loadProfile(profileName)
		if err != nil {
			return err
		}

		p := pin.New(
			fmt.Sprintf("Linking %s...", profileName),
			pin.WithSpinnerColor(pin.ColorCyan),
			pin.WithWriter(os.Stdout),
		)
		cancel := p.Start(context.Background())
		defer cancel()

		err = profile.Link(overwrite)
		if err != nil {
			return err
		}

		p.Stop("Done!")
		return nil
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
