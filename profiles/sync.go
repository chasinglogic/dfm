package profiles

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/chasinglogic/dfm/logger"
)

type SyncOptions struct {
	CommitMessage string
	SkipModules   bool
}

func (p Profile) Sync(opts SyncOptions) error {
	if err := p.RunHook("before_sync"); err != nil {
		return err
	}

	var err error

	if p.config.PullOnly {
		err = p.pull()
	} else {
		err = p.doSync(opts)
	}

	if err != nil {
		return err
	}

	if err := p.RunHook("after_sync"); err != nil {
		return err
	}

	if opts.SkipModules {
		return nil
	}

	for _, module := range p.modules {
		err = module.Sync(opts)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p Profile) hasOrigin() bool {
	proc := exec.Command("git", "remote", "-v")
	proc.Dir = p.config.Location
	output, err := proc.CombinedOutput()
	if err != nil {
		// Something unexpected happened while running git so let's assume we
		// can't run anymore git commands and skip trying to sync.
		return false
	}

	return strings.Contains(string(output), "origin")
}

func (p Profile) branch() string {
	if p.config.Branch != "" {
		return p.config.Branch
	}

	proc := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	proc.Dir = p.config.Location
	output, err := proc.CombinedOutput()
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(string(output))
}

func (p Profile) pull() error {
	if !p.hasOrigin() {
		logger.Debug.Printf("%s does not have a remote named origin so cannot pull\n", p.config.Location)
		return nil
	}

	logger.Debug.Printf("updating: %s\n", p.config.Location)
	return p.Git("pull", "--rebase", "origin", p.branch())
}

func (p Profile) doSync(opts SyncOptions) error {
	logger.Debug.Printf("syncing: %s\n", p.config.Location)

	dirty := p.IsDirty()
	logger.Debug.Printf("repo is dirty: %t", dirty)

	if dirty {
		// Display the diff to the user.
		p.Git("--no-pager", "diff")

		commitMsg := "Files managed by DFM! https://github.com/chasinglogic/dfm"
		if opts.CommitMessage != "" {
			commitMsg = opts.CommitMessage
		} else if p.config.PromptForCommitMessage {
			fmt.Print("Commit message: ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			if err := scanner.Err(); err != nil {
				return err
			}

			commitMsg = scanner.Text()
		}

		logger.Debug.Printf("commit msg: %s\n", commitMsg)
		p.Git("add", "--all")
		p.Git("commit", "-m", commitMsg)
	}

	err := p.pull()
	if err != nil {
		return err
	}

	if dirty {
		p.Git("push", "origin", p.branch())
	}

	return nil
}

func (p Profile) IsDirty() bool {
	proc := exec.Command("git", "status", "--porcelain")
	proc.Dir = p.config.Location
	output, err := proc.CombinedOutput()
	if err != nil {
		// Something unexpected happened while running git so let's assume we
		// can't run anymore git commands and skip trying to sync.
		return false
	}

	return string(output) != ""
}
