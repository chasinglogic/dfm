package cmd

import (
	"fmt"
	"os"

	"github.com/chasinglogic/dfm/profiles"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List available profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		profiles, err := os.ReadDir(profiles.ProfileDir)
		if err != nil {
			return err
		}

		for _, profile := range profiles {
			if profile.IsDir() {
				fmt.Println(profile.Name())
			}
		}

		return nil
	},
}
