package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chasinglogic/dfm/config"
	"github.com/chasinglogic/dfm/utils"
	"github.com/spf13/cobra"
)

// Add will add the specified profile to the current profile, linking it as
// necessary.
var Add = &cobra.Command{
	Use:   "add",
	Short: "Add a file to the current profile.",
	Run: func(cmd *cobra.Command, args []string) {
		if Verbose {
			fmt.Println("Adding files:", args)
		}

		userDir := filepath.Join(config.ProfileDir(), config.CurrentProfile)

		for _, f := range args {
			file, err := filepath.Abs(f)
			if err != nil {
				fmt.Println("ERROR:", err.Error())
				os.Exit(1)
			}

			if Verbose {
				fmt.Println("Absolute path:", file)
			}

			err = utils.RenameAndLink(userDir, file)
			if err != nil {
				fmt.Println("ERROR:", err.Error())
				os.Exit(1)
			}
		}

		err := Backend.Sync(userDir)
		if err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(1)
		}
	},
}
