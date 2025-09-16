/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/chasinglogic/dfm/internal/state"
	"github.com/spf13/cobra"
)

var whereCmd = &cobra.Command{
	Use:     "where",
	Aliases: []string{"w"},
	Short:   "Prints the location of the current dotfile profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, err := loadProfile(state.State.CurrentProfile)
		if err != nil {
			return err
		}

		fmt.Println(profile.GetLocation())
		return nil
	},
}

func init() {
	RootCmd.AddCommand(whereCmd)
}
