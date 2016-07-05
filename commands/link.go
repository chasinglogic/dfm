package commands

import (
	"path/filepath"

	"github.com/urfave/cli"
)

// Link will generate and create the symlinks to the dotfiles in the repo.
func Link(c *cli.Context) error {
	setGlobalOptions(c.Parent())

	userDir := filepath.Join(c.Parent().String("config"), "profiles", getUser(c))
	links := generateSymlinks(userDir)
	return createSymlinks(links, c.Bool("overwrite"))
}
