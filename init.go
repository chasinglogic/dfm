package dfm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	cli "gopkg.in/urfave/cli.v1"
)

func Init(c *cli.Context) error {
	profile := c.Args().First()
	if profile == "" {
		return cli.NewExitError("Please specify a profile name.", 1)
	}

	userDir := filepath.Join(getProfileDir(), profile)
	err := os.Mkdir(userDir, os.ModeDir)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}

	initCMD := exec.Command("git", "init")
	initCMD.Dir = userDir
	output, err := initCMD.CombinedOutput()

	if err != nil {
		return cli.NewExitError(string(output), 128)
	}

	fmt.Println(string(output))
	return nil
}
