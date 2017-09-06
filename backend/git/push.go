package git

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/chasinglogic/dfm/config"
	cli "gopkg.in/urfave/cli.v1"
)

// Push performs a git push origin master in the profile's directory
func Push(c *cli.Context) error {
	profile := c.Args().First()
	if profile == "" {
		profile = config.CONFIG.CurrentProfile
	}

	userDir := filepath.Join(config.ProfileDir(), profile)
	pullCMD := exec.Command("git", "push", "origin", "master")
	pullCMD.Dir = userDir
	output, err := pullCMD.CombinedOutput()

	if err != nil {
		return cli.NewExitError(string(output), 128)
	}

	fmt.Println(string(output))
	return nil
}
