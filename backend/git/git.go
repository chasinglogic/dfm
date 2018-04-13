// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package git

import (
	"fmt"
	"os"
	"os/exec"
)

func getUserMsg() string {
	return os.Getenv("DFM_GIT_COMMIT_MSG")
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
	dirty, err := isDirty(userDir)
	if err != nil {
		return err
	}

	if dirty {
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
	}

	err = runGitCMD(userDir, "pull", "--rebase", "origin", "master")
	if err != nil {
		return err
	}

	if dirty {
		return runGitCMD(userDir, "push", "origin", "master")
	}

	return nil
}

// NewProfile will run git init in the directory
func (b Backend) NewProfile(userDir string) error {
	return runGitCMD(userDir, "init")
}

func isDirty(userDir string) (bool, error) {
	command := exec.Command("git", "status", "--porcelain")
	command.Dir = userDir
	out, err := command.Output()
	return string(out) != "", err
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
