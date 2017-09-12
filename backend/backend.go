package backend

import "gopkg.in/urfave/cli.v1"

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
	Commands() []cli.Command
}
