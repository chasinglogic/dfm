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

// Added this to make testing easier.
func buildApp() *cli.App {
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
				cli.BoolFlag{
					Name:  "overwrite, o",
					Usage: "Overwrites existing files when creating links.",
				},
				cli.BoolFlag{
					Name:  "link, l",
					Usage: "Links the profile after creation. Use --use instead if there is already an active profile.",
				},
				cli.BoolFlag{
					Name:  "use, u",
					Usage: "Switches to the newly created profile after downloaded.",
				},
			},
		},
		{
			Name:    "link",
			Aliases: []string{"l"},
			Usage:   "Recreate the links from the dotfiles profile.",
			Action:  commands.Link,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "overwrite, o",
					Usage: "Overwrites existing files when creating links.",
				},
			},
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "List available profiles",
			Action:  commands.List,
		},
		{
			Name:    "update",
			Aliases: []string{"up"},
			Usage:   "Pull the latest version of the profile from origin master.",
			Action:  commands.Update,
		},
		{
			Name:    "use",
			Aliases: []string{"u"},
			Usage:   "Switch to the profile",
			Action:  commands.Use,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:   "overwrite, o",
					Usage:  "Overwrites existing files when creating links.",
					Hidden: true,
				},
			},
		},
		{
			Name:        "remove",
			Aliases:     []string{"rm"},
			Usage:       "Remove the profile and all it's symlinks.",
			Description: "Removeds the profile and all it's symlinks, if there is another profile on this system we will switch to it. Otherwise will do nothing.",
			Action:      commands.Remove,
		},
	}

	return app
}

func main() {
	app := buildApp()
	app.Run(os.Args)
}
