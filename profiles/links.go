package profiles

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/chasinglogic/dfm/logger"
)

type LinkOptions struct {
	DryRun    bool
	Overwrite bool
}

func (p Profile) Link(opts LinkOptions) error {
	if err := p.RunHook("before_link"); err != nil {
		return err
	}

	for _, module := range p.modules {
		if module.config.LinkMode == "before" {
			err := module.Link(opts)
			if err != nil {
				return err
			}
		}
	}

	if p.config.LinkMode != "none" {
		logger.Debug.Printf("linking profile %s because it's linkmode is %s", p.Name(), p.config.LinkMode)
		err := p.symlinkFiles(opts)
		if err != nil {
			return err
		}
	} else {
		logger.Verbose.Printf("not linking profile %s because it's linkmode is none", p.Name())
	}

	for _, module := range p.modules {
		if module.config.LinkMode == "after" {
			err := module.Link(opts)
			if err != nil {
				return err
			}
		}
	}

	if err := p.RunHook("after_link"); err != nil {
		return err
	}

	return nil
}

var errIsNotSymlink = errors.New("file exists and is not symlink")

func removeIfSymlink(homefile string, opts LinkOptions) error {
	info, err := os.Lstat(homefile)
	if err != nil && !os.IsNotExist(err) {
		return err
	} else if err != nil {
		return nil
	}

	isSymlink := info.Mode()&os.ModeSymlink == os.ModeSymlink
	if !isSymlink && !opts.Overwrite {
		logger.Info.Printf("%s is exists and is not a symlink, cowardly refusing to remove.\n", homefile)
		return errIsNotSymlink
	} else if err == nil && !opts.DryRun {
		logger.Debug.Printf("%s is exists and is a symlink removing...\n", homefile)
		rmErr := os.Remove(homefile)
		if rmErr != nil {
			return rmErr
		}
	}

	return nil
}

func (p Profile) symlinkFiles(opts LinkOptions) error {
	if p.config.LinkMode == "none" {
		return nil
	}

	git := exec.Command("git", "ls-files", "--exclude-standard", "--others", "--cached")
	git.Dir = p.config.Location
	output, err := git.Output()
	if err != nil {
		return err
	}

	files := strings.Split(string(output), "\n")

	logger.Debug.Printf("linking files: %v", files)

	linkAsDir := make(map[string]bool)

linker:
	for _, file := range files {
		logger.Debug.Printf("linking file: %s\n", file)
		if file == "" {
			continue
		}

		dotfile, _ := filepath.Abs(path.Join(p.config.Location, file))

		for _, mapping := range p.config.GetMappings() {
			logger.Debug.Printf("checking if file %s matches: %s\n", file, mapping.Match)
			matches := mapping.Matches(dotfile)
			logger.Debug.Printf("mapping matches file: %t\n", matches)
			if !matches {
				continue
			}

			if mapping.ShouldSkip() {
				logger.Verbose.Printf("skipping file %s because it matches mapping: %s\n", file, mapping.Match)
				continue linker
			}

			if mapping.ShouldLinkAsDir() {
				dir := path.Dir(file)
				if _, ok := linkAsDir[dir]; !ok {
					logger.Verbose.Printf("linking directory %s because it matches mapping: %s\n", dir, mapping.Match)
					err := p.doLink(dir, opts)
					if err != nil {
						return err
					}

				} else {
					logger.Debug.Printf("directory %s has already been linked\n", dir)
				}

				continue linker
			}
		}

		err := p.doLink(file, opts)
		if err != nil {
			return err
		}

	}

	return nil
}

func (p Profile) doLink(fileOrDir string, opts LinkOptions) error {
	dotfile, _ := filepath.Abs(path.Join(p.config.Location, fileOrDir))
	homefile, _ := filepath.Abs(path.Join(os.Getenv("HOME"), fileOrDir))
	err := removeIfSymlink(homefile, opts)
	if err != nil && err == errIsNotSymlink {
		return nil
	} else if err != nil {
		return err
	}

	logger.Debug.Printf("creating new symlink: %s", homefile)
	if opts.DryRun {
		fmt.Println(homefile, "->", dotfile)
		return nil
	}

	if _, err := os.Stat(path.Dir(homefile)); os.IsNotExist(err) {
		err := os.MkdirAll(path.Dir(homefile), 0700)
		if err != nil {
			return err
		}
	}

	logger.Verbose.Println(homefile, "->", dotfile)
	linkErr := os.Symlink(dotfile, homefile)
	if linkErr != nil {
		return linkErr
	}

	return nil
}
