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

// Sync will add and commit all files in the repo then push.
func Sync(workingDir string) error {
	dirty, err := isDirty(workingDir)
	if err != nil {
		return err
	}

	if dirty {
		err := RunGitCMD(workingDir, "add", "--all")
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

		err = RunGitCMD(workingDir, "commit", "-m", msg)
		if err != nil {
			return err
		}
	}

	err = Pull(workingDir)
	if err != nil {
		return err
	}

	if dirty {
		return RunGitCMD(workingDir, "push", "origin", "master")
	}

	return nil
}

// Pull runs git pull --rebase origin master in the given workingDir
func Pull(workingDir string) error {
	err := RunGitCMD(workingDir, "pull", "--rebase", "origin", "master")
	if err != nil {
		return err
	}

	return nil
}

// Init will run git init in the directory
func Init(workingDir string) error {
	return RunGitCMD(workingDir, "init")
}

func isDirty(workingDir string) (bool, error) {
	command := exec.Command("git", "status", "--porcelain")
	fmt.Println(workingDir)
	command.Dir = workingDir
	out, err := command.Output()
	if err != nil {
		fmt.Println("ERROR Running Git Command: git status --porcelain", string(out))
	}

	return string(out) != "", err
}

// RunGitCMD runs git with the given args in workingDir
func RunGitCMD(workingDir string, args ...string) error {
	command := exec.Command("git", args...)
	command.Dir = workingDir
	out, err := command.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		fmt.Println("ERROR Running Git Command:", "git", args)
	}

	return err
}
