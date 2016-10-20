package commands

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli"
)

// Update performs a git pull origin master in the profile's directory
func Update(c *cli.Context) error {
	userDir := filepath.Join(c.Parent().String("config"), "profiles", getUser(c))
	pullCMD := exec.Command("git", "pull", "origin", "master")
	pullCMD.Dir = userDir
	output, err := pullCMD.CombinedOutput()

	if err != nil {
		return cli.NewExitError(string(output), 128)
	}

	fmt.Println(string(output))
	return nil
}
