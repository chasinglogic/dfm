package state

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type appState struct {
	CurrentProfile string
}

var State *appState

func DfmDir() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	d := filepath.Join(cacheDir, "dfm")
	return d, os.MkdirAll(d, 0744)
}

func stateFile() (string, error) {
	d, err := DfmDir()
	return filepath.Join(d, "state.json"), err
}

func subDir(name string) (string, error) {
	d, err := DfmDir()
	if err != nil {
		return "", err
	}

	subDir := filepath.Join(d, name)
	return subDir, os.MkdirAll(subDir, 0744)
}

func ModulesDir() (string, error) {
	return subDir("modules")
}

func ProfilesDir() (string, error) {
	return subDir("profiles")
}

func Load() error {
	if State != nil {
		return nil
	}

	State = &appState{}

	file, err := stateFile()
	if err != nil {
		return err
	}

	content, err := os.ReadFile(file)
	if err != nil && os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	return json.Unmarshal(content, State)
}

func Save() error {
	content, err := json.Marshal(State)
	if err != nil {
		return err
	}

	file, err := stateFile()
	if err != nil {
		return err
	}

	return os.WriteFile(file, content, 0644)
}
