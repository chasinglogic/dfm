package commands

import "github.com/urfave/cli"

// Use simply calls link forcing an overwrite.
func Use(c *cli.Context) error {
	c.Set("overwrite", "true")
	return Link(c)
}
