package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chasinglogic/dfm/config"
	"github.com/chasinglogic/dfm/utils"
	cli "gopkg.in/urfave/cli.v1"
)

// Link will generate and create the symlinks to the dotfiles in the repo.
func Link(c *cli.Context) error {
	userDir := filepath.Join(config.ProfileDir(), c.Args().First())
	fmt.Println("Linking profile", c.Args().First())

	if err := utils.CreateSymlinks(userDir, os.Getenv("HOME"), DRYRUN, c.Bool("overwrite")); err != nil {
		return cli.NewExitError(err.Error(), 2)
	}

	config.CONFIG.CurrentProfile = c.Args().First()
	return nil
}
