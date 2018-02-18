// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chasinglogic/dfm/backend"
	"github.com/chasinglogic/dfm/backend/dropbox"
	"github.com/chasinglogic/dfm/backend/git"
)

func RenameAndLink(userDir, file string) error {
	s := strings.Split(file, string(filepath.Separator))
	newFile := s[len(s)-1]
	newFile = strings.TrimPrefix(newFile, ".")

	// Check if file is in XDG_config.CONFIG_HOME
	xdgConfigHome, _ := filepath.Abs(os.Getenv("XDG_CONFIG_HOME"))
	if s[len(s)-2] == ".config" || s[len(s)-2] == xdgConfigHome {
		newFile = "config" + string(filepath.Separator) + s[len(s)-1]
	}

	newFile = filepath.Join(userDir, newFile)

	err := os.Rename(file, newFile)
	if err != nil {
		fmt.Println("Encountered error:", err)
		fmt.Println("Trying to create intermediate directories...")

		err = os.MkdirAll(filepath.Dir(newFile), 0700)
		if err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(1)
		}

		err = os.Rename(file, newFile)
		if err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(1)
		}
	}

	CreateSymlinks(userDir, os.Getenv("HOME"), false, false, nil)
	return nil
}

// LoadBackend loads the appropriate backend based on string name
func LoadBackend(backendName string) backend.Backend {
	switch backendName {
	case "git":
		return git.Backend{}
	case "dropbox":
		return dropbox.Backend{}
	default:
		fmt.Printf("Backend \"%s\" not found defaulting to git\n.", backendName)
		return git.Backend{}
	}
}
