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

	"github.com/chasinglogic/dfm/filemap"
	"github.com/chasinglogic/dfm/git"
	yaml "gopkg.in/yaml.v2"
)

// Module is a profile module, it will be linked as it if was a top level
// profile.
type Module struct {
	Repo         string           `yaml:"repo"`
	Branch       string           `yaml:"branch,omitempty"`
	Link         string           `yaml:"link,omitempty"`
	UserName     string           `yaml"name"`
	UserLocation string           `yaml:"location,omitempty"`
	Mappings     filemap.Mappings `yaml:"mappings,omitempty"`
	PullOnly     bool             `yaml:"pull_only"`
}

// Name of the module, based on what git would use as the directory name from
// the URL. A specific name can be specified in the module configuration
func (m Module) Name() string {
	if m.UserName != "" {
		return m.UserName
	}

	split := strings.Split(m.Repo, "/")
	return split[len(split)-1]
}

func (m Module) Location() string {
	location := filepath.Join(ModuleDir(), m.Name())
	if m.UserLocation != "" {
		location = ExpandFilePath(m.UserLocation)
	}

	if _, err := os.Stat(location); os.IsNotExist(err) {
		branch := m.Branch
		if branch == "" {
			branch = "master"
		}

		err := git.RunGitCMD(
			ModuleDir(),
			"clone",
			"-b", branch,
			m.Repo,
			ExpandFilePath(location),
		)
		if err != nil {
			fmt.Println("ERROR: Unable to clone module:", err)
			os.Exit(1)
		}
	}

	currentBranch, err := git.Branch(location)
	if err != nil {
		fmt.Printf("ERROR: Unable to get current branch for module %s: %s\n", m.Name, err)
		os.Exit(1)
	}

	desiredBranch := m.Branch
	if desiredBranch == "" {
		desiredBranch = "master"
	}

	if currentBranch != desiredBranch {
		err := git.RunGitCMD(
			location,
			"checkout",
			desiredBranch,
		)

		if err != nil {
			fmt.Println("ERROR: Unable to clone module:", err)
			os.Exit(1)
		}
	}

	return location
}

// DFMYml is used for extending and configuring DFM. It is the .dfm.yml file
// found in the root of a profile.
type DFMYml struct {
	Hooks       Hooks            `yaml:"hooks"`
	Mappings    filemap.Mappings `yaml:"mappings"`
	Modules     []Module         `yaml:"modules"`
	SyncModules bool             `yaml:"always_sync_modules"`
}

func (yml DFMYml) Validate() {
	for _, module := range yml.Modules {
		if module.Link != "" &&
			module.Link != "pre" &&
			module.Link != "post" &&
			module.Link != "none" {
			fmt.Println("ERROR: Unknown link value found:", module.Link)
			fmt.Println("This will cause the module to be effectively ignored.")
			fmt.Println("Valid values are: \"pre\", \"post\", and \"none\"")
		}
	}
}

// Return modules which should be linked before the parent profile
func (yml DFMYml) PreLinkModules() []Module {
	var prelinkModules []Module

	for _, module := range yml.Modules {
		if module.Link == "pre" {
			prelinkModules = append(prelinkModules, module)
		}
	}

	return prelinkModules
}

// Return modules which should be linked before the parent profile
func (yml DFMYml) PostLinkModules() []Module {
	var postlinkModules []Module

	for _, module := range yml.Modules {
		if module.Link == "post" || module.Link == "" {
			postlinkModules = append(postlinkModules, module)
		}
	}

	return postlinkModules
}

// LoadDotDFM will load the hooks file for the given Profile
func LoadDotDFM(profileDir string) DFMYml {
	dfmyml, err := ioutil.ReadFile(filepath.Join(profileDir, ".dfm.yml"))
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("ERROR loading .dfm.yml:", err.Error())
		}

		return DFMYml{}
	}

	var yml DFMYml

	err = yaml.Unmarshal(dfmyml, &yml)
	if err != nil {
		fmt.Println("ERROR loading .dfm.yml:", err.Error())
		return yml
	}

	yml.Validate()

	return yml
}

// ModuleDir will return the module directory for the given profile
func ModuleDir() string {
	moduleDir := filepath.Join(GetDefaultConfigDir(), "modules")
	if _, err := os.Stat(moduleDir); os.IsNotExist(err) {
		err = os.MkdirAll(moduleDir, os.ModePerm)
		if err != nil {
			fmt.Println("ERROR: Unable to created module directory:", err)
			return ""
		}
	}

	return moduleDir
}

// ExpandFilePath does bash-esque expansions on a filepath
func ExpandFilePath(path string) string {
	home := os.Getenv("HOME")
	return strings.Replace(path, "~", home, 1)
}
