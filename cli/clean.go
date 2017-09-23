package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	cli "gopkg.in/urfave/cli.v1"
)

// Clean looks for any symlinks which are broken and cleans them up. This can happen
// if files are removed from the profile directory manually instead of through dfm.
func Clean(c *cli.Context) error {
	homeDir := os.Getenv("HOME")
	if err := cleanDeadLinks(homeDir); err != nil {
		return err
	}

	xdg := os.Getenv("XDG_CONFIG_HOME")
	if xdg == "" {
		xdg = filepath.Join(os.Getenv("HOME"), ".config")
	}

	return cleanDeadLinks(xdg)
}

func cleanDeadLinks(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.Mode()&os.ModeSymlink == os.ModeSymlink {
			p, err := os.Readlink(filepath.Join(dir, file.Name()))
			if err != nil {
				return err
			}

			_, err = os.Stat(p)
			if err == nil {
				continue
			}

			err = os.Remove(filepath.Join(dir, file.Name()))
			if err != nil {
				fmt.Println("Removing", file.Name())
				return err
			}
		}
	}

	return nil
}
