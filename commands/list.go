package commands

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/chasinglogic/dfm/config"
	"github.com/spf13/cobra"
)

// List will list the available profiles and aliases
var List = &cobra.Command{
	Use:   "list",
	Short: "list available profiles",
	Run: func(cmd *cobra.Command, args []string) {
		files, err := ioutil.ReadDir(config.ProfileDir())
		if err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(1)
		}

		for _, f := range files {
			fmt.Println(f.Name())
		}
	},
}
