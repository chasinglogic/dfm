package git

import (
	"fmt"
	"os/exec"
	"path/filepath"

	cli "gopkg.in/urfave/cli.v1"
)

// Remote will run git remote -v if given no arguments otherwise will set the
// remote for origin
func Remote(c *cli.Context) error {
	remote := c.Args().First()
	userDir := filepath.Join(getProfileDir(), CONFIG.CurrentProfile)

	// No args means run git remote -v
	if remote == "" {
		cmd := exec.Command("git", "remote", "-v")
		cmd.Dir = userDir

		output, err := cmd.CombinedOutput()
		if err != nil {
			return cli.NewExitError(string(output), 128)
		}

		fmt.Println(string(output))
		return nil
	}

	cmd := exec.Command("git", "remote", "add", "origin", remote)
	cmd.Dir = userDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Might be the remote already exists so try again with set-url
		cmd = exec.Command("git", "remote", "set-url", "origin", remote)
		cmd.Dir = userDir

		output, err = cmd.CombinedOutput()
		if err != nil {
			return cli.NewExitError(err.Error(), 128)
		}
	}

	fmt.Println(string(output))
	return nil
}
