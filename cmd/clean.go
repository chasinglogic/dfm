package cmd

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/chasinglogic/dfm/logger"
	"github.com/chasinglogic/dfm/profiles"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean dead symlinks created by DFM",
	RunE: func(cmd *cobra.Command, args []string) error {
		home := os.Getenv("HOME")
		if home == "" {
			return errors.New("$HOME is not set!")
		}

		filepath.WalkDir(home, func(path string, d fs.DirEntry, err error) error {
			isSymlink := d.Type()&os.ModeSymlink == os.ModeSymlink
			logger.Debug.Printf("checking if %s is a symlink", path)
			if !isSymlink {
				logger.Debug.Printf("%s is not a symlink", path)
				return nil
			}

			logger.Debug.Printf("%s is a symlink, checking if dfm created it", path)

			realpath, err := os.Readlink(path)
			if err != nil {
				return err
			}

			if !strings.Contains(realpath, profiles.DFMDir) {
				logger.Verbose.Printf("%s is not a dfm symlink", path)
				return nil
			}

			if _, err := os.Stat(realpath); os.IsNotExist(err) {
				logger.Verbose.Printf("%s is a dfm symlink and is dead, removing...", path)
				return os.Remove(path)
			}

			logger.Verbose.Printf("%s is a dfm symlink and is not dead doing nothing", path)
			return nil
		})

		return nil
	},
}
