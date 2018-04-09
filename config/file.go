package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

func loadConfig() {
	cfgFile := filepath.Join(os.Getenv("HOME"), ".dfm.yml")
	yamlBytes, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		setupWizard()
		return
	}

	err = yaml.Unmarshal(yamlBytes, &global)
	if err != nil {
		fmt.Println("ERROR:", err.Error())
		os.Exit(1)
	}

	err = os.MkdirAll(ProfileDir(), os.ModePerm)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}

// GetDefaultConfigDir will return the default location to store profiles which
// is $XDG_CONFIG_HOME/dfm or $HOME/.config/dfm
func GetDefaultConfigDir() string {
	xdg := os.Getenv("XDG_CONFIG_HOME")

	if xdg == "" {
		xdg = filepath.Join(os.Getenv("HOME"), ".config")
	}

	return filepath.Join(xdg, "dfm")
}
