package dfm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	cli "gopkg.in/urfave/cli.v1"
)

// Link will generate and create the symlinks to the dotfiles in the repo.
func Link(c *cli.Context) error {
	userDir := filepath.Join(getProfileDir(), c.Args().First())
	links := GenerateSymlinks(userDir)
	if err := CreateSymlinks(links, c.Bool("overwrite")); err != nil {
		return cli.NewExitError(err.Error(), 2)
	}

	CONFIG.CurrentProfile = c.Args().First()
	return nil
}

// LinkInfo simulates a tuple for our symbolic link
type LinkInfo struct {
	Src  string
	Dest string
}

func (l *LinkInfo) String() string {
	return fmt.Sprintf("%s -> %s", l.Dest, l.Src)
}

func GenerateSymlinks(profileDir string) []LinkInfo {
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

			if DRYRUN {
				fmt.Printf("Generated symlink %s\n", ln.String())
			}

			links = append(links, ln)
		}
	}

	return links
}

func CheckLinks(l []LinkInfo, overwrite bool, files chan LinkInfo) {
	for _, link := range l {
		if _, err := os.Stat(link.Dest); err == nil {
			if overwrite {
				if CONFIG.Verbose || DRYRUN {
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
				files <- LinkInfo{}
				continue
			}
		}

		files <- link
	}
}

func CreateLinks(numOfLinks int, files chan LinkInfo, errors chan error) {
	count := 0

	for count < numOfLinks {
		link := <-files
		// Means we had an error
		if link.Src == "" {
			continue
		}

		if CONFIG.Verbose || DRYRUN {
			fmt.Println("Creating symlink", link)
		}

		if !DRYRUN {
			if err := os.Symlink(link.Src, link.Dest); err != nil {
				fmt.Println(err)
			}
		}

		count++
	}
}

func CreateSymlinks(l []LinkInfo, overwrite bool) error {
	files := make(chan LinkInfo, 3)
	errors := make(chan error)

	go CheckLinks(l, overwrite, files)
	go CreateLinks(len(l), files, errors)

	return nil
}
