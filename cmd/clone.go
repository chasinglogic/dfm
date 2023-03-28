package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	link = false
)

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone a git repo as a dotfile profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("TODO")
		return nil
	},
}

func init() {
	cloneCmd.Flags().BoolVarP(&link, "link", "l", false, "immediately link the profile after cloning")
	cloneCmd.Flags().BoolVarP(&linkOverwrite, "overwrite", "o", false, "remove regular files if they exist at the link path")
	cloneCmd.Flags().BoolVarP(&linkDryRun, "dry-run", "d", false, "print what would happen instead of creating symlinks")
}
