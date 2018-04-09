package dotfiles

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/chasinglogic/dfm/backend"
	"github.com/chasinglogic/dfm/backend/dropbox"
	"github.com/chasinglogic/dfm/backend/git"
	"github.com/chasinglogic/dfm/filemap"
	"github.com/chasinglogic/dfm/linking"
)

// loadBackend loads the appropriate backend based on string name
func loadBackend(backendName string) backend.Backend {
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

// Profile represents a DFM dotfile profile
type Profile struct {
	Name      string   `yaml:"name"`
	Backend   string   `yaml:"backend"`
	Locations []string `yaml:"locations"`
}

// Sync this profile using the configured backend
func (p Profile) Sync() error {
	backend := loadBackend(p.Backend)
	var errs []error

	for _, location := range p.Locations {
		err := backend.Sync(location)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return collectErrors(errs)
}

// Link this profile to $HOME
func (p Profile) Link(target string, mappings filemap.Mappings, config linking.Config) error {
	var errs []error

	for _, location := range p.Locations {
		err := linking.CreateSymlinks(location, target, config, mappings)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return collectErrors(errs)
}

func collectErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	var errMsgs []string

	for _, err := range errs {
		errMsgs = append(errMsgs, err.Error())
	}

	return errors.New(strings.Join(errMsgs, "\n"))
}

func (p Profile) Init() error {
	for _, location := range p.Locations {
		err := os.Mkdir(location, os.ModePerm)
		if err != nil {
			return err
		}

		backend := loadBackend(p.Backend)
		err = backend.NewProfile(location)
		if err != nil {
			return err
		}
	}

	return nil
}
