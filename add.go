package dfm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	cli "gopkg.in/urfave/cli.v1"
)

func Add(c *cli.Context) error {
	config, cerr := loadConfig(c.Parent())
	if cerr != nil {
		return cli.NewExitError(cerr.Error(), 3)
	}

	if VERBOSE {
		fmt.Println("Adding files:", c.Args())
	}

	userDir := filepath.Join(getProfileDir(c), config.CurrentProfile)

	for _, f := range c.Args() {
		file, err := filepath.Abs(f)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		if VERBOSE {
			fmt.Println("Absolute path:", file)
		}

		nodot := strings.TrimPrefix(f, ".")
		newFile := filepath.Join(userDir, nodot)

		err = move(file, newFile)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		err = os.Link(newFile, file)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
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
