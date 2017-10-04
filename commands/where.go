// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package commands

import (
	"fmt"
	"path/filepath"

	"github.com/chasinglogic/dfm/config"
	"github.com/spf13/cobra"
)

// Where simply prints the current profile directory path
var Where = &cobra.Command{
	Use:   "where",
	Short: "prints the current profile directory path",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(filepath.Join(config.ProfileDir(), config.CurrentProfile))
	},
}
