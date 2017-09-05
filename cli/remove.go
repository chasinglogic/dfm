package dfm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	cli "gopkg.in/urfave/cli.v1"
)

// Remove will remove the specified profile
func Remove(c *cli.Context) error {
	profile := c.Args().First()
	if profile == "" {
		return cli.NewExitError("Please specify a profile.", 1)
	}

	userDir := filepath.Join(getProfileDir(), profile)
	links := GenerateSymlinks(userDir, os.Getenv("HOME"))

	rmerr := os.RemoveAll(userDir)
	if rmerr != nil {
		return rmerr
	}

	if CONFIG.Verbose {
		fmt.Println("Removed profile directory:", userDir)
	}

	return RemoveSymlinks(links, c.Args().First())
}

// GenerateSymlinks will create the symlinks so we know what they were supposed
// to be prior to removing the profile.
func GenerateSymlinks(profileDir, target string) *[]LinkInfo {
	var lnks []LinkInfo

	files, err := ioutil.ReadDir(profileDir)
	if err != nil {
		fmt.Println(err)
		return &lnks
	}

	for _, file := range files {
		// Handle the XDG_CONFIG_HOME special case
		if file.Name() == "config" && file.IsDir() {
			xdg := os.Getenv("XDG_CONFIG_HOME")
			if xdg == "" {
				xdg = filepath.Join(os.Getenv("HOME"), ".config")
			}

			lnks = append(lnks,
				*GenerateSymlinks(filepath.Join(profileDir, file.Name()), xdg)...)
		}

		lnks = append(lnks, *GenerateSymlink(profileDir, target, file))
	}

	return &lnks
}

// RemoveSymlinks will remove all of the symlinks after removing a profile it
// will first Check if the link is still valid after removing the profile, and
// if so just verify that it doesn't contain the username of the profile
// we're deleting. If the profile we're removing is the one that was currently
// in use then both conditions should be true.
func RemoveSymlinks(l *[]LinkInfo, username string) error {
	for _, link := range *l {
		// Check if the link is still valid after removing the profile, and if
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
