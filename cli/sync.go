package cli

import (
	"path/filepath"

	"github.com/chasinglogic/dfm/config"
	"gopkg.in/urfave/cli.v1"
)

func Sync(c *cli.Context) error {
	userDir := filepath.Join(filepath.Join(config.ProfileDir(),
		config.CONFIG.CurrentProfile))
	return Backend.Sync(userDir)
}
