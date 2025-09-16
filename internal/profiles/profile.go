package profiles

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/chasinglogic/dfm/internal/config"
	"github.com/chasinglogic/dfm/internal/mapping"
	"github.com/chasinglogic/dfm/internal/utils"
)

type Profile struct {
	config  *config.Config
	modules []*Profile
}

func New(config *config.Config) (*Profile, error) {
	profile := Profile{
		config:  config,
		modules: make([]*Profile, len(config.Modules)),
	}

	return &profile, profile.loadModules()
}

func Load(profilePath string) (*Profile, error) {
	config, err := config.Load(filepath.Join(profilePath, ".dfm.yml"))
	if err != nil {
		return nil, err
	}

	return New(config)
}

func (p *Profile) loadModules() error {
	for idx, moduleConfig := range p.config.Modules {
		module, err := New(&moduleConfig)
		if err != nil {
			return err
		}

		if err := module.ensureDownloaded(); err != nil {
			return err
		}

		p.modules[idx] = module
	}

	return nil
}

func (p *Profile) ensureDownloaded() error {
	if _, err := os.Stat(p.config.Location); os.IsNotExist(err) {
		return utils.Run("git", "clone", p.config.Repo, p.config.Location)
	}

	return nil
}

func (p *Profile) Link(overwrite bool) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	for _, profile := range p.modules {
		if profile.config.LinkMode == "pre" {
			if err := profile.Link(overwrite); err != nil {
				return err
			}
		}
	}

	err = filepath.WalkDir(
		p.config.Location,
		func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				if filepath.Base(path) == ".git" {
					return filepath.SkipDir
				}

				return nil
			}

			if filepath.Base(path) == ".dfm.yml" {
				return nil
			}

			for _, mapper := range p.config.Mappings {
				if mapper.IsMatch(path) {
					return p.handleMapping(
						overwrite,
						path,
						mapper,
						home,
					)
				}
			}

			return p.linkTo(overwrite, path, home)
		},
	)
	if err != nil {
		return err
	}

	for _, profile := range p.modules {
		if profile.config.LinkMode != "pre" {
			if err := profile.Link(overwrite); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Profile) handleMapping(
	overwrite bool,
	path string,
	m *mapping.Mapping,
	home string,
) error {
	switch m.Action() {
	case mapping.ActionSkip, mapping.ActionNone:
		return nil
	case mapping.ActionLinkAsDir:
		// TODO: would be nice if we could skip the dir but I don't see an obvious way
		// to make that not suck from a code maintenance perspective yet.
		return p.linkTo(overwrite, filepath.Dir(path), home)
	case mapping.ActionTranslate:
		return p.linkTo(overwrite, path, m.Dest)
	default:
		return fmt.Errorf("unhandled map action: %s", m.Action())
	}
}

func (p *Profile) linkTo(overwrite bool, path, target string) error {
	rel, err := filepath.Rel(p.config.Location, path)
	if err != nil {
		return err
	}

	targetPath := filepath.Join(target, rel)

	fmt.Println("linking", path, "->", targetPath)
	if err := deleteIfExists(overwrite, targetPath); err != nil {
		return err
	}

	return nil
	// return os.Symlink(path, targetPath)
}

func deleteIfExists(overwrite bool, path string) error {
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	if info.IsDir() {
		return errors.New("refusing to remove a directory")
	}

	if info.Mode().IsRegular() && !overwrite {
		return fmt.Errorf(
			"refusing to remove %s because it is a regular file and --overwrite not provided",
			path,
		)
	}

	// return os.Remove(path)
	return nil
}
