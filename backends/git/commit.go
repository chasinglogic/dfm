package dfm

import (
	"os"
	"os/exec"
	"path/filepath"

	cli "gopkg.in/urfave/cli.v1"
)

// Commit takes the first argument as a commit message and runs git commit in
// the current profile directory.
func Commit(c *cli.Context) error {
	profile := CONFIG.CurrentProfile
	userDir := filepath.Join(getProfileDir(), profile)

	args := append([]string{"commit", c.Args().First()}, c.Args().Tail()...)
	commit := exec.Command("git", args...)
	commit.Dir = userDir
	commit.Stdin = os.Stdin
	commit.Stdout = os.Stdout
	commit.Stderr = os.Stderr

	err := commit.Start()
	if err != nil {
		return cli.NewExitError(err.Error(), 128)

	}

	err = commit.Wait()
	if err != nil {
		return cli.NewExitError(err.Error(), 128)
	}

	return nil
}
