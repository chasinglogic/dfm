package cli

import (
	"path/filepath"

	"github.com/chasinglogic/dfm/config"
	"github.com/urfave/cli"
)

func Sync(c *cli.Context) error {
	userDir := filepath.Join(filepath.Join(config.ProfileDir(),
		config.CONFIG.CurrentProfile))
	return Backend.Sync(userDir)
}
