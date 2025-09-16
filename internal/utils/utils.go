package utils

import (
	"os"
	"os/exec"
)

func Run(args ...string) error {
	return RunIn("", args...)
}

func RunIn(dir string, args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	return cmd.Run()
}
