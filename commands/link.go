package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chasinglogic/dfm/config"
	"github.com/chasinglogic/dfm/utils"
	"github.com/spf13/cobra"
)

func init() {
	Link.Flags().BoolVarP(&overwrite, "overwrite", "o", false,
		"whether dfm should remove files that exist where a link should go")
}

// Link will generate and create the symlinks to the dotfiles in the repo.
var Link = &cobra.Command{
	Use:   "link",
	Short: "link the profile with `NAME`",
	Long:  "will generate and create the symlinks to the dotfiles in the profile",
	Run: func(cmd *cobra.Command, args []string) {
		userDir := filepath.Join(config.ProfileDir(), args[0])
		fmt.Println("Linking profile", args[0])

		if err := utils.CreateSymlinks(userDir, os.Getenv("HOME"), DryRun, overwrite); err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(1)
		}

		config.CurrentProfile = args[0]
	},
}
