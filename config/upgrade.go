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


package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func getInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return text[:len(text)-1]
}

func checkForOldConfig() bool {
	_, err := getOldConfig()
	return err != os.ErrNotExist
}

func getOldConfig() (string, error) {
	cfgFile := filepath.Join(os.Getenv("HOME"), ".dfm")
	_, err := os.Stat(cfgFile)
	if err == nil {
		return cfgFile, err
	}

	ymlFile := filepath.Join(os.Getenv("HOME"), ".dfm.yml")
	_, ymlErr := os.Open(ymlFile)
	if ymlErr == nil {
		return ymlFile, ymlErr
	}

	return "", os.ErrNotExist
}

func upgradeConfig() {
	fmt.Println("It looks like you have an old style DFM config.")
	fmt.Println("We've removed global configs in favor of per-profile configs")
	fmt.Println("See https://github.com/chasinglogic/dfm for more info.")
	ans := getInput("Should we remove the old config? ")
	if strings.HasPrefix(strings.ToLower(ans), "y") {
		cfgFile, _ := getOldConfig()
		fmt.Println("Removing...")
		err := os.Remove(cfgFile)
		if err != nil {
			fmt.Println("ERROR:", err)
		}
	}
}
