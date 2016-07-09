package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

// Link will generate and create the symlinks to the dotfiles in the repo.
func Link(c *cli.Context) error {
	setGlobalOptions(c.Parent())

	userDir := filepath.Join(getProfileDir(c), getUser(c))
	links := generateSymlinks(userDir)
	return createSymlinks(links, c.Bool("overwrite"))
}

func createSymlinks(l []LinkInfo, overwrite bool) error {
	ok := true

	for _, link := range l {
		if _, err := os.Stat(link.Dest); err == nil {
			if overwrite {
				if VERBOSE || DRYRUN {
					fmt.Printf("%s already exists, removing.\n", link.Dest)
				}

				if !DRYRUN {
					if rmerr := os.Remove(link.Dest); rmerr != nil {
						fmt.Printf("Unable to remove %s: %s\n",
							link.Dest,
							rmerr.Error())
					}
				}
			} else {
				fmt.Printf("%s already exists.\n", link.Dest)
				ok = false
			}
		}
	}

	if ok {
		for _, link := range l {
			if DRYRUN || VERBOSE {
				fmt.Printf("Creating symlink %s -> %s\n", link.Src, link.Dest)
			}

			if !DRYRUN {
				if err := os.Symlink(link.Src, link.Dest); err != nil {
					return err
				}
			}
		}

		return nil
	}

	return cli.NewExitError("Symlink targets exist. Refusing to create a broken state please remove the targets then rerun command.", 68)
}
