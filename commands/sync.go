package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chasinglogic/dfm/config"
	"github.com/spf13/cobra"
)

// Sync will sync the current profile with the configured backend
var Sync = &cobra.Command{
	Use:   "sync",
	Short: "sync the current profile with the configured backend",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		userDir := filepath.Join(filepath.Join(config.ProfileDir(), config.CurrentProfile))
		if err := Backend.Sync(userDir); err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(1)
		}
	},
}
