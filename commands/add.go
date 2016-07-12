package commands

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

func Add(c *cli.Context) error {
	setGlobalOptions(c)

	userDir := filepath.Join(getProfileDir(c), getUser(c))

	for _, f := range c.Args() {
		file, err := filepath.Abs(f)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		nodot := strings.TrimPrefix(f, ".")
		newFile := filepath.Join(userDir, nodot)

		move(file, newFile)
		os.Link(newFile, file)
	}

	addCMD := exec.Command("git", "add", "--all")
	commitCMD := exec.Command("git", "commit", "-m", "File added by Dotfile Manager! :D")

	addCMD.Dir = userDir
	commitCMD.Dir = userDir

	output, err := addCMD.CombinedOutput()
	if err != nil {
		return cli.NewExitError(string(output), 128)
	}

	output, err = commitCMD.CombinedOutput()
	if err != nil {
		return cli.NewExitError(string(output), 128)
	}

	return nil
}

func move(oldfile, newfile string) error {
	return os.Rename(oldfile, newfile)
}
