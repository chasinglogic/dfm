// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

func Init() {
	if checkForOldConfig() {
		upgradeConfig()
		return
	}

	loadConfig()
	err := saveConfig()
	if err != nil {
		fmt.Println("ERROR: Unable to save config:", err)
		os.Exit(1)
	}
}

type Config struct {
	// Dir is where dfm will keep internal files and state.
	Dir string `yaml:"dir"`
	// CurrentProfileName is the currently loaded profile.
	CurrentProfileName string `yaml:"current_profile"`
}

func (c Config) ProfileDir() string {
	return filepath.Join(c.Dir, "profiles")
}

func (c Config) CurrentProfile() string {
	profile := filepath.Join(c.ProfileDir(), c.CurrentProfileName)
	if profile == c.ProfileDir() {
		files, err := ioutil.ReadDir(c.ProfileDir())
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

			return filepath.Join(c.ProfileDir(), file.Name())
		}
	}

	return profile
}

func (c Config) GetProfileByName(name string) string {
	return filepath.Join(c.ProfileDir(), name)
}

func (c Config) AddProfile(name string) error {
	return os.Mkdir(c.GetProfileByName(name), os.ModePerm)
}

func (c Config) AvailableProfiles() []string {
	files, err := ioutil.ReadDir(c.ProfileDir())
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

var global Config

func CurrentProfile() string {
	return global.CurrentProfile()
}

func GetProfileByName(name string) string {
	return global.GetProfileByName(name)
}

func AddProfile(profile string) error {
	return global.AddProfile(profile)
}

func AvailableProfiles() []string {
	return global.AvailableProfiles()
}

func Dir() string {
	return global.Dir
}

// saveConfig should be run after every command in dfm.
func saveConfig() error {
	yml, merr := yaml.Marshal(global)
	if merr != nil {
		fmt.Println(merr)
		return merr
	}

	err := os.MkdirAll(global.Dir, os.ModePerm)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configFile(), yml, 0644)
}

// ProfileDir will return the config.Dir joined with profiles.
func ProfileDir() string {
	return global.ProfileDir()
}
