package main

import (
	"os"
	"path/filepath"

	"github.com/chasinglogic/dfm/commands"
	"github.com/urfave/cli"
)

func defaultConfigDir() string {
	xdg := os.Getenv("XDG_CONFIG_HOME")

	if xdg == "" {
		xdg = filepath.Join(os.Getenv("HOME"), ".config")
	}

	return filepath.Join(xdg, "dfm")
}

func main() {
	app := cli.NewApp()
	app.Name = "dfm"
	app.Usage = "Manage dotfiles."
	app.Version = "1.0-dev"
	app.Authors = []cli.Author{
		{
			Name:  "Mathew Robinson",
			Email: "mathew.robinson3114@gmail.com",
		},
		{
			Name:  "Mark Chandler",
			Email: "mark.allen.chandler@gmail.com",
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Use `DIR` for storing dfm configuration",
			Value: defaultConfigDir(),
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Print verbose messaging.",
		},
		cli.BoolFlag{
			Name:  "dry-run",
			Usage: "Don't create symlinks just print what would be done.",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "Create a dotfiles profile from a git repo.",
			Action:  commands.Create,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "alias, a",
					Usage: "Creates `ALIAS` for the profile instead of username",
				},
			},
		},
	}

	app.Run(os.Args)
}
