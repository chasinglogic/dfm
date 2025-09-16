/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/chasinglogic/dfm/internal/logger"
	"github.com/chasinglogic/dfm/internal/state"
	"github.com/spf13/cobra"
)

func cleanDeadSymlinks(rootPath, targetDir string) error {
	targetDir = filepath.Clean(targetDir) + string(os.PathSeparator)
	return filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logger.Debug().Str("path", path).Err(err).Msg("error accessing directory")
			return nil
		}

		if d.Type()&fs.ModeSymlink != 0 {
			linkTarget, err := os.Readlink(path)
			if err != nil {
				return err
			}

			// Don't need to make linkTarget absolute because dfm only ever creates
			// absolute links so if we got a relative one it's not one of ours.

			_, err = os.Stat(path)
			if err != nil && os.IsNotExist(err) {
				if strings.HasPrefix(linkTarget, targetDir) {
					fmt.Println("deleting dead link:", path)
					return os.Remove(path)
				}
			}
		}

		return nil
	})
}

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean dead symlinks. Will ignore symlinks unrelated to DFM.",
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		dfmDir, err := state.DfmDir()
		if err != nil {
			return err
		}

		return cleanDeadSymlinks(home, dfmDir)
	},
}

func init() {
	RootCmd.AddCommand(cleanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
