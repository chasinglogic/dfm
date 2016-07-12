package main

import (
	"os"
	"path/filepath"
	"testing"
)

func testFilesExistence(files ...string) (string, bool) {
	for _, f := range files {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			return f, false
		}
	}

	return "", true
}

func createDefault(t *testing.T) {
	createErr, _ := testCommand("create", "chasinglogic/dotfiles")
	if createErr != nil {
		t.Errorf("Failed with error: %s\n", createErr.Error())
	}
}

func testCommand(args ...string) (error, string) {
	e := os.Setenv("HOME", filepath.Join(os.TempDir(), "dfm_test"))
	if e != nil {
		return e, ""
	}

	app := buildApp()
	profileDir := filepath.Join(defaultConfigDir(), "profiles", "chasinglogic")
	return app.Run(append([]string{"dfm"}, args...)), profileDir
}

func cleanup(t *testing.T) {
	err := os.RemoveAll(os.Getenv("HOME"))
	if err != nil {
		t.Errorf("Error cleaning up: %s\n", err.Error())
	}
}

func TestCreate(t *testing.T) {
	e, profileDir := testCommand("create", "chasinglogic/dotfiles")

	if e != nil {
		t.Errorf("Failed with error: %s\n", e.Error())
	}

	testFilesExistence(filepath.Join(profileDir, "bashrc"),
		filepath.Join(profileDir, "vimrc"),
		filepath.Join(profileDir, "vim"))

	cleanup(t)
}

func TestCreateWithAlias(t *testing.T) {
	e, profileDir := testCommand("create", "-a", "cl", "chasinglogic/dotfiles")

	if e != nil {
		t.Errorf("Failed with error: %s\n", e.Error())
	}

	testFilesExistence(filepath.Join(profileDir, "bashrc"),
		filepath.Join(profileDir, "vimrc"),
		filepath.Join(profileDir, "vim"),
		filepath.Join(defaultConfigDir(), "cl"))

	cleanup(t)
}

func TestLink(t *testing.T) {
	createDefault(t)

	e, _ := testCommand("link", "chasinglogic")
	if e != nil {
		t.Errorf("Failed with error: %s\n", e.Error())
	}

	testFilesExistence(filepath.Join(os.Getenv("HOME"), ".bashrc"),
		filepath.Join(os.Getenv("HOME"), ".vimrc"),
		filepath.Join(os.Getenv("HOME"), ".vim"))

	er, _ := testCommand("link", "-o", "chasinglogic")
	if er != nil {
		t.Errorf("Failed with error: %s\n", e.Error())
	}

	testFilesExistence(filepath.Join(os.Getenv("HOME"), ".bashrc"),
		filepath.Join(os.Getenv("HOME"), ".vimrc"),
		filepath.Join(os.Getenv("HOME"), ".vim"))

	cleanup(t)
}

func TestUse(t *testing.T) {
	createDefault(t)

	e, _ := testCommand("use", "chasinglogic")
	if e != nil {
		t.Errorf("Failed with error: %s\n", e.Error())
	}

	testFilesExistence(filepath.Join(os.Getenv("HOME"), ".bashrc"),
		filepath.Join(os.Getenv("HOME"), ".vimrc"),
		filepath.Join(os.Getenv("HOME"), ".vim"))

	cleanup(t)
}

func TestUpdate(t *testing.T) {
	createDefault(t)

	e, _ := testCommand("update", "chasinglogic")
	if e != nil {
		t.Errorf("Failed to update %s\n", e.Error())
	}

	cleanup(t)
}
