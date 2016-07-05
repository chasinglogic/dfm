package common

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Link struct {
	Src  string
	Dest string
}

func GenerateSymlinks(profileDir string) []Link {
	links := []Link{}
	// TODO: Handle the config dir special case
	filepath.Walk(profileDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println("Error walking directory:", err)
				return err
			}

			l := filepath.Join(os.Getenv("HOME"), filepath.Base(path))
			links = append(links, Link{path, l})
			return nil
		})

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
			fmt.Printf("Creating symlink %s -> %s", link.Src, link.Dest)
			if err := os.Symlink(link.Src, link.Dest); err != nil {
				return err
			}
		}

		return nil
	}

	return errors.New("Symlink targets exist.")
}
