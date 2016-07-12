package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

func Remove(c *cli.Context) error {
	config, cerr := loadConfig(c.Parent())
	if cerr != nil {
		return cli.NewExitError(cerr.Error(), 3)
	}

	profile := getUser(c)
	if profile == "" {
		profile = config.CurrentProfile
	}

	userDir := filepath.Join(getProfileDir(c), profile)
	links := generateSymlinks(userDir)

	rmerr := os.RemoveAll(userDir)
	if rmerr != nil {
		return rmerr
	}

	if VERBOSE {
		fmt.Println("Removed profile directory:", userDir)
	}

	return removeSymlinks(links, getUser(c))
}

func removeSymlinks(l []LinkInfo, username string) error {
	for _, link := range l {
		// Check if the link is still valid after removing the profile, and if
		// so just verify thta it doesn't contain the username of the profile
		// we're deleting. If the profile we're removing is the one that was
		// currently in use then both conditions should be true.
		if path, err := os.Readlink(link.Dest); err != nil ||
			strings.Contains(path, username) {

			if DRYRUN || VERBOSE {
				fmt.Printf("Removing symlink %s\n", link.Dest)
			}

			if !DRYRUN {
				if err := os.Remove(link.Dest); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
