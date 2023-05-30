package profiles

import (
	"os"
	"os/exec"

	"github.com/chasinglogic/dfm/logger"
)

const (
	LinkModePre  = "pre"
	LinkModePost = "post"
	LinkModeNone = "none"
)

type Profile struct {
	config  ProfileConfig
	modules []Profile
}

func (p Profile) Name() string {
	return p.config.Name
}

func (p Profile) Git(args ...string) error {
	git := exec.Command("git", args...)
	git.Dir = p.config.Location
	git.Stdin = os.Stdin
	git.Stdout = os.Stdout
	git.Stderr = os.Stderr
	return git.Run()
}

// Create will ensure that the given path exists and run git init there
// returning a Profile object representing this location.
func Create(where string) (Profile, error) {
	return Profile{}, nil
}

// Load will load the .dfm.yml file in where and create a Profile object
// representing this location.
func FromConfig(cfg ProfileConfig) Profile {
	logger.Debug.Printf("creating profile from config: %s", cfg.Name)

	p := Profile{
		config:  cfg,
		modules: make([]Profile, len(cfg.Modules)),
	}

	for idx, moduleConfig := range cfg.Modules {
		p.modules[idx] = FromConfig(moduleConfig)
		err := p.modules[idx].ensureExists()
		if err != nil {
			panic(err)
		}
	}

	return p
}

func (p Profile) RunHook(name string) error {
	return p.config.Hooks.RunHook(name, p.config.Location, false)
}

func (p Profile) Where() string {
	p.RunHook("before_where")
	defer p.RunHook("after_where")
	return p.config.Location
}

func (p Profile) ensureExists() error {
	_, err := os.Stat(p.config.Location)
	if !os.IsNotExist(err) {
		return err
	}

	git := exec.Command("git", "clone", p.config.GetRepo(), p.config.Location)
	git.Stdin = os.Stdin
	git.Stdout = os.Stdout
	git.Stderr = os.Stderr
	return git.Run()
}
