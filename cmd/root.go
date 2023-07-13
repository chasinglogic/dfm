package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/chasinglogic/dfm/profiles"
	"github.com/spf13/cobra"
)

var state struct {
	CurrentProfile string `json:"current_profile"`
}

var errProfileNotFound = errors.New("no profile selected")

var statefile = path.Join(profiles.DFMDir, "state.json")
var currentProfile profiles.Profile

func loadCurrentProfile() (profiles.Profile, error) {
	if state.CurrentProfile == "" {
		return profiles.Profile{}, nil
	}

	return loadProfileByName(state.CurrentProfile)
}

func loadProfileByName(name string) (profiles.Profile, error) {
	return profiles.Load(path.Join(profiles.ProfileDir, name))
}

func fail(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = cobra.Command{
	Use:          "dfm",
	Short:        "A dotfile manager for lazy people and pair programmers",
	SilenceUsage: true,
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		fh, err := os.Create(statefile)
		if err != nil {
			return err
		}

		encoder := json.NewEncoder(fh)

		return encoder.Encode(&state)
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile(statefile)
		if err != nil && !os.IsNotExist(err) {
			return err
		} else if os.IsNotExist(err) {
			dfmErr := os.MkdirAll(profiles.DFMDir, 0700)
			if dfmErr != nil {
				return dfmErr
			}

			profilesErr := os.MkdirAll(profiles.ProfileDir, 0700)
			if profilesErr != nil {
				return profilesErr
			}

			modulesErr := os.MkdirAll(profiles.ModulesDir, 0700)
			if modulesErr != nil {
				return modulesErr
			}

			return nil
		}

		err = json.Unmarshal(data, &state)
		if err != nil {
			return err
		}

		currentProfile, err = loadCurrentProfile()
		return err
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(cloneCmd)
	rootCmd.AddCommand(gitCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(linkCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(runHookCmd)
	rootCmd.AddCommand(whereCmd)
	rootCmd.AddCommand(updateCmd)
}

func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		fail(err)
	}
}
