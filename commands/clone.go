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
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/chasinglogic/dfm/config"
	"github.com/spf13/cobra"
)

var (
	link bool
	name string
)

func init() {
	Clone.Flags().BoolVarP(&overwrite, "overwrite", "o", false,
		"whether dfm should remove files that exist where a link should go if --link is given")
	Clone.Flags().BoolVarP(&link, "link", "l", false,
		"whether the profile should be linked after being cloned")
	Clone.Flags().StringVarP(&name, "name", "n", "",
		"name of the profile, this will be automatically computed based on the git url if not provided")
}

// Clone will clone the given git repo to the profiles directory, it optionally
// will call link or use depending on the flag given.
var Clone = &cobra.Command{
	Use:   "clone",
	Short: "git clone an existing profile from `URL`",
	Run: func(cmd *cobra.Command, args []string) {
		url, user := CreateURL(strings.Split(args[0], "/"))
		userDir := filepath.Join(config.ProfileDir(), user)
		if name != "" {
			userDir = filepath.Join(config.ProfileDir(), name)
		}

		if err := CloneRepo(url, userDir); err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(1)
		}

		yml := config.LoadDotDFM(userDir)
		for _, module := range yml.Modules {
			// Location will clone the module if it's not there
			module.Location()
		}

		if link {
			args := []string{"dfm", "link", user}
			if overwrite {
				args = []string{"dfm", "link", "-o", user}
			}

			c := exec.Command(args[0], args[1:]...)
			_, err := c.CombinedOutput()
			if err != nil {
				fmt.Println("ERROR:", err.Error())
				os.Exit(1)
			}
		}
	},
}

// CreateURL will add the missing github.com for the shorthand version of
// links.
func CreateURL(s []string) (string, string) {
	if len(s) == 2 {
		return fmt.Sprintf("https://github.com/%s", strings.Join(s, "/")), s[0]
	}

	return strings.Join(s, "/"), s[len(s)-2]
}

// CloneRepo will git clone the provided url into the appropriate profileDir
func CloneRepo(url, profileDir string) error {
	fmt.Printf("Creating profile in %s\n", profileDir)

	c := exec.Command("git", "clone", url, profileDir)
	_, err := c.CombinedOutput()
	if err != nil && err.Error() == "exit status 128" {
		fmt.Println("Profile exists, perhaps you meant dfm link?")
		os.Exit(128)
	}

	return err
}
