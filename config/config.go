// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Init() {
	if checkForOldConfig() {
		upgradeConfig()
		return
	}

	loadConfig()
}

type Config struct {
	// CurrentProfileName is the currently loaded profile.
	CurrentProfileName string `yaml:"current_profile"`
}

var global Config

func CurrentProfile() string {
	profile := filepath.Join(ProfileDir(), global.CurrentProfileName)
	if profile == ProfileDir() {
		files, err := ioutil.ReadDir(ProfileDir())
		if err != nil {
			fmt.Println("ERROR: Unable to load profiles:", err)
			os.Exit(1)
		}

		if len(files) == 0 {
			fmt.Println("ERROR: No dfm profiles found")
			os.Exit(1)
		}

		for _, file := range files {
			if strings.HasPrefix(file.Name(), ".") {
				continue
			}

			SetCurrentProfile(file.Name())
			return filepath.Join(ProfileDir(), file.Name())
		}
	}

	return profile
}

func GetProfileByName(name string) string {
	return filepath.Join(ProfileDir(), name)
}

func AddProfile(profile string) error {
	return os.Mkdir(GetProfileByName(profile), os.ModePerm)
}

func AvailableProfiles() []string {
	files, err := ioutil.ReadDir(ProfileDir())
	if err != nil {
		fmt.Println("ERROR: Unable to read config dir:", err)
	}

	var profiles []string

	for _, f := range files {
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}

		profiles = append(profiles, f.Name())
	}

	return profiles
}

func Dir() string {
	return GetDefaultConfigDir()
}

// SetCurrentProfile will save the config.json
func SetCurrentProfile(profile string) error {
	global.CurrentProfileName = profile
	jsn, merr := json.Marshal(global)
	if merr != nil {
		fmt.Println(merr)
		return merr
	}

	return ioutil.WriteFile(configFile(), jsn, 0644)
}

// ProfileDir will return the config.Dir joined with profiles.
func ProfileDir() string {
	return filepath.Join(GetDefaultConfigDir(), "profiles")
}

// ProfileName returns the name of the profile from it's full path
func ProfileName(path string) string {
	split := strings.Split(path, string(filepath.Separator))
	return split[len(split)-1]
}
