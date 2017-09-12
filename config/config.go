package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/urfave/cli.v1"
)

var DefaultConfigDir = GetDefaultConfigDir()

// Config holds some basic information about the state of dfm.
type Config struct {
	ConfigDir      string
	CurrentProfile string
	Backend        string
	Etc            map[string]interface{} `json:",omitempty"`
}

var CONFIG Config

func getInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return text[:len(text)-1]
}

func setupWizard() Config {
	confirmation := "no"

	cfg := Config{
		Backend:   "git",
		ConfigDir: DefaultConfigDir,
	}

	for strings.HasPrefix(strings.ToLower(confirmation), "n") {
		cfg.Backend = getInput(fmt.Sprintf("Which backend would you like to use? [git, dropbox] (Default: %s) ", cfg.Backend))
		if cfg.Backend == "" {
			cfg.Backend = "git"
		}

		fmt.Println("\nNOTE: Some backends will change this automatically and override your setting. (i.e. dropbox)")
		cfg.ConfigDir = getInput(fmt.Sprintf("Where would you like to store profiles? (Default: %s) ", cfg.ConfigDir))
		fmt.Println(cfg.ConfigDir)
		if cfg.ConfigDir == "" {
			cfg.ConfigDir = DefaultConfigDir
		}

		jsn, _ := json.MarshalIndent(cfg, "", "\t")
		fmt.Println(string(jsn))

		confirmation = getInput("Does this look correct? Y/n: ")
	}

	return cfg
}

// LoadConfig runs as a hook when dfm runs to load the config object from
// ~/.dfm
func LoadConfig() error {
	cfgFile := filepath.Join(os.Getenv("HOME"), ".dfm")
	configJSON, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		CONFIG = setupWizard()
	} else {
		err = json.Unmarshal(configJSON, &CONFIG)
		if err != nil {
			return err
		}
	}

	return nil
}

// SaveConfig will save the config to the configDir/config.json
func SaveConfig(c *cli.Context) error {
	jsn, merr := json.MarshalIndent(CONFIG, "", "\t")
	if merr != nil {
		return merr
	}

	return ioutil.WriteFile(filepath.Join(os.Getenv("HOME"), ".dfm"), jsn, 0644)
}

func ProfileDir() string {
	return filepath.Join(CONFIG.ConfigDir, "profiles")
}

func CurrentProfile() string {
	return CONFIG.CurrentProfile
}

// DefaultConfigDir will return the default location to store profiles which is
// $XDG_CONFIG_HOME/dfm or $HOME/.config/dfm
func GetDefaultConfigDir() string {
	xdg := os.Getenv("XDG_CONFIG_HOME")

	if xdg == "" {
		xdg = filepath.Join(os.Getenv("HOME"), ".config")
	}

	return filepath.Join(xdg, "dfm")
}
