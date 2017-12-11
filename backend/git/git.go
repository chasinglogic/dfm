// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/chasinglogic/dfm/config"
	"github.com/spf13/cobra"
)

func getUserMsg() string {
	etc, ok := config.Etc["DFM_GIT_COMMIT_MSG"]
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
	err := runGitCMD(userDir, "add", "--all")
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

	err = runGitCMD(userDir, "commit", "-m", msg)
	if err != nil {
		return err
	}

	err = runGitCMD(userDir, "pull", "--rebase", "origin", "master")
	if err != nil {
		return err
	}

	return runGitCMD(userDir, "push", "origin", "master")
}

// NewProfile will run git init in the directory
func (b Backend) NewProfile(userDir string) error {
	return runGitCMD(userDir, "init")
}

// Commands adds some git specific funtionality
func (b Backend) Commands() []*cobra.Command {
	return []*cobra.Command{
		Git,
		Clone,
	}
}

func runGitCMD(userDir string, args ...string) error {
	command := exec.Command("git", args...)
	command.Dir = userDir
	out, err := command.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		fmt.Println("ERROR Running Git Command:", "git", args)
	}

	return err
}

// Git runs arbitrary git commands on the current profile
var Git = &cobra.Command{
	Use:                "git",
	Args:               cobra.ArbitraryArgs,
	Short:              "run the given git command on the current profile",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		userDir := filepath.Join(config.ProfileDir(), config.CurrentProfile)
		if err := runGitCMD(userDir, args...); err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(1)
		}
	},
}
