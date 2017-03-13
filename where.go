package dfm

import (
	"fmt"
	"path/filepath"

	cli "gopkg.in/urfave/cli.v1"
)

// Where simply prints the current profile directory path
func Where(c *cli.Context) error {
	fmt.Println(filepath.Join(getProfileDir(), CONFIG.CurrentProfile))
	return nil
}
