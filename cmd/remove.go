package cmd

import (
	"errors"
	"os"
	"path"

	"github.com/chasinglogic/dfm/profiles"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "Remove a profile from your system",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profile := path.Join(profiles.ProfileDir, args[0])
		if _, err := os.Stat(profile); os.IsNotExist(err) {
			return errors.New("ERROR: profile does not exist.")
		} else if err != nil {
			return err
		}

		return os.RemoveAll(profile)
	},
}
