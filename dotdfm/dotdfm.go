// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package dotdfm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/chasinglogic/dfm/filemap"
	"github.com/chasinglogic/dfm/hooks"
	yaml "gopkg.in/yaml.v2"
)

// DFMYml is used for extending DFM. It is the .dfm.yml file found
// in the root of a profile.
type DFMYml struct {
	Hooks    hooks.Hooks      `yaml:"hooks"`
	Mappings filemap.Mappings `yaml:"mappings"`
}

// LoadDotDFM will load the hooks file for the current Profile
func LoadDotDFM(userDir string) DFMYml {
	dfmyml, err := ioutil.ReadFile(filepath.Join(userDir, ".dfm.yml"))
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("ERROR loading hooks:", err.Error())
		}

		return DFMYml{}
	}

	var yml DFMYml

	err = yaml.Unmarshal(dfmyml, &yml)
	if err != nil {
		fmt.Println("ERROR loading hooks:", err.Error())
		return yml
	}

	return yml
}
