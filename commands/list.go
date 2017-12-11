// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

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
