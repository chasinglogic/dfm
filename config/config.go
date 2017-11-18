// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	// Dir is where dfm will keep internal files and state.
	// TODO: Move the .dfm config json here
	Dir = GetDefaultConfigDir()

	// Backend is a string use for loading the appropriate backend.
	Backend = "git"

	// CurrentProfile is the currently loaded profile.
	CurrentProfile = ""

	// Etc contains additionall information that the various backends can
	// reference.
	Etc map[string]interface{}
)

func init() {
	cfgFile := filepath.Join(os.Getenv("HOME"), ".dfm")
	jsonBytes, err := ioutil.ReadFile(cfgFile)

	var config configJSON

	if err != nil {
		setupWizard()
		return
	}

	err = json.Unmarshal(jsonBytes, &config)
	if err != nil {
		fmt.Println("ERROR:", err.Error())
		os.Exit(1)
	}

	Backend = config.Backend
	CurrentProfile = config.CurrentProfile
	Dir = config.ConfigDir
	Etc = config.Etc

	err = os.MkdirAll(ProfileDir(), os.ModePerm)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}

func getInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return text[:len(text)-1]
}

type configJSON struct {
	Backend        string
	ConfigDir      string
	CurrentProfile string
	Etc            map[string]interface{}
}

func setupWizard() {
	confirmation := "no"

	for strings.HasPrefix(strings.ToLower(confirmation), "n") {
		Backend = getInput(fmt.Sprintf("Which backend would you like to use? [git, dropbox] (Default: %s) ", Backend))
		if Backend == "" {
			Backend = "git"
		}

		fmt.Println("\nNOTE: Some backends will change this automatically and override your setting. (i.e. dropbox)")
		Dir = getInput(fmt.Sprintf("Where would you like to store profiles? (Default: %s) ", Dir))
		if Dir == "" {
			Dir = GetDefaultConfigDir()
		}

		cfg := configJSON{
			Backend:   Backend,
			ConfigDir: Dir,
		}

		jsn, _ := json.MarshalIndent(cfg, "", "\t")
		fmt.Println(string(jsn))

		confirmation = getInput("Does this look correct? Y/n: ")
	}

	err := SaveConfig()
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
}

// SaveConfig should be run after every command in dfm.
func SaveConfig() error {
	config := configJSON{
		Backend:        Backend,
		ConfigDir:      Dir,
		CurrentProfile: CurrentProfile,
		Etc:            Etc,
	}

	jsn, merr := json.MarshalIndent(config, "", "\t")
	if merr != nil {
		return merr
	}

	return ioutil.WriteFile(filepath.Join(os.Getenv("HOME"), ".dfm"), jsn, 0644)
}

// ProfileDir will return the config.Dir joined with profiles.
func ProfileDir() string {
	return filepath.Join(Dir, "profiles")
}

// GetDefaultConfigDir will return the default location to store profiles which
// is $XDG_CONFIG_HOME/dfm or $HOME/.config/dfm
func GetDefaultConfigDir() string {
	xdg := os.Getenv("XDG_CONFIG_HOME")

	if xdg == "" {
		xdg = filepath.Join(os.Getenv("HOME"), ".config")
	}

	return filepath.Join(xdg, "dfm")
}
