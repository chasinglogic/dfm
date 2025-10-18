package profiles

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/chasinglogic/dfm/internal/config"
	"github.com/chasinglogic/dfm/internal/logger"
	"github.com/chasinglogic/dfm/internal/mapping"
	"github.com/chasinglogic/dfm/internal/utils"
	"github.com/chzyer/readline"
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
	if err := p.RunHook("pre_link"); err != nil {
		return err
	}

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

	logger.Debug().Interface("config", p.config).Msg("starting link")

	err = filepath.WalkDir(
		p.config.GetDotfileDirectory(),
		func(path string, d fs.DirEntry, err error) error {
			if d == nil {
				return nil
			}

			if d.IsDir() {
				if filepath.Base(path) == ".git" {
					logger.Debug().
						Str("path", path).
						Msg("skipping because it is the git directory")
					return filepath.SkipDir
				}

				return nil
			}

			if filepath.Base(path) == ".dfm.yml" {
				logger.Debug().
					Str("path", path).
					Msg("skipping because it is the dfm config file")

				return nil
			}

			for _, m := range p.config.Mappings {
				if m.IsMatch(path) {
					logger.Debug().
						Str("mapping", m.String()).
						Str("path", path).
						Msg("matched mapping")

					return p.handleMapping(
						overwrite,
						path,
						m,
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

	return p.RunHook("post_link")
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
	rel, err := filepath.Rel(p.config.GetDotfileDirectory(), path)
	if err != nil {
		return err
	}

	targetPath := filepath.Join(target, rel)

	logger.Debug().
		Str("relativePath", rel).
		Str("targetDirectory", target).
		Str("path", path).
		Str("targetPath", targetPath).
		Msg("link")
	if err := deleteIfExists(overwrite, targetPath); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0744); err != nil {
		return err
	}

	return os.Symlink(path, targetPath)
}

func deleteIfExists(overwrite bool, path string) error {
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("refusing to remove a directory: %s", path)
	}

	if info.Mode().IsRegular() && !overwrite {
		return fmt.Errorf(
			"refusing to remove %s because it is a regular file and --overwrite not provided",
			path,
		)
	}

	return os.Remove(path)
}

func (p *Profile) GetLocation() string {
	return p.config.Location
}

func (p *Profile) GetDotfileDirectory() string {
	return p.config.GetDotfileDirectory()
}

func (p *Profile) isDirty() bool {
	buf := bytes.NewBuffer([]byte{})

	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = p.config.Location
	cmd.Stdout = buf
	cmd.Stderr = buf
	_ = cmd.Run()

	return buf.String() != ""
}

func (p *Profile) Sync() error {
	if err := p.RunHook("pre_sync"); err != nil {
		return err
	}

	fmt.Println("Syncing", p.GetLocation())
	if !p.isDirty() || p.config.PullOnly {
		if err := utils.RunIn(p.config.Location, "git", "pull", "--ff-only"); err != nil {
			return err
		}
	} else {
		if err := utils.RunIn(p.config.Location, "git", "diff"); err != nil {
			return err
		}

		commitMessage := "Dotfiles managed by DFM!"
		if p.config.PromptForCommitMessage {
			rl, err := readline.New("Commit message: ")
			if err != nil {
				panic(err)
			}

			commitMessage, err = rl.Readline()
			if err != nil {
				return err
			}

			_ = rl.Close()
		}

		cmds := [][]string{
			{"git", "add", "--all"},
			{"git", "commit", "--message", commitMessage},
			{"git", "pull", "--rebase"},
			{"git", "push"},
		}

		for _, cmd := range cmds {
			if err := utils.RunIn(p.config.Location, cmd...); err != nil {
				return err
			}
		}
	}
	fmt.Println("")

	for _, module := range p.modules {
		if err := module.Sync(); err != nil {
			return err
		}
	}

	return p.RunHook("post_sync")
}

func (p *Profile) RunHook(hookName string) error {
	return p.config.Hooks.Execute(p.config.Location, hookName)
}
