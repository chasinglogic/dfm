package cmd

import (
	"github.com/chasinglogic/dfm/profiles"
	"github.com/spf13/cobra"
)

var (
	linkDryRun    = false
	linkOverwrite = false
)

var linkCmd = &cobra.Command{
	Use:     "link",
	Aliases: []string{"l"},
	Short:   "Create links for a profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		profile := currentProfile
		if len(args) == 1 {
			profile, err = loadProfileByName(args[0])
			if err != nil {
				return err
			}
		}

		err = profile.Link(profiles.LinkOptions{
			DryRun:    linkDryRun,
			Overwrite: linkOverwrite,
		})
		if err != nil {
			return err
		}

		state.CurrentProfile = profile.Name()
		return nil
	},
}

func init() {
	linkCmd.Flags().BoolVarP(&linkOverwrite, "overwrite", "o", false, "remove regular files if they exist at the link path")
	linkCmd.Flags().BoolVarP(&linkDryRun, "dry-run", "d", false, "print what would happen instead of creating symlinks")
}
