package cli

import (
	"fmt"
	"os"

	"github.com/chasinglogic/dfm/backend"
	"github.com/chasinglogic/dfm/backend/git"
	"github.com/chasinglogic/dfm/config"
	"github.com/urfave/cli"
)

func loadBackend(backendName string) backend.Backend {
	switch backendName {
	case "git":
		return git.Backend{}
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
	app.Version = "1.1"
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

	e := config.LoadConfig()
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}

	Backend = loadBackend(config.CONFIG.Backend)
	Backend.Init()
	app.Commands = append(app.Commands, Backend.Commands()...)

	return app
}

// Run is the entry point for the app
func Run() int {
	a := buildApp()
	a.Run(os.Args)
	return 0
}
