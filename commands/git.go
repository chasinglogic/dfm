package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/chasinglogic/dfm/config"
	"github.com/spf13/cobra"
)

func runGitCMD(userDir string, args ...string) error {
	command := exec.Command("git", args...)
	command.Dir = userDir
	out, err := command.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		fmt.Println("ERROR Running Git Command:", "git", args)
	}

	return err
}

// Git runs arbitrary git commands on the current profile
var Git = &cobra.Command{
	Use:                "git",
	Args:               cobra.ArbitraryArgs,
	Short:              "run the given git command on the current profile",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		profile := config.CurrentProfile()

		for _, location := range profile.Locations {
			if err := runGitCMD(location, args...); err != nil {
				fmt.Println("ERROR:", err.Error())
				os.Exit(1)
			}
		}
	},
}
