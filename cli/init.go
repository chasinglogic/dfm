package cli

import (
	"os"
	"path/filepath"

	"github.com/chasinglogic/dfm/config"
	"github.com/urfave/cli"
)

// Init will create a new profile with the given name.
func Init(c *cli.Context) error {
	profile := c.Args().First()
	if profile == "" {
		return cli.NewExitError("Please specify a profile name.", 1)
	}

	userDir := filepath.Join(config.ProfileDir(), profile)
	err := os.Mkdir(userDir, os.ModePerm)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}

	return Backend.NewProfile(userDir)
}
