package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/chasinglogic/dfm/config"
	"github.com/urfave/cli"
)

type Backend struct{}

func (b Backend) Init() error {
	if _, err := exec.LookPath("git"); err != nil {
		fmt.Println("ERROR: git is required for this backend.")
		fmt.Println("Please install git then try again.")
		os.Exit(1)
	}

	return nil
}

func (b Backend) Sync(userDir string) error {
	err := runGitCMD(userDir, "git", "add", "--all")
	if err != nil {
		return err
	}

	return runGitCMD(userDir, "git", "commit", "-m", "File added by Dotfile Manager!")
}

func (b Backend) NewProfile(userDir string) error {
	return runGitCMD(userDir, "git", "init")
}

func (b Backend) Commands() []cli.Command {
	return []cli.Command{
		{
			Name:   "pull",
			Usage:  "Run git pull on the current profile.",
			Action: Pull,
		},
		{
			Name:   "push",
			Usage:  "Run git push on the current profile.",
			Action: Push,
		},
		{
			Name:   "clone",
			Usage:  "Run git clone on the current profile.",
			Action: Clone,
		},
		{
			Name:   "status",
			Usage:  "Run git status on the current profile.",
			Action: Status,
		},
		{
			Name:   "commit",
			Usage:  "Run git commit on the current profile.",
			Action: Commit,
		},
		{
			Name:   "git",
			Usage:  "Run arbritrary git commands on the current profile.",
			Action: Git,
		},
	}
}

func runGitCMD(userDir string, args ...string) error {
	command := exec.Command("git", args...)
	command.Dir = userDir

	err := command.Start()
	if err != nil {
		return cli.NewExitError(err.Error(), 128)
	}

	err = command.Wait()
	if err != nil {
		return cli.NewExitError(err.Error(), 128)
	}

	return nil
}

// Git passes directly through and runs the given git command on the current profile.
func Git(c *cli.Context) error {
	cmd := append([]string{c.Args().First()}, c.Args().Tail()...)
	userDir := filepath.Join(config.ProfileDir())
	return runGitCMD(userDir, cmd...)
}
