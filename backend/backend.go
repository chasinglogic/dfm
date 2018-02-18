// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package backend

import (
	"github.com/spf13/cobra"
)

// Backend represents any syncing service or store that DFM can use.
type Backend interface {
	// This is called on dfm start once the backend to use is determined. Any
	// setup code or checking for available tools should happen in this
	// function
	Init() error

	// Sync the users current profile with the remote backend
	Sync(userDir string) error

	// NewProfile allows the backend to do any necessary setup on a new profile.
	NewProfile(userDir string) error

	// This is used to register "extra" commands to dfm proper from the backend
	// it allows for backends to expose internal functionality but their use
	// should not be required
	Commands() []*cobra.Command
}
