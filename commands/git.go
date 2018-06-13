package commands

import (
	"fmt"
	"os"

	"github.com/chasinglogic/dfm/config"
	"github.com/chasinglogic/dfm/git"
	"github.com/spf13/cobra"
)

// Git runs arbitrary git commands on the current profile
var Git = &cobra.Command{
	Use:                "git",
	Args:               cobra.ArbitraryArgs,
	Short:              "run the given git command on the current profile",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		profile := config.CurrentProfile()

		if err := git.RunGitCMD(profile, args...); err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(1)
		}
	},
}
