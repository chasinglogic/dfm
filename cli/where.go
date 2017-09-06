package cli

import (
	"fmt"
	"path/filepath"

	"github.com/chasinglogic/dfm/config"
	cli "gopkg.in/urfave/cli.v1"
)

// Where simply prints the current profile directory path
func Where(c *cli.Context) error {
	fmt.Println(filepath.Join(config.ProfileDir(), config.CONFIG.CurrentProfile))
	return nil
}
