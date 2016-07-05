package commands

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/urfave/cli"
)

// List will list the available profiles and aliases
func List(c *cli.Context) error {
	profileDir := filepath.Join(c.Parent().String("config"), "profiles")
	files, err := ioutil.ReadDir(profileDir)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}

	return nil
}
