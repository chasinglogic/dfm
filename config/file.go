package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

// XDG returns the appropriate XDG_CONFIG_HOME directory
func XDG() string {
	xdg := os.Getenv("XDG_CONFIG_HOME")
	if xdg == "" {
		xdg = filepath.Join(os.Getenv("HOME"), ".config")
	}

	return xdg
}

func configFile() string {
	return filepath.Join(GetDefaultConfigDir(), "config.json")
}

func loadConfig() {
	yamlBytes, err := ioutil.ReadFile(configFile())
	if err != nil {
		global = Config{
			CurrentProfileName: "",
		}

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
	return filepath.Join(XDG(), "dfm")
}
