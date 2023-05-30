package cmd

import (
	"os"
	"path"

	"github.com/chasinglogic/dfm/profiles"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create and initialise a new dotfile profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		profilePath := path.Join(profiles.ProfileDir, name)
		if err := os.Mkdir(profilePath, 0700); err != nil {
			return err
		}

		cfg := profiles.DefaultConfig(profilePath)
		profile := profiles.FromConfig(cfg)

		if err := profile.Git("init"); err != nil {
			return err
		}

		fh, err := os.Create(path.Join(profilePath, ".dfm.yml"))
		if err != nil {
			return err
		}

		encoder := yaml.NewEncoder(fh)
		if err := encoder.Encode(&cfg); err != nil {
			return err
		}

		if err := profile.Git("add", ".dfm.yml"); err != nil {
			return err
		}

		if err := profile.Git("commit", "-m", "add default .dfm.yml"); err != nil {
			return err
		}

		state.CurrentProfile = profile.Name()

		// TODO: Add some helpful output about next steps.
		return nil
	},
}
