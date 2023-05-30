package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var whereCmd = &cobra.Command{
	Use:     "where",
	Aliases: []string{"w"},
	Short:   "print the location of the currently active dfm profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(currentProfile.Where())
		return nil
	},
}
