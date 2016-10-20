package dfm

import (
	"fmt"
	"os/exec"
	"path/filepath"

	cli "gopkg.in/urfave/cli.v1"
)

func Commit(c *cli.Context) error {
	profile := c.Args().First()
	if profile == "" {
		profile = CONFIG.CurrentProfile
	}

	userDir := filepath.Join(getProfileDir(), profile)
	commit := exec.Command("git", "commit", "-m", c.Args().First())
	commit.Dir = userDir
	output, err := commit.CombinedOutput()

	if err != nil {
		return cli.NewExitError(string(output), 128)
	}

	fmt.Println(string(output))
	return nil
}
