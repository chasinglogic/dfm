package hooks

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/chasinglogic/dfm/config"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

type hooks map[string][]string

// DFMYml is used for extending DFM. It is the .dfm.yml file found
// in the root of a profile.
type DFMYml struct {
	Hooks hooks `yaml:"hooks"`
}

func loadHooks() hooks {
	userDir := filepath.Join(config.ProfileDir(), config.CurrentProfile)

	dfmyml, err := ioutil.ReadFile(filepath.Join(userDir, ".dfm.yml"))
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("ERROR loading hooks:", err.Error())
		}

		return nil
	}

	var yml DFMYml

	err = yaml.Unmarshal(dfmyml, &yml)
	if err != nil {
		fmt.Println("ERROR loading hooks:", err.Error())
		return yml.Hooks
	}

	return yml.Hooks
}

// AddHooks will add before and after hooks to the given command.
func AddHooks(command *cobra.Command) *cobra.Command {
	// Store this for later use
	runFunc := command.Run

	command.Run = func(cmd *cobra.Command, args []string) {
		hooks := loadHooks()

		commands, preHooks := hooks["before_"+command.Use]
		if preHooks {
			runCommands(commands)
		}

		// Run the real command
		runFunc(cmd, args)

		commands, postHooks := hooks["after_"+command.Use]
		if postHooks {
			runCommands(commands)
		}
	}

	return command
}
