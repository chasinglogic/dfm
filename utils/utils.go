package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/urfave/cli.v1"
)

func RenameAndLink(userDir, file string) error {
	s := strings.Split(file, string(filepath.Separator))
	newFile := s[len(s)-1]
	newFile = strings.TrimPrefix(newFile, ".")

	// Check if file is in XDG_config.CONFIG_HOME
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

	CreateSymlinks(userDir, os.Getenv("HOME"), false, false)
	return nil
}
