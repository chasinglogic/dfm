package backends

import "github.com/urfave/cli"

// DFMBackend represents any syncing service or store that DFM can use.
type DFMBackend interface {
	// This is called on dfm start once the backend to use is determined. Any
	// setup code or checking for available tools should happen in this
	// function
	Init() error

	// Sync the users current profile with the remote backend
	Sync() error

	// Get a new profile from the remote backend
	Clone(profileName string) error

	// Remove a profile
	RemoveProfile(profileName string) error

	// Get the full path for the given profile name
	ProfileDir(profileName string) (string, error)

	// Add a file to the current profile
	AddFile(file string) error

	// Get all profiles
	GetProfiles() []string

	// This is used to register "extra" commands to dfm proper from the backend
	// it allows for backends to expose internal functionality but their use
	// should not be required
	Commands() []cli.Command
}
