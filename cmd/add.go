package cmd

import (
	"os"
	"path"
	"path/filepath"

	"github.com/chasinglogic/dfm/profiles"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add files to the currently active profile",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, file := range args {
			homefile, err := filepath.Abs(file)
			if err != nil {
				return err
			}

			homedir := os.Getenv("HOME")

			relpath, err := filepath.Rel(homedir, homefile)
			if err != nil {
				return err
			}

			dotfile := path.Join(currentProfile.Where(), relpath)

			precedingDirectories := path.Dir(dotfile)
			if err := os.MkdirAll(precedingDirectories, 0700); err != nil {
				return err
			}

			if err := os.Rename(homefile, dotfile); err != nil {
				return err
			}
		}

		return currentProfile.Link(profiles.LinkOptions{})
	},
}
