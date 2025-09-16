package cmd

import (
	"path/filepath"

	"github.com/chasinglogic/dfm/internal/config"
	"github.com/chasinglogic/dfm/internal/logger"
	"github.com/chasinglogic/dfm/internal/profiles"
	"github.com/chasinglogic/dfm/internal/state"
	"github.com/chasinglogic/dfm/internal/utils"
	"github.com/spf13/cobra"
)

var link bool
var overwrite bool
var profileName string

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone a dotfile repo",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profilesDir, err := state.ProfilesDir()
		if err != nil {
			return err
		}

		repo := args[0]

		if profileName == "" {
			profileName = config.RepoToName(repo)
		}

		profilePath := filepath.Join(profilesDir, profileName)

		logger.Debug().
			Str("profilePath", profilePath).
			Str("repo", repo).
			Msg("cloning repository")
		if err := utils.Run("git", "clone", args[0], profilePath); err != nil {
			return err
		}

		profile, err := profiles.Load(profilePath)
		if err != nil {
			return err
		}

		if link {
			return profile.Link(overwrite)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().StringVarP(
		&profileName,
		"name",
		"n",
		"",
		"Name of the profile, if not provided is derived from the clone URL.",
	)
	cloneCmd.Flags().BoolVarP(
		&link,
		"link",
		"l",
		false,
		"After cloning immediately link the profile",
	)
}
