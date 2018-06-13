// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package commands

import (
	"fmt"

	"github.com/chasinglogic/dfm/config"
	"github.com/spf13/cobra"
)

// Where simply prints the current profile directory path
var Where = &cobra.Command{
	Use:   "where",
	Short: "prints the first location for the current profile directory path",
	Run: func(cmd *cobra.Command, args []string) {
		profile := config.CurrentProfile()
		fmt.Println(profile)
	},
}
