package dfm

import (
	"fmt"
	"os/exec"
	"path/filepath"

	cli "gopkg.in/urfave/cli.v1"
)

// Status will run git status for the current profile.
func Status(c *cli.Context) error {
	profile := c.Args().First()
	if profile == "" {
		profile = CONFIG.CurrentProfile
	}

	userDir := filepath.Join(getProfileDir(), profile)

	status := exec.Command("git", "status")
	status.Dir = userDir

	output, err := status.CombinedOutput()
	if err != nil {
		return cli.NewExitError(string(output), 128)
	}

	fmt.Println(string(output))
	return nil
}
