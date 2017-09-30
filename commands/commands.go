package commands

import (
	"fmt"

	"github.com/chasinglogic/dfm/backend"
	"github.com/chasinglogic/dfm/backend/dropbox"
	"github.com/chasinglogic/dfm/backend/git"
	"github.com/chasinglogic/dfm/config"
	"github.com/spf13/cobra"
)

// Global variables to represent root flags available to sub commands
var (
	Verbose bool
	DryRun  bool

	// Whether or not to overwrite existing files when linking
	overwrite bool

	Backend = loadBackend(config.Backend)
)

func init() {
	Root.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	Root.PersistentFlags().BoolVarP(&DryRun, "dry-run", "d", false, "don't make changes just print what would happen")

	Root.AddCommand(Init)
	Root.AddCommand(Add)
	Root.AddCommand(Link)
	Root.AddCommand(List)
	Root.AddCommand(Remove)
	Root.AddCommand(Where)
	Root.AddCommand(Sync)
	Root.AddCommand(Clean)

	for _, c := range Backend.Commands() {
		Root.AddCommand(c)
	}
}

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

// Root is the root dfm command.
var Root = &cobra.Command{
	Use:   "dfm",
	Short: "Manage dotfiles.",
	Long: `Dotfile management written for pair programmers. Examples on getting
started with dfm are avialable at https://github.com/chasinglogic/dfm`,
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		_ = config.SaveConfig()
	},
}

// Execute aliases to running Execute on the root command
func Execute() error {
	return Root.Execute()
}
