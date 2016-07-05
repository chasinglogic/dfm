package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

var VERBOSE = false
var DRY_RUN = false

type Link struct {
	Src  string
	Dest string
}

func (l *Link) String() string {
	return fmt.Sprintf("Link( %s, %s )", l.Src, l.Dest)
}

func setGlobalOptions(c *cli.Context) {
	VERBOSE = c.Bool("verbose")
	DRY_RUN = c.Bool("dry-run")
}

func GenerateSymlinks(profileDir string) []Link {
	links := []Link{}
	// TODO: Handle the config dir special case
	files, err := ioutil.ReadDir(profileDir)
	if err != nil {
		return links
	}

	for _, file := range files {
		if !strings.HasPrefix(file.Name(), ".") {
			ln := Link{
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

func CreateSymlinks(l []Link) error {
	ok := true

	for _, link := range l {
		if _, err := os.Stat(link.Dest); err == nil {
			fmt.Printf("%s already exists.\n", link.Dest)
			ok = false
		}
	}

	if ok {
		for _, link := range l {
			if DRY_RUN || VERBOSE {
				fmt.Printf("Creating symlink %s -> %s\n", link.Src, link.Dest)
			}

			if !DRY_RUN {
				if err := os.Symlink(link.Src, link.Dest); err != nil {
					return err
				}
			}
		}

		return nil
	}

	return cli.NewExitError("Symlink targets exist. Refusing to create a broken state please remove the targets then rerun command.", 68)
}
