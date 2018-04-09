// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package dropbox

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// Backend implements backend.Backend for a dropbox based remote.
type Backend struct{}

func getDropboxDir() string {
	dropboxDir := os.Getenv("DFM_DROPBOX_DIR")

	defaultDir := filepath.Join(os.Getenv("HOME"), "Dropbox")
	_, err := os.Stat(defaultDir)
	if err != nil && os.IsNotExist(err) && dropboxDir == "" {
		fmt.Println("Default Dropbox location found.")
		fmt.Println("Set the DFM_DROPBOX_DIR environment variable.")
	}

	defaultDir = dropboxDir
	return defaultDir
}

func printWarning(userDir string) {
	if !strings.Contains(userDir, getDropboxDir()) {
		fmt.Println(`WARNING: You are using the Dropbox but it doesn't appear
that your Config Directory is not pointed at your Dropbox folder. If you
do not set ConfigDir to somewhere inside your Dropbox folder this
backend will not work.`)
	}
}

// Init determines where the Dropbox folder is and sets up dfm
func (b Backend) Init() error {
	return nil
}

// Sync has nothing to do. Dropbox handles it all
func (b Backend) Sync(userDir string) error {
	printWarning(userDir)
	return nil
}

// NewProfile has nothing to do. Dropbox handles it all
func (b Backend) NewProfile(userDir string) error {
	printWarning(userDir)
	return nil
}

// Commands has nothing to do. Dropbox handles it all
func (b Backend) Commands() []*cobra.Command { return nil }
