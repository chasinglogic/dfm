package dfm

import (
	"fmt"
	"os/exec"
	"path/filepath"

	cli "gopkg.in/urfave/cli.v1"
)

// Git passes directly through and runs the given git command on the current profile.
func Git(c *cli.Context) error {
	cmd := append([]string{c.Args().First()}, c.Args().Tail()...)
	userDir := filepath.Join(getProfileDir(), CONFIG.CurrentProfile)

	command := exec.Command("git", cmd...)
	command.Dir = userDir

	output, err := command.CombinedOutput()
	if err != nil {
		return cli.NewExitError(string(output), 128)
	}

	fmt.Println(string(output))
	return nil
}
