package cli

import (
	"fmt"
	"path/filepath"

	"github.com/chasinglogic/dfm/config"
	"github.com/chasinglogic/dfm/utils"
	cli "gopkg.in/urfave/cli.v1"
)

// Add will add the specified profile to the current profile, linking it as
// necessary.
func Add(c *cli.Context) error {
	if Verbose {
		fmt.Println("Adding files:", c.Args())
	}

	userDir := filepath.Join(config.ProfileDir(), config.CONFIG.CurrentProfile)

	for _, f := range c.Args() {
		file, err := filepath.Abs(f)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		if Verbose {
			fmt.Println("Absolute path:", file)
		}

		err = utils.RenameAndLink(userDir, file)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	}

	return Backend.Sync(userDir)
}
