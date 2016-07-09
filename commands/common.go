package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

// VERBOSE is used to globally control verbosity
var VERBOSE = false

// DRYRUN is used to globally control whether changes should be made
var DRYRUN = false

// LinkInfo simulates a tuple for our symbolic link
type LinkInfo struct {
	Src  string
	Dest string
}

func (l *LinkInfo) String() string {
	return fmt.Sprintf("Link( %s, %s )", l.Src, l.Dest)
}

func setGlobalOptions(c *cli.Context) {
	VERBOSE = c.Bool("verbose")
	DRYRUN = c.Bool("dry-run")
}

func getProfileDir(c *cli.Context) string {
	return filepath.Join(c.Parent().String("config"), "profiles")
}

func getUser(c *cli.Context) string {
	// This handles the case when create passes us it's context
	if len(strings.Split(c.Args().First(), "/")) > 1 {
		_, user := createURL(strings.Split(c.Args().First(), "/"))
		return user
	}

	return c.Args().First()
}

func generateSymlinks(profileDir string) []LinkInfo {
	links := []LinkInfo{}
	// TODO: Handle the config dir special case
	files, err := ioutil.ReadDir(profileDir)
	if err != nil {
		return links
	}

	for _, file := range files {
		if !strings.HasPrefix(file.Name(), ".") {
			ln := LinkInfo{
				filepath.Join(profileDir, file.Name()),
				filepath.Join(os.Getenv("HOME"), "."+file.Name()),
			}

			if VERBOSE {
				fmt.Printf("Generated symlink %s\n", ln.String())
			}

			links = append(links, ln)
		}
	}

	return links
}
