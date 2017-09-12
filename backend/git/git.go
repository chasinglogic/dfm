package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/chasinglogic/dfm/config"
	"gopkg.in/urfave/cli.v1"
)

func getUserMsg() string {
	etc, ok := config.CONFIG.Etc["DFM_GIT_COMMIT_MSG"]
	if !ok {
		return ""
	}

	msg, _ := etc.(*string)
	return *msg
}

// Backend implements backend.Backend for a git based remote.
type Backend struct{}

// Init checks for the existence of git as it's a requirement for this backend.
func (b Backend) Init() error {
	if _, err := exec.LookPath("git"); err != nil {
		fmt.Println("ERROR: git is required for this backend.")
		fmt.Println("Please install git then try again.")
		os.Exit(1)
	}

	return nil
}

// Sync will add and commit all files in the repo then push.
func (b Backend) Sync(userDir string) error {
	err := runGitCMD(userDir, "git", "add", "--all")
	if err != nil {
		return err
	}

	msg := "Files managed by DFM! https://github.com/chasinglogic/dfm"
	if userMsg := os.Getenv("DFM_GIT_COMMIT_MSG"); userMsg != "" {
		msg = userMsg
	}

	if userMsg := getUserMsg(); userMsg != "" {
		msg = userMsg
	}

	err = runGitCMD(userDir, "git", "commit", "-m", msg)
	if err != nil {
		return err
	}

	return runGitCMD(userDir, "git", "push", "origin", "master")
}

// NewProfile will run git init in the directory
func (b Backend) NewProfile(userDir string) error {
	return runGitCMD(userDir, "git", "init")
}

// Commands adds some git specific funtionality
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
