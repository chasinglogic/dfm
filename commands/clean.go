package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Clean looks for any symlinks which are broken and cleans them up. This can happen
// if files are removed from the profile directory manually instead of through dfm.
var Clean = &cobra.Command{
	Use:   "clean",
	Short: "clean dead symlinks",
	Long:  "Clean looks for any symlinks which are broken and cleans them up. This can happen if files are removed from the profile directory manually instead of through dfm.",
	Run: func(cmd *cobra.Command, args []string) {
		homeDir := os.Getenv("HOME")
		if err := cleanDeadLinks(homeDir); err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(1)
		}

		xdg := os.Getenv("XDG_CONFIG_HOME")
		if xdg == "" {
			xdg = filepath.Join(os.Getenv("HOME"), ".config")
		}

		if err := cleanDeadLinks(xdg); err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(1)
		}
	},
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
