// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package dropbox

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chasinglogic/dfm/config"
	"github.com/spf13/cobra"
)

// Backend implements backend.Backend for a dropbox based remote.
type Backend struct{}

func getDropboxDir() string {
	etc, ok := config.Etc["DROPBOX_DIR"]

	defaultDir := filepath.Join(os.Getenv("HOME"), "Dropbox")
	_, err := os.Stat(defaultDir)
	if err != nil && os.IsNotExist(err) && !ok {
		fmt.Println("Default Dropbox location found.")
		fmt.Println("Set DROPBOX_DIR in your config's etc section.")
		fmt.Println(`Example:
{
    "Etc": {
         "DROPBOX_DIR": "<path to dropbox folder>
    }
}
`)
		os.Exit(1)
	}

	var ed *string

	if ed, ok = etc.(*string); !ok {
		fmt.Println("Error: Etc DROPBOX_DIR is not correct type.")
		fmt.Println("Got:", etc)
		os.Exit(1)
	}

	defaultDir = *ed
	return defaultDir
}

// Init determines where the Dropbox folder is and sets up dfm
func (b Backend) Init() error {
	if !strings.Contains(config.Dir, "Dropbox") {
		fmt.Println(`WARNING: You are using the Dropbox but it doesn't appear
that your Config Directory is not pointed at your Dropbox folder. If you
do not set ConfigDir to somewhere inside your Dropbox folder this
backend will not work.`)
	}

	return nil
}

// Sync has nothing to do. Dropbox handles it all
func (b Backend) Sync(userDir string) error { return nil }

// NewProfile has nothing to do. Dropbox handles it all
func (b Backend) NewProfile(userDir string) error { return nil }

// Commands has nothing to do. Dropbox handles it all
func (b Backend) Commands() []*cobra.Command { return nil }
