package cli

import (
	"os"

	cli "gopkg.in/urfave/cli.v1"

	"github.com/chasinglogic/dfm"
)

// Added this to make testing easier.
func buildApp() *cli.App {
	app := cli.NewApp()
	app.Name = "dfm"
	app.Usage = "Manage dotfiles."
	app.Version = "1.0"
	app.Authors = []cli.Author{
		{
			Name:  "Mathew Robinson",
			Email: "mathew.robinson3114@gmail.com",
		},
	}

	app.Before = dfm.LoadConfig
	app.After = dfm.SaveConfig

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Use `DIR` for storing dfm configuration and profiles",
			Value: dfm.DefaultConfigDir(),
		},
		cli.BoolFlag{
			Name:  "verbose, vv",
			Usage: "Print verbose messaging.",
		},
		cli.BoolFlag{
			Name:  "dry-run, dr",
			Usage: "Don't create symlinks just print what would be done.",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Add a file to the current profile.",
			Action:  dfm.Add,
		},
		{
			Name:    "clone",
			Aliases: []string{"c"},
			Usage:   "Create a dotfiles profile from a git repo.",
			Action:  dfm.Clone,
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
					Usage: "Links the profile after creation.",
				},
			},
		},
		{
			Name:    "link",
			Aliases: []string{"l"},
			Usage:   "Recreate the links from the dotfiles profile.",
			Action:  dfm.Link,
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
			Action:  dfm.List,
		},
		{
			Name:    "pull",
			Aliases: []string{"pl"},
			Usage:   "Pull the latest version of the profile from origin master.",
			Action:  dfm.Pull,
		},
		{
			Name:    "push",
			Aliases: []string{"ps"},
			Usage:   "Push your local version of the profile to the remote.",
			Action:  dfm.Push,
		},
		{
			Name:        "remove",
			Aliases:     []string{"rm"},
			Usage:       "Remove the profile and all it's symlinks.",
			Description: "Removes the profile and all it's symlinks, if there is another profile on this system we will switch to it. Otherwise will do nothing.",
			Action:      dfm.Remove,
		},
		{
			Name:    "remote",
			Aliases: []string{"re"},
			Usage:   "Will show the remote if given no arguments otherwise will set the remote.",
			Action:  dfm.Remote,
		},

		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "Create a new profile with `NAME`",
			Action:  dfm.Init,
		},
		{
			Name:            "commit",
			Aliases:         []string{"cm"},
			Usage:           "Runs git commit for the profile using `MSG` as the message",
			SkipFlagParsing: true,
			Action:          dfm.Commit,
		},
		{
			Name:    "status",
			Aliases: []string{"st"},
			Usage:   "Runs git status for the current or given profile.",
			Action:  dfm.Status,
		},
		{
			Name:            "git",
			Aliases:         []string{"g"},
			Usage:           "Runs the git command given in the current profile dir directly.",
			SkipFlagParsing: true,
			Action:          dfm.Git,
		},
	}

	return app
}

// Run is the entry point for the app
func Run() int {
	a := buildApp()
	a.Run(os.Args)
	return 0
}
