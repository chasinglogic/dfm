package cli

import (
	"fmt"
	"io/ioutil"

	"github.com/chasinglogic/dfm/config"
	cli "gopkg.in/urfave/cli.v1"
)

// List will list the available profiles and aliases
func List(c *cli.Context) error {
	files, err := ioutil.ReadDir(config.ProfileDir())
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}

	return nil
}
