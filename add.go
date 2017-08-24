package dfm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	cli "gopkg.in/urfave/cli.v1"
)

func renameAndLink(userDir, file string) error {
	s := strings.Split(file, string(filepath.Separator))
	newFile := s[len(s)-1]
	newFile = strings.TrimPrefix(newFile, ".")

	// Check if file is in XDG_CONFIG_HOME
	xdgConfigHome, _ := filepath.Abs(os.Getenv("XDG_CONFIG_HOME"))
	if s[len(s)-2] == ".config" || s[len(s)-2] == xdgConfigHome {
		newFile = "config" + string(filepath.Separator) + s[len(s)-1]
	}

	newFile = filepath.Join(userDir, newFile)
	err := os.Rename(file, newFile)
	if err != nil {
		fmt.Println("Encountered error:", err)
		fmt.Println("Trying to create intermediate directories...")

		err = os.MkdirAll(filepath.Dir(newFile), 0700)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		err = os.Rename(file, newFile)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	}

	return os.Link(newFile, file)
}

// Add will add the specified profile to the current profile, linking it as
// necessary.
func Add(c *cli.Context) error {
	if CONFIG.Verbose {
		fmt.Println("Adding files:", c.Args())
	}

	userDir := filepath.Join(getProfileDir(), CONFIG.CurrentProfile)

	for _, f := range c.Args() {
		file, err := filepath.Abs(f)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		if CONFIG.Verbose {
			fmt.Println("Absolute path:", file)
		}

		err = renameAndLink(userDir, file)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	}

	addCMD := exec.Command("git", "add", "--all")
	commitCMD := exec.Command("git", "commit", "-m", "File added by Dotfile Manager!")

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
