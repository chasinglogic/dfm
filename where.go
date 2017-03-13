package dfm

import (
	"fmt"

	cli "gopkg.in/urfave/cli.v1"
)

// Where simply prints the current profile directory path
func Where(c *cli.Context) error {
	fmt.Println(getProfileDir())
	return nil
}
