/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:     "status",
	Short:   "Print the git status of the current dotfile profile",
	Aliases: []string{"st"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return gitCmd.RunE(cmd, []string{"status"})
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
}
