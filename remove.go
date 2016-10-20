package dfm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cli "gopkg.in/urfave/cli.v1"
)

func Remove(c *cli.Context) error {
	profile := c.Args().First()
	if profile == "" {
		return cli.NewExitError("Please specify a profile.", 1)
	}

	userDir := filepath.Join(getProfileDir(), profile)
	links := GenerateSymlinks(userDir)

	rmerr := os.RemoveAll(userDir)
	if rmerr != nil {
		return rmerr
	}

	if CONFIG.Verbose {
		fmt.Println("Removed profile directory:", userDir)
	}

	return RemoveSymlinks(links, c.Args().First())
}

func RemoveSymlinks(l []LinkInfo, username string) error {
	for _, link := range l {
		// Check if the link is still valid after removing the profile, and if
		// so just verify thta it doesn't contain the username of the profile
		// we're deleting. If the profile we're removing is the one that was
		// currently in use then both conditions should be true.
		if path, err := os.Readlink(link.Dest); err != nil ||
			strings.Contains(path, username) {

			if DRYRUN || CONFIG.Verbose {
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
