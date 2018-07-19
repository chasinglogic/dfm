// Copyright 2018 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.


package commands

import (
	"fmt"
	"os"

	"github.com/chasinglogic/dfm/config"
	"github.com/chasinglogic/dfm/git"
	"github.com/spf13/cobra"
)

// Git runs arbitrary git commands on the current profile
var Git = &cobra.Command{
	Use:                "git",
	Args:               cobra.ArbitraryArgs,
	Short:              "run the given git command on the current profile",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		profile := config.CurrentProfile()

		if err := git.RunGitCMD(profile, args...); err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(1)
		}
	},
}
