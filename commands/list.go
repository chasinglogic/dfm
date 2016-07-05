package commands

import (
	"fmt"
	"io/ioutil"

	"github.com/urfave/cli"
)

// List will list the available profiles and aliases
func List(c *cli.Context) error {
	files, err := ioutil.ReadDir(getProfileDir(c))
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}

	return nil
}
