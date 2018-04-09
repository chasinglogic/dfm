// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/chasinglogic/dfm/dotfiles"
	yaml "gopkg.in/yaml.v2"
)

func init() {
	if checkForOldConfig() {
		upgradeConfig()
		return
	}

	loadConfig()
}

type Config struct {
	// Dir is where dfm will keep internal files and state.
	// TODO: Move the .dfm config json here
	Dir string `yaml:"dir"`

	// Backend is the backend used for new profiles.
	DefaultBackend string `yaml:"default_backend"`

	// Profiles is the available profiles and their
	// associated configs
	Profiles []dotfiles.Profile `yaml:"profiles"`

	// CurrentProfileName is the currently loaded profile.
	CurrentProfileName string `yaml:"current_profile"`
}

func (c Config) CurrentProfile() dotfiles.Profile {
	return c.GetProfileByName(c.CurrentProfileName)
}

func (c Config) GetProfileByName(name string) dotfiles.Profile {
	for _, profile := range c.Profiles {
		if profile.Name == name {
			return profile
		}
	}

	return dotfiles.Profile{}
}

func (c Config) AddProfile(profile dotfiles.Profile) {
	c.Profiles = append(c.Profiles, profile)
}

var global Config

func SetCurrentProfile(profile dotfiles.Profile) {
	global.CurrentProfileName = profile.Name
	SaveConfig()
}

func CurrentProfile() dotfiles.Profile {
	return global.CurrentProfile()
}

func GetProfileByName(name string) dotfiles.Profile {
	return global.GetProfileByName(name)
}

func AddProfile(profile dotfiles.Profile) {
	global.AddProfile(profile)
	SaveConfig()
}

func AvailableProfiles() []dotfiles.Profile {
	return global.Profiles
}

func Dir() string {
	return global.Dir
}

func DefaultBackend() string {
	return global.DefaultBackend
}

// SaveConfig should be run after every command in dfm.
func SaveConfig() error {
	yml, merr := yaml.Marshal(global)
	if merr != nil {
		fmt.Println(merr)
		return merr
	}

	return ioutil.WriteFile(filepath.Join(os.Getenv("HOME"), ".dfm.yml"), yml, 0644)
}

// ProfileDir will return the config.Dir joined with profiles.
func ProfileDir() string {
	return filepath.Join(global.Dir, "profiles")
}
