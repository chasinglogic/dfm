package dfm

import (
	"fmt"
	"os/exec"
	"path/filepath"

	cli "gopkg.in/urfave/cli.v1"
)

// Pull performs a git pull origin master in the profile's directory
func Pull(c *cli.Context) error {
	profile := c.Args().First()
	if profile == "" {
		profile = CONFIG.CurrentProfile
	}

	userDir := filepath.Join(getProfileDir(), profile)
	pullCMD := exec.Command("git", "pull", "origin", "master")
	pullCMD.Dir = userDir
	output, err := pullCMD.CombinedOutput()

	if err != nil {
		return cli.NewExitError(string(output), 128)
	}

	fmt.Println(string(output))
	return nil
}
