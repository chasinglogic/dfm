package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

var verbose = false

type Link struct {
	Src  string
	Dest string
}

func setVerbosity(c *cli.Context) {
	verbose = c.Bool("verbose")
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
			l := filepath.Join(os.Getenv("HOME"), "."+file.Name())
			if verbose {
				fmt.Printf("Geneated symlink %s\n", l)
			}
			links = append(links, Link{profileDir + file.Name(), l})
		}
	}

	os.Exit(0)
	return links
}

func CreateSymlinks(l []Link) error {
	ok := true

	for _, link := range l {
		if _, err := os.Stat(link.Dest); err == nil {
			fmt.Printf("%s already exists. Please remove and rerun this command.\n",
				link.Dest)
			ok = false
		}
	}

	if ok {
		for _, link := range l {
			fmt.Printf("Creating symlink %s -> %s\n", link.Src, link.Dest)
			if err := os.Symlink(link.Src, link.Dest); err != nil {
				return err
			}
		}

		return nil
	}

	return errors.New("Symlink targets exist.")
}
