package cli

import (
	"fmt"
	"os"

	"github.com/chasinglogic/dfm/backend"
	"github.com/chasinglogic/dfm/backend/dropbox"
	"github.com/chasinglogic/dfm/backend/git"
	"github.com/chasinglogic/dfm/config"
	"gopkg.in/urfave/cli.v1"
)

func loadBackend(backendName string) backend.Backend {
	switch backendName {
	case "git":
		return git.Backend{}
	case "dropbox":
		return dropbox.Backend{}
	default:
		fmt.Printf("Backend \"%s\" not found defaulting to git\n.", backendName)
		return git.Backend{}
	}
}

// Added this to make testing easier.
func buildApp() *cli.App {
	app := cli.NewApp()
	app.Name = "dfm"
	app.Usage = "Manage dotfiles."
	app.Version = "3.0.0"
	app.Authors = []cli.Author{
		{
			Name:  "Mathew Robinson",
			Email: "chasinglogic@gmail.com",
		},
	}

	app.Before = func(c *cli.Context) error {
		Verbose = c.Bool("verbose")
		DRYRUN = c.Bool("dry-run")
		return nil
	}

	app.After = config.SaveConfig

	app.Flags = []cli.Flag{
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
			Name:  "license",
			Usage: "Print version and licensing info about dfm.",
			Action: func(c *cli.Context) error {
				fmt.Printf(`Name:
    dfm - Manage dotfiles.

Version:
    %s

License:
    Copyright (C) 2017 Mathew Robinson

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
`, c.App.Version)
				return nil
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Add a file to the current profile.",
			Action:  Add,
		},
		{
			Name:    "link",
			Aliases: []string{"l"},
			Usage:   "Recreate the links from the dotfiles profile.",
			Action:  Link,
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
			Action:  List,
		},
		{
			Name:        "delete",
			Aliases:     []string{"d"},
			Usage:       "Delete the profile and all it's symlinks.",
			Description: "Deletes the profile and all it's symlinks, if there is another profile on this system we will switch to it. Otherwise will do nothing.",
			Action:      Remove,
		},
		{
			Name:        "remove",
			Aliases:     []string{"rm"},
			Usage:       "Remove the file from the profile.",
			Description: "Removes the file pointed at by the link and syncs the current profile. Can restore the file if given the --restore flag.",
			Action:      Rm,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "restore, re",
					Usage: "Restore the file to it's original location instead of deleting it permanently.",
				},
			},
		},
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "Create a new profile with `NAME`",
			Action:  Init,
		},
		{
			Name:    "where",
			Aliases: []string{"w"},
			Usage:   "Prints the CurrentProfile directory, useful for using with other bash commands",
			Action:  Where,
		},
		{
			Name:    "sync",
			Aliases: []string{"s"},
			Usage:   "Sync your config with the configured backend.",
			Action:  Sync,
		},
	}

	err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	Backend = loadBackend(config.CONFIG.Backend)
	err = Backend.Init()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}

	err = os.MkdirAll(config.CONFIG.ConfigDir, os.ModePerm)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}

	app.Commands = append(app.Commands, Backend.Commands()...)

	return app
}

// Run is the entry point for the app
func Run() int {
	a := buildApp()
	a.Run(os.Args)
	return 0
}
